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

package vm

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/params"
	"github.com/holiman/uint256"
)

// StateDB is an EVM database for full state querying.
type StateDB interface {
	state.StateDBI
	// TODO: streamline this interface
	// CreateAccount(common.Address)

	// SubBalance(common.Address, *big.Int)
	// AddBalance(common.Address, *big.Int)
	// GetBalance(common.Address) *big.Int

	// GetNonce(common.Address) uint64
	// SetNonce(common.Address, uint64)

	// GetCodeHash(common.Address) common.Hash
	// GetCode(common.Address) []byte
	// SetCode(common.Address, []byte)
	// GetCodeSize(common.Address) int

	// AddRefund(uint64)
	// SubRefund(uint64)
	// GetRefund() uint64

	// GetCommittedState(common.Address, common.Hash) common.Hash
	// GetState(common.Address, common.Hash) common.Hash
	// SetState(common.Address, common.Hash, common.Hash)

	// GetTransientState(addr common.Address, key common.Hash) common.Hash
	// SetTransientState(addr common.Address, key, value common.Hash)

	// SelfDestruct(common.Address)
	// HasSelfDestructed(common.Address) bool

	// Selfdestruct6780(common.Address)

	// // Exist reports whether the given account exists in state.
	// // Notably this should also return true for self-destructed accounts.
	// Exist(common.Address) bool
	// // Empty returns whether the given account is empty. Empty
	// // is defined according to EIP161 (balance = nonce = code = 0).
	// Empty(common.Address) bool

	// AddressInAccessList(addr common.Address) bool
	// SlotInAccessList(addr common.Address, slot common.Hash) (addressOk bool, slotOk bool)
	// // AddAddressToAccessList adds the given address to the access list. This operation is safe to perform
	// // even if the feature/fork is not active yet
	// AddAddressToAccessList(addr common.Address)
	// // AddSlotToAccessList adds the given (address,slot) to the access list. This operation is safe to perform
	// // even if the feature/fork is not active yet
	// AddSlotToAccessList(addr common.Address, slot common.Hash)
	// Prepare(rules params.Rules, sender, coinbase common.Address, dest *common.Address, precompiles []common.Address, txAccesses types.AccessList)

	// RevertToSnapshot(int)
	// Snapshot() int

	// AddLog(*types.Log)
	// AddPreimage(common.Hash, []byte)

	// Error() error
}

// CallContext provides a basic interface for the EVM calling conventions. The EVM
// depends on this context being implemented for doing subcalls and initialising new EVM contracts.
type CallContext interface {
	// Call calls another contract.
	Call(env *EVM, me ContractRef, addr common.Address, data []byte, gas, value *big.Int) ([]byte, error)
	// CallCode takes another contracts code and execute within our own context
	CallCode(env *EVM, me ContractRef, addr common.Address, data []byte, gas, value *big.Int) ([]byte, error)
	// DelegateCall is same as CallCode except sender and value is propagated from parent to child scope
	DelegateCall(env *EVM, me ContractRef, addr common.Address, data []byte, gas *big.Int) ([]byte, error)
	// Create creates a new contract
	Create(env *EVM, me ContractRef, data []byte, gas, value *big.Int) ([]byte, common.Address, error)
}

type (
	// PrecompileManager allows the EVM to execute a precompiled contract.
	PrecompileManager interface {
		// `Has` returns if a precompiled contract was found at `addr`.
		Has(addr common.Address) bool

		// `Get` returns the precompiled contract at `addr`. Returns nil if no
		// contract is found at `addr`.
		Get(addr common.Address) PrecompiledContract

		GetActive(*params.Rules) []common.Address

		// `Run` runs a precompiled contract and returns the remaining gas.
		Run(evm PrecompileEVM, p PrecompiledContract, input []byte, caller common.Address,
			value *big.Int, suppliedGas uint64, readonly bool,
		) (ret []byte, remainingGas uint64, err error)
	}

	// PrecompileEVM is the interface through which stateful precompiles can call back into the EVM.
	PrecompileEVM interface {
		GetStateDB() StateDB

		Call(caller ContractRef, addr common.Address, input []byte, gas uint64, value *big.Int) (ret []byte, leftOverGas uint64, err error)
		StaticCall(caller ContractRef, addr common.Address, input []byte, gas uint64) (ret []byte, leftOverGas uint64, err error)
		Create(caller ContractRef, code []byte, gas uint64, value *big.Int) (ret []byte, contractAddr common.Address, leftOverGas uint64, err error)
		Create2(caller ContractRef, code []byte, gas uint64, endowment *big.Int, salt *uint256.Int) (ret []byte, contractAddr common.Address, leftOverGas uint64, err error)
		GetContext() *BlockContext
	}
)
