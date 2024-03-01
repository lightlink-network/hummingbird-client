package contracts

import "math/big"

const (
	// Status for a DA challenge that has not been initiated
	ChallengeDAStatusNone = 0
	// Status for a DA challenge initiated by a challenger
	ChallengeDAStatusChallengerInitiated = 1
	// Status for a DA challenge that has been won by the challenger
	ChallengeDAStatusChallengerWon = 2
	// Status for a DA challenge that has been won by the defender
	ChallengeDAStatusDefenderWon = 3
)

// Helper struct for pretty printing
type ChallengeDaInfo struct {
	BlockIndex *big.Int `pretty:"Block Index"`
	Challenger string   `pretty:"Challenger"`
	Expiry     *big.Int `pretty:"Expiry"`
	Status     uint8    `pretty:"Status"`
}

// Helper to convert challenge status enum to string
func DAChallengeStatusToString(c uint8) string {
	switch c {
	case ChallengeDAStatusNone:
		return "None"
	case ChallengeDAStatusChallengerInitiated:
		return "ChallengerInitiated"
	case ChallengeDAStatusChallengerWon:
		return "ChallengerWon"
	case ChallengeDAStatusDefenderWon:
		return "DefenderWon"
	default:
		return "Unknown"
	}
}
