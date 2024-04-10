package contracts

import (
	"math/big"

	challenge "hummingbird/node/contracts/Challenge.sol"

	"github.com/ethereum/go-ethereum/common"
)

const (
	// Status for a DA challenge that has not been initiated
	ChallengeDAStatusNone = 0
	// Status for a DA challenge initiated by a challenger
	ChallengeDAStatusChallengerInitiated = 1
	// Status for a DA challenge that has been won by the challenger
	ChallengeDAStatusChallengerWon = 2
	// Status for a DA challenge that has been won by the defender
	ChallengeDAStatusDefenderWon = 3

	ChallengeL2HeaderStatusNone                = 0
	ChallengeL2HeaderStatusChallengerInitiated = 1
	ChallengeL2HeaderStatusChallengerWon       = 3
	ChallengeL2HeaderStatusDefenderWon         = 2
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

type L2HeaderChallengeInfo struct {
	Header       challenge.ChallengeL2HeaderL2HeaderPointer
	PrevHeader   challenge.ChallengeL2HeaderL2HeaderPointer
	ChallengeEnd *big.Int
	Challenger   common.Address
	Status       uint8
}
