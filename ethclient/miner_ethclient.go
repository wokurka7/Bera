// Copyright 2016 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

// Package ethclient provides a client for the Ethereum RPC API.
package ethclient

import (
	"context"

	"github.com/ethereum/go-ethereum/beacon/engine"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/miner"
)

func (ec *Client) Etherbase(ctx context.Context) (common.Address, error) {
	var miner common.Address
	if err := ec.c.CallContext(ctx, &miner, "miner_etherbase"); err != nil {
		return common.Address{}, err
	}
	return miner, nil
}

func (ec *Client) BuildBlock(ctx context.Context, attrs *miner.BuildPayloadArgs) (*engine.ExecutionPayloadEnvelope, error) {
	var payload engine.ExecutionPayloadEnvelope
	if err := ec.c.CallContext(ctx, &payload, "miner_buildBlock", attrs); err != nil {
		return nil, err
	}
	return &payload, nil
}

func (ec *Client) NewPayloadV3(ctx context.Context, params engine.ExecutableData, versionedHashes []common.Hash, beaconRoot *common.Hash) (engine.PayloadStatusV1, error) {
	var payloadStatus engine.PayloadStatusV1
	err := ec.c.CallContext(ctx, &payloadStatus, "miner_newPayloadV3", params, versionedHashes, beaconRoot)
	return payloadStatus, err
}
