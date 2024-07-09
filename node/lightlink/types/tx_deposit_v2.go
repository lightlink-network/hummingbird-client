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
	"bytes"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
)

const DepositTxV2Type = 0x7D

type DepositTxV2 struct {
	// SourceHash uniquely identifies the source of the deposit
	SourceHash common.Hash
	// From is exposed through the types.Signer, not through TxData
	From common.Address
	// nil means contract creation
	To *common.Address `rlp:"nil"`
	// Mint is minted on L2, locked on L1, nil if no minting.
	Mint *big.Int `rlp:"nil"`
	// Value is transferred from L2 balance, executed after Mint (if any)
	Value *big.Int
	// gas limit
	Gas uint64
	// Field indicating if this transaction is exempt from the L2 gas limit.
	IsSystemTransaction bool
	// Normal Tx data
	Data []byte
}

// copy creates a deep copy of the transaction data and initializes all fields.
func (tx *DepositTxV2) copy() TxData {
	cpy := &DepositTxV2{
		SourceHash:          tx.SourceHash,
		From:                tx.From,
		To:                  copyAddressPtr(tx.To),
		Mint:                nil,
		Value:               new(big.Int),
		Gas:                 tx.Gas,
		IsSystemTransaction: tx.IsSystemTransaction,
		Data:                common.CopyBytes(tx.Data),
	}
	if tx.Mint != nil {
		cpy.Mint = new(big.Int).Set(tx.Mint)
	}
	if tx.Value != nil {
		cpy.Value.Set(tx.Value)
	}
	return cpy
}

// accessors for innerTx.
func (tx *DepositTxV2) txType() byte              { return DepositTxV2Type }
func (tx *DepositTxV2) chainID() *big.Int         { return common.Big0 }
func (tx *DepositTxV2) accessList() AccessList    { return nil }
func (tx *DepositTxV2) data() []byte              { return tx.Data }
func (tx *DepositTxV2) gas() uint64               { return tx.Gas }
func (tx *DepositTxV2) gasFeeCap() *big.Int       { return new(big.Int) }
func (tx *DepositTxV2) gasTipCap() *big.Int       { return new(big.Int) }
func (tx *DepositTxV2) gasPrice() *big.Int        { return new(big.Int) }
func (tx *DepositTxV2) value() *big.Int           { return tx.Value }
func (tx *DepositTxV2) nonce() uint64             { return 0 }
func (tx *DepositTxV2) to() *common.Address       { return tx.To }
func (tx *DepositTxV2) isSystemTx() bool          { return tx.IsSystemTransaction }
func (tx *DepositTxV2) blobGas() uint64           { return 0 }
func (tx *DepositTxV2) blobGasFeeCap() *big.Int   { return nil }
func (tx *DepositTxV2) blobHashes() []common.Hash { return nil }

func (tx *DepositTxV2) effectiveGasPrice(dst *big.Int, baseFee *big.Int) *big.Int {
	return dst.Set(new(big.Int))
}

func (tx *DepositTxV2) effectiveNonce() *uint64 { return nil }

func (tx *DepositTxV2) rawSignatureValues() (v, r, s *big.Int) {
	return common.Big0, common.Big0, common.Big0
}

func (tx *DepositTxV2) setSignatureValues(chainID, v, r, s *big.Int) {
	// this is a noop for deposit transactions
}

func (tx *DepositTxV2) encode(b *bytes.Buffer) error {
	return rlp.Encode(b, tx)
}

func (tx *DepositTxV2) decode(input []byte) error {
	return rlp.DecodeBytes(input, tx)
}
