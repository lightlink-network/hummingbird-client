package utils

import (
	"github.com/celestiaorg/celestia-node/blob"
	"github.com/celestiaorg/celestia-openrpc/types/appconsts"
	squareblob "github.com/celestiaorg/go-square/blob"
	"github.com/celestiaorg/go-square/shares"
	"github.com/celestiaorg/go-square/v2/share"
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
	_b := &squareblob.Blob{
		NamespaceId:      b.Namespace().ID(),
		Data:             b.Data(),
		ShareVersion:     uint32(b.ShareVersion()),
		NamespaceVersion: uint32(b.Namespace().Version()),
	}
	return shares.SplitBlobs(_b)
}

func NSSharesToShares(ns []share.Share) []shares.Share {
	s := []shares.Share{}

	for _, _nsShare := range ns {
		_share, err := shares.NewShare(_nsShare.ToBytes())
		if err != nil {
			panic(err)
		}

		s = append(s, *_share)
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
