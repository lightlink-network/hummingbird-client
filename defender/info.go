package defender

import (
	"hummingbird/node/contracts"

	"github.com/ethereum/go-ethereum/common"
)

func (d *Defender) InfoDA(block common.Hash, pointer uint8, share uint32) (contracts.ChallengeDaInfo, error) {
	return d.Ethereum.GetDataRootInclusionChallenge(block, pointer, share)
}
