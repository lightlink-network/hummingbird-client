// Copyright 2021 The go-ethereum Authors
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

package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

const DepositTxType = 0x7E

// DepositTx is
type DepositTx struct {
	ChainID  *big.Int        // destination chain ID
	Nonce    uint64          // nonce of sender account
	GasPrice *big.Int        // wei per gas
	Gas      uint64          // gas limit
	To       *common.Address `rlp:"nil"` // nil means contract creation
	Value    *big.Int        // wei amount
	Data     []byte          // contract invocation input data
	V, R, S  *big.Int        // signature values
}

// copy creates a deep copy of the transaction data and initializes all fields.
func (tx *DepositTx) copy() TxData {
	cpy := &DepositTx{
		Nonce: tx.Nonce,
		To:    copyAddressPtr(tx.To),
		Data:  common.CopyBytes(tx.Data),
		Gas:   tx.Gas,
		// These are initialized below.
		Value:    new(big.Int),
		ChainID:  new(big.Int),
		GasPrice: new(big.Int),
		V:        new(big.Int),
		R:        new(big.Int),
		S:        new(big.Int),
	}
	if tx.Value != nil {
		cpy.Value.Set(tx.Value)
	}
	if tx.ChainID != nil {
		cpy.ChainID.Set(tx.ChainID)
	}
	if tx.GasPrice != nil {
		cpy.GasPrice.Set(tx.GasPrice)
	}
	if tx.V != nil {
		cpy.V.Set(tx.V)
	}
	if tx.R != nil {
		cpy.R.Set(tx.R)
	}
	if tx.S != nil {
		cpy.S.Set(tx.S)
	}
	return cpy
}

// accessors for innerTx.
func (tx *DepositTx) txType() byte              { return DepositTxType }
func (tx *DepositTx) chainID() *big.Int         { return tx.ChainID }
func (tx *DepositTx) accessList() AccessList    { return nil }
func (tx *DepositTx) data() []byte              { return tx.Data }
func (tx *DepositTx) gas() uint64               { return tx.Gas }
func (tx *DepositTx) gasPrice() *big.Int        { return tx.GasPrice }
func (tx *DepositTx) gasTipCap() *big.Int       { return tx.GasPrice }
func (tx *DepositTx) gasFeeCap() *big.Int       { return tx.GasPrice }
func (tx *DepositTx) value() *big.Int           { return tx.Value }
func (tx *DepositTx) nonce() uint64             { return tx.Nonce }
func (tx *DepositTx) to() *common.Address       { return tx.To }
func (tx *DepositTx) blobGas() uint64           { return 0 }
func (tx *DepositTx) blobGasFeeCap() *big.Int   { return nil }
func (tx *DepositTx) blobHashes() []common.Hash { return nil }

func (tx *DepositTx) effectiveGasPrice(dst *big.Int, baseFee *big.Int) *big.Int {
	return dst.Set(tx.GasPrice)
}

func (tx *DepositTx) rawSignatureValues() (v, r, s *big.Int) {
	return tx.V, tx.R, tx.S
}

func (tx *DepositTx) setSignatureValues(chainID, v, r, s *big.Int) {
	tx.ChainID, tx.V, tx.R, tx.S = chainID, v, r, s
}
