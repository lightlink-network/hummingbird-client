package contracts

import "math/big"

// Helper struct for pretty printing
type ChallengeDaInfo struct {
	BlockIndex *big.Int `pretty:"Block Index"`
	Challenger string   `pretty:"Challenger"`
	Expiry     *big.Int `pretty:"Expiry"`
	Status     uint8    `pretty:"Status"`
}
