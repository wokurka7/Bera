package eth

import (
	"errors"
	"time"

	"github.com/ethereum/go-ethereum/beacon/engine"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/miner"
	"github.com/ethereum/go-ethereum/params"
)

// BuildBlock is a convenience function to build a block.
func (s *MinerAPI) BuildBlock(attrs *miner.BuildPayloadArgs) (*engine.ExecutionPayloadEnvelope, error) {
	// Send a request to generate a full block in the background.
	// The result can be obtained via the returned channel.
	args := &miner.BuildPayloadArgs{
		Parent:       attrs.Parent,
		Timestamp:    uint64(attrs.Timestamp),
		FeeRecipient: attrs.FeeRecipient,
		Random:       attrs.Random,
		Withdrawals:  attrs.Withdrawals,
		BeaconRoot:   attrs.BeaconRoot,
	}

	payload, err := s.e.Miner().BuildPayload(args)
	if err != nil {
		log.Error("Failed to build payload", "err", err)
		return nil, err
	}

	resCh := make(chan *engine.ExecutionPayloadEnvelope, 1)
	go func() {
		resCh <- payload.ResolveFull()
	}()

	timer := time.NewTimer(4 * time.Second)
	defer timer.Stop()

	select {
	case payload := <-resCh:
		if payload == nil {
			return nil, errors.New("received nil payload from sealing work")
		}
		return payload, nil
	case <-timer.C:
		log.Error("timeout waiting for block", "parent hash", attrs.Parent)
		return nil, errors.New("timeout waiting for block result")
	}
}

func (s *MinerAPI) GetBlockByHash(hash common.Hash) *types.Block {
	return s.e.BlockChain().GetBlockByHash(hash)
}

func (s *MinerAPI) Config() *params.ChainConfig {
	return s.e.BlockChain().Config()
}

func (s *MinerAPI) Etherbase() common.Address {
	return s.e.etherbase
}

func (s *MinerAPI) NewPayloadV3(
	params engine.ExecutableData, versionedHashes []common.Hash, beaconRoot *common.Hash,
) (engine.PayloadStatusV1, error) {
	var (
		block *types.Block
		err   error
	)

	if block, err = engine.ExecutableDataToBlock(
		params, nil, nil,
	); err != nil {
		log.Error("failed to build evm block", "err", err)
		return engine.STATUS_INVALID.PayloadStatus, err
	}
	if err = s.e.BlockChain().InsertBlockWithoutSetHead(block); err != nil {
		log.Error("failed to insert evm block", "err", err)
		return engine.STATUS_INVALID.PayloadStatus, err
	}

	if _, err = s.e.BlockChain().SetCanonical(block); err != nil {
		return engine.STATUS_INVALID.PayloadStatus, err
	}

	insertedHash := s.e.BlockChain().CurrentBlock().Hash()

	return engine.PayloadStatusV1{
		Status: engine.VALID, LatestValidHash: &insertedHash,
	}, nil
}
