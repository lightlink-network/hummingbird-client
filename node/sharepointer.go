package node

import (
	"hummingbird/utils"

	"github.com/celestiaorg/celestia-app/pkg/shares"
)

type ShareRange struct {
	Start uint64
	End   uint64
}

// SharePointer is a pointer to some data inside a group of shares
type SharePointer struct {
	shares     []shares.Share
	StartShare int
	Ranges     []ShareRange
}

func NewSharePointer(_shares []shares.Share, startShare int, startIndex int, endShare int, endIndex int) *SharePointer {
	// add the start range
	ranges := []ShareRange{
		{
			Start: uint64(startIndex) + uint64(utils.ShareDataStart(_shares[startShare])),
			End:   uint64(_shares[startShare].Len()),
		},
	}

	// add the middle ranges
	for i := startShare + 1; i < endShare; i++ {
		r := ShareRange{
			Start: uint64(utils.ShareDataStart(_shares[i])),
			End:   uint64(_shares[i].Len()),
		}

		ranges = append(ranges, r)
	}

	// add the end range
	if startShare != endShare {
		ranges = append(ranges, ShareRange{
			Start: uint64(utils.ShareDataStart(_shares[endShare])),
		})
	}
	ranges[len(ranges)-1].End = uint64(endIndex) + uint64(utils.ShareDataStart(_shares[endShare]))

	return &SharePointer{
		shares:     _shares,
		StartShare: startShare,
		Ranges:     ranges,
	}
}

func (s *SharePointer) EndShare() int {
	return s.StartShare + len(s.Ranges) - 1
}

func (s *SharePointer) Bytes() []byte {
	data := []byte{}
	for i := 0; i < len(s.Ranges); i++ {
		data = append(data, s.shares[s.StartShare+i].ToBytes()[s.Ranges[i].Start:s.Ranges[i].End]...)
	}

	return data
}

func (s *SharePointer) Shares() []shares.Share {
	return s.shares[s.StartShare : s.EndShare()+1]
}

func (s *SharePointer) AllShares() []shares.Share {
	return s.shares
}
