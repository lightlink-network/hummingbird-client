package utils

import (
	"github.com/celestiaorg/celestia-app/pkg/appconsts"
	"github.com/celestiaorg/celestia-app/pkg/shares"
	"github.com/celestiaorg/celestia-node/blob"
	openshare "github.com/celestiaorg/celestia-openrpc/types/share"
	"github.com/celestiaorg/go-square/v2/share"
	coretypes "github.com/tendermint/tendermint/types"
)

func BytesToBlob(ns string, buf []byte) (*blob.Blob, error) {
	// get the namespace
	_ns, err := share.NewV0Namespace([]byte(ns))
	if err != nil {
		return nil, err
	}

	return blob.NewBlobV0(_ns, buf)
}

func BlobToShares(b *blob.Blob) ([]shares.Share, error) {
	_b := coretypes.Blob{
		NamespaceID:      b.Namespace().ID(),
		Data:             b.Data(),
		ShareVersion:     uint8(b.ShareVersion()),
		NamespaceVersion: uint8(b.Namespace().Version()),
	}

	return shares.SplitBlobs(_b)
}

func NSSharesToShares(ns openshare.NamespacedShares) []shares.Share {
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

func ExtractDataFromShares(s []shares.Share) []byte {
	data := []byte{}

	for _, share := range s {
		d, err := share.RawData()
		if err != nil {
			panic(err)
		}

		data = append(data, d...)
	}

	return data
}

func BytesToShares(buf [][]byte) ([]shares.Share, error) {
	s := []shares.Share{}
	for _, b := range buf {
		share, err := shares.NewShare(b)
		if err != nil {
			return nil, err
		}

		s = append(s, *share)
	}

	return s, nil
}
