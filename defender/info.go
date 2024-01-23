package defender

import (
	"hummingbird/node/contracts"

	"github.com/ethereum/go-ethereum/common"
)

func (d *Defender) InfoDA(block common.Hash) (contracts.ChallengeDaInfo, error) {
	return d.Ethereum.GetDataRootInclusionChallenge(block)
}
