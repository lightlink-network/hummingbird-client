package contracts

import "math/big"

// Helper struct for pretty printing
type ChallengeDaInfo struct {
	BlockIndex *big.Int `pretty:"Block Index"`
	Challenger string   `pretty:"Challenger"`
	Expiry     *big.Int `pretty:"Expiry"`
	Status     uint8    `pretty:"Status"`
}

// Helper to convert challenge status enum to string
func StatusString(c uint8) string {
	switch c {
	case 0:
		return "None"
	case 1:
		return "ChallengerInitiated"
	case 2:
		return "DefenderResponded"
	case 3:
		return "ChallengerWon"
	default:
		return "Unknown"
	}
}
