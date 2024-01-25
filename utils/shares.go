package utils

import (
	"github.com/celestiaorg/celestia-app/pkg/appconsts"
	"github.com/celestiaorg/celestia-app/pkg/shares"
	"github.com/celestiaorg/celestia-node/blob"
	"github.com/celestiaorg/celestia-node/share"
	coretypes "github.com/tendermint/tendermint/types"
)

func BytesToBlob(ns string, buf []byte) (*blob.Blob, error) {
	// get the namespace
	_ns, err := share.NewBlobNamespaceV0([]byte(ns))
	if err != nil {
		return nil, err
	}

	return blob.NewBlob(0, _ns, buf)
}

func BlobToShares(b *blob.Blob) ([]shares.Share, error) {
	_b := coretypes.Blob{
		NamespaceID:      b.NamespaceId,
		Data:             b.Data,
		ShareVersion:     uint8(b.ShareVersion),
		NamespaceVersion: uint8(b.NamespaceVersion),
	}

	return shares.SplitBlobs(_b)
}

func NSSharesToShares(ns share.NamespacedShares) []shares.Share {
	s := []shares.Share{}

	for _, row := range ns {
		for _, _nsShare := range row.Shares {
			_share, err := shares.NewShare(_nsShare)
			if err != nil {
				panic(err)
			}

			s = append(s, *_share)
		}
	}

	return s
}

// ShareDataStart returns the index of the first byte of the shares raw data.
// It is after the namespace, share info byte, sequence number (if present).
func ShareDataStart(s shares.Share) int {
	isStart, _ := s.IsSequenceStart()
	isCompact, _ := s.IsCompactShare()

	index := appconsts.NamespaceSize + appconsts.ShareInfoBytes
	if isStart {
		index += appconsts.SequenceLenBytes
	}
	if isCompact {
		index += appconsts.CompactShareReservedBytes
	}
	return index
}

// RawIndexToSharesIndex converts a raw index to a shares index, and a share index.
// Taking into account the namespace, share info byte, sequence number (if present).
func RawIndexToSharesIndex(rawIndex int, s []shares.Share) (share int, shareIndex int) {
	dataRead := 0

	for i := 0; i < len(s); i++ {
		rawData, _ := s[i].RawData()
		if rawIndex < dataRead+len(rawData) {
			share = i
			shareIndex = rawIndex - dataRead
			return
		}
		dataRead += len(rawData)
	}

	return
}
