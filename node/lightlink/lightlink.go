package lightlink

import (
	"encoding/json"
	"errors"
	"math/big"

	"hummingbird/node/lightlink/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

const DepositTxType = 0x7E
const DepositV2TxType = 125

type txJSON struct {
	Type hexutil.Uint64 `json:"type"`

	ChainID              *hexutil.Big      `json:"chainId,omitempty"`
	Nonce                *hexutil.Uint64   `json:"nonce"`
	To                   *common.Address   `json:"to"`
	Gas                  *hexutil.Uint64   `json:"gas"`
	GasPrice             *hexutil.Big      `json:"gasPrice"`
	MaxPriorityFeePerGas *hexutil.Big      `json:"maxPriorityFeePerGas"`
	MaxFeePerGas         *hexutil.Big      `json:"maxFeePerGas"`
	MaxFeePerDataGas     *hexutil.Big      `json:"maxFeePerDataGas,omitempty"`
	Value                *hexutil.Big      `json:"value"`
	Input                *hexutil.Bytes    `json:"input"`
	AccessList           *types.AccessList `json:"accessList,omitempty"`
	BlobVersionedHashes  []common.Hash     `json:"blobVersionedHashes,omitempty"`
	V                    *hexutil.Big      `json:"v"`
	R                    *hexutil.Big      `json:"r"`
	S                    *hexutil.Big      `json:"s"`

	// Deposit
	SourceHash          *common.Hash    `json:"sourceHash,omitempty"`
	From                *common.Address `json:"from,omitempty"`
	Mint                *hexutil.Big    `json:"mint,omitempty"`
	IsSystemTransaction bool            `json:"isSystemTransaction,omitempty"`

	// Only used for encoding:
	Hash common.Hash `json:"hash"`
}

func UnMarshallTx(input []byte) (*types.Transaction, error) {

	var dec txJSON
	if err := json.Unmarshal(input, &dec); err != nil {
		return nil, err
	}

	// Decode / verify fields according to transaction type.
	var inner types.TxData
	switch dec.Type {
	case types.LegacyTxType:
		var itx types.LegacyTx
		inner = &itx
		if dec.Nonce == nil {
			return nil, errors.New("missing required field 'nonce' in transaction")
		}
		itx.Nonce = uint64(*dec.Nonce)
		if dec.To != nil {
			itx.To = dec.To
		}
		if dec.Gas == nil {
			return nil, errors.New("missing required field 'gas' in transaction")
		}
		itx.Gas = uint64(*dec.Gas)
		if dec.GasPrice == nil {
			return nil, errors.New("missing required field 'gasPrice' in transaction")
		}
		itx.GasPrice = (*big.Int)(dec.GasPrice)
		if dec.Value == nil {
			return nil, errors.New("missing required field 'value' in transaction")
		}
		itx.Value = (*big.Int)(dec.Value)
		if dec.Input == nil {
			return nil, errors.New("missing required field 'input' in transaction")
		}
		itx.Data = *dec.Input
		if dec.V == nil {
			return nil, errors.New("missing required field 'v' in transaction")
		}
		itx.V = (*big.Int)(dec.V)
		if dec.R == nil {
			return nil, errors.New("missing required field 'r' in transaction")
		}
		itx.R = (*big.Int)(dec.R)
		if dec.S == nil {
			return nil, errors.New("missing required field 's' in transaction")
		}
		itx.S = (*big.Int)(dec.S)
		withSignature := itx.V.Sign() != 0 || itx.R.Sign() != 0 || itx.S.Sign() != 0
		if withSignature {
			if err := sanityCheckSignature(itx.V, itx.R, itx.S, true); err != nil {
				return nil, err
			}
		}

	case types.AccessListTxType:
		var itx types.AccessListTx
		inner = &itx
		if dec.ChainID == nil {
			return nil, errors.New("missing required field 'chainId' in transaction")
		}
		itx.ChainID = (*big.Int)(dec.ChainID)
		if dec.Nonce == nil {
			return nil, errors.New("missing required field 'nonce' in transaction")
		}
		itx.Nonce = uint64(*dec.Nonce)
		if dec.To != nil {
			itx.To = dec.To
		}
		if dec.Gas == nil {
			return nil, errors.New("missing required field 'gas' in transaction")
		}
		itx.Gas = uint64(*dec.Gas)
		if dec.GasPrice == nil {
			return nil, errors.New("missing required field 'gasPrice' in transaction")
		}
		itx.GasPrice = (*big.Int)(dec.GasPrice)
		if dec.Value == nil {
			return nil, errors.New("missing required field 'value' in transaction")
		}
		itx.Value = (*big.Int)(dec.Value)
		if dec.Input == nil {
			return nil, errors.New("missing required field 'input' in transaction")
		}
		itx.Data = *dec.Input
		if dec.V == nil {
			return nil, errors.New("missing required field 'v' in transaction")
		}
		if dec.AccessList != nil {
			itx.AccessList = *dec.AccessList
		}
		itx.V = (*big.Int)(dec.V)
		if dec.R == nil {
			return nil, errors.New("missing required field 'r' in transaction")
		}
		itx.R = (*big.Int)(dec.R)
		if dec.S == nil {
			return nil, errors.New("missing required field 's' in transaction")
		}
		itx.S = (*big.Int)(dec.S)
		withSignature := itx.V.Sign() != 0 || itx.R.Sign() != 0 || itx.S.Sign() != 0
		if withSignature {
			if err := sanityCheckSignature(itx.V, itx.R, itx.S, false); err != nil {
				return nil, err
			}
		}

	case types.DynamicFeeTxType:
		var itx types.DynamicFeeTx
		inner = &itx
		if dec.ChainID == nil {
			return nil, errors.New("missing required field 'chainId' in transaction")
		}
		itx.ChainID = (*big.Int)(dec.ChainID)
		if dec.Nonce == nil {
			return nil, errors.New("missing required field 'nonce' in transaction")
		}
		itx.Nonce = uint64(*dec.Nonce)
		if dec.To != nil {
			itx.To = dec.To
		}
		if dec.Gas == nil {
			return nil, errors.New("missing required field 'gas' for txdata")
		}
		itx.Gas = uint64(*dec.Gas)
		if dec.MaxPriorityFeePerGas == nil {
			return nil, errors.New("missing required field 'maxPriorityFeePerGas' for txdata")
		}
		itx.GasTipCap = (*big.Int)(dec.MaxPriorityFeePerGas)
		if dec.MaxFeePerGas == nil {
			return nil, errors.New("missing required field 'maxFeePerGas' for txdata")
		}
		itx.GasFeeCap = (*big.Int)(dec.MaxFeePerGas)
		if dec.Value == nil {
			return nil, errors.New("missing required field 'value' in transaction")
		}
		itx.Value = (*big.Int)(dec.Value)
		if dec.Input == nil {
			return nil, errors.New("missing required field 'input' in transaction")
		}
		itx.Data = *dec.Input
		if dec.V == nil {
			return nil, errors.New("missing required field 'v' in transaction")
		}
		if dec.AccessList != nil {
			itx.AccessList = *dec.AccessList
		}
		itx.V = (*big.Int)(dec.V)
		if dec.R == nil {
			return nil, errors.New("missing required field 'r' in transaction")
		}
		itx.R = (*big.Int)(dec.R)
		if dec.S == nil {
			return nil, errors.New("missing required field 's' in transaction")
		}
		itx.S = (*big.Int)(dec.S)
		withSignature := itx.V.Sign() != 0 || itx.R.Sign() != 0 || itx.S.Sign() != 0
		if withSignature {
			if err := sanityCheckSignature(itx.V, itx.R, itx.S, false); err != nil {
				return nil, err
			}
		}

	case DepositTxType:
		if dec.AccessList != nil || dec.MaxFeePerGas != nil ||
			dec.MaxPriorityFeePerGas != nil {
			return nil, errors.New("unexpected field(s) in deposit transaction")
		}
		if dec.GasPrice != nil && dec.GasPrice.ToInt().Cmp(common.Big0) != 0 {
			return nil, errors.New("deposit transaction GasPrice must be 0")
		}
		var itx types.DynamicFeeTx
		inner = &itx
		// if dec.ChainID == nil {
		// 	return nil, errors.New("missing required field 'chainId' in transaction")
		// }
		// // uh oh: WE NEED THIS.
		// // fix later
		// itx.ChainID = (*big.Int)(dec.ChainID)
		// if dec.Nonce == nil {
		// 	return nil, errors.New("missing required field 'nonce' in transaction")
		// }
		itx.Nonce = uint64(*dec.Nonce)
		if dec.To != nil {
			itx.To = dec.To
		}
		if dec.Gas == nil {
			return nil, errors.New("missing required field 'gas' in transaction")
		}
		itx.Gas = uint64(*dec.Gas)
		if dec.GasPrice == nil {
			return nil, errors.New("missing required field 'gasPrice' in transaction")
		}
		itx.GasFeeCap = (*big.Int)(dec.GasPrice)
		itx.Value = (*big.Int)(dec.Value)
		if dec.Input == nil {
			return nil, errors.New("missing required field 'input' in transaction")
		}
		itx.Data = *dec.Input
		if dec.V == nil {
			return nil, errors.New("missing required field 'v' in transaction")
		}
		itx.V = (*big.Int)(dec.V)
		if dec.R == nil {
			return nil, errors.New("missing required field 'r' in transaction")
		}
		itx.R = (*big.Int)(dec.R)
		if dec.S == nil {
			return nil, errors.New("missing required field 's' in transaction")
		}
		itx.S = (*big.Int)(dec.S)
		withSignature := itx.V.Sign() != 0 || itx.R.Sign() != 0 || itx.S.Sign() != 0
		if withSignature {
			if err := sanityCheckSignature(itx.V, itx.R, itx.S, false); err != nil {
				return nil, err
			}
		}

	case DepositV2TxType:
		if dec.AccessList != nil || dec.MaxFeePerGas != nil ||
			dec.MaxPriorityFeePerGas != nil {
			return nil, errors.New("unexpected field(s) in deposit transaction")
		}
		if dec.GasPrice != nil && dec.GasPrice.ToInt().Cmp(common.Big0) != 0 {
			return nil, errors.New("deposit transaction GasPrice must be 0")
		}
		/*
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
		} */
		var itx types.DepositTxV2
		inner = &itx

		if dec.SourceHash == nil {
			return nil, errors.New("missing required field 'sourceHash' in transaction")
		}
		itx.SourceHash = *dec.SourceHash

		if dec.From == nil {
			return nil, errors.New("missing required field 'from' in transaction")
		}
		itx.From = *dec.From

		if dec.To != nil {
			itx.To = dec.To
		}

		if dec.Mint != nil {
			itx.Mint = (*big.Int)(dec.Mint)
		}

		if dec.Value != nil {
			itx.Value = (*big.Int)(dec.Value)
		}

		if dec.Gas == nil {
			return nil, errors.New("missing required field 'gas' in transaction")
		}
		itx.Gas = uint64(*dec.Gas)

		if dec.IsSystemTransaction {
			itx.IsSystemTransaction = dec.IsSystemTransaction
		}

		itx.Data = *dec.Input

	default:
		return nil, types.ErrTxTypeNotSupported
	}

	// Now set the inner transaction.
	return types.NewTx(inner), nil
}

func sanityCheckSignature(v *big.Int, r *big.Int, s *big.Int, maybeProtected bool) error {
	if isProtectedV(v) && !maybeProtected {
		return types.ErrUnexpectedProtection
	}

	var plainV byte
	if isProtectedV(v) {
		chainID := deriveChainId(v).Uint64()
		plainV = byte(v.Uint64() - 35 - 2*chainID)
	} else if maybeProtected {
		// Only EIP-155 signatures can be optionally protected. Since
		// we determined this v value is not protected, it must be a
		// raw 27 or 28.
		plainV = byte(v.Uint64() - 27)
	} else {
		// If the signature is not optionally protected, we assume it
		// must already be equal to the recovery id.
		plainV = byte(v.Uint64())
	}
	if !crypto.ValidateSignatureValues(plainV, r, s, false) {
		return types.ErrInvalidSig
	}

	return nil
}

func isProtectedV(V *big.Int) bool {
	if V.BitLen() <= 8 {
		v := V.Uint64()
		return v != 27 && v != 28 && v != 1 && v != 0
	}
	// anything not 27 or 28 is considered protected
	return true
}

// deriveChainId derives the chain id from the given v parameter
func deriveChainId(v *big.Int) *big.Int {
	if v.BitLen() <= 64 {
		v := v.Uint64()
		if v == 27 || v == 28 {
			return new(big.Int)
		}
		return new(big.Int).SetUint64((v - 35) / 2)
	}
	v = new(big.Int).Sub(v, big.NewInt(35))
	return v.Div(v, big.NewInt(2))
}
