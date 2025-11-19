package utils

import (
	"github.com/celestiaorg/celestia-node/blob"
	"github.com/celestiaorg/celestia-openrpc/types/appconsts"
	"github.com/celestiaorg/go-square/v3/share"
)

func BytesToBlob(ns string, buf []byte) (*blob.Blob, error) {
	// get the namespace
	_ns, err := share.NewV0Namespace([]byte(ns))
	if err != nil {
		return nil, err
	}

	return blob.NewBlobV0(_ns, buf)
}

func BlobToShares(b *blob.Blob) ([]share.Share, error) {
	return b.Blob.ToShares()
}

func NSSharesToShares(ns []share.Share) []share.Share {
	// Already v3 shares, just return them
	return ns
}

// ShareDataStart returns the index of the first byte of the shares raw data.
// It is after the namespace, share info byte, sequence number (if present).
func ShareDataStart(s share.Share) int {
	isStart := s.IsSequenceStart()
	isCompact := s.IsCompactShare()

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
func RawIndexToSharesIndex(rawIndex int, s []share.Share) (shareIdx int, shareIndex int) {
	dataRead := 0

	for i := 0; i < len(s); i++ {
		rawData := s[i].RawData()
		if rawIndex < dataRead+len(rawData) {
			shareIdx = i
			shareIndex = rawIndex - dataRead
			return
		}
		dataRead += len(rawData)
	}

	return
}

func ExtractDataFromShares(s []share.Share) []byte {
	data := []byte{}

	for _, sh := range s {
		d := sh.RawData()

		data = append(data, d...)
	}

	return data
}

func BytesToShares(buf [][]byte) ([]share.Share, error) {
	s := []share.Share{}
	for _, b := range buf {
		sh, err := share.NewShare(b)
		if err != nil {
			return nil, err
		}

		s = append(s, *sh)
	}

	return s, nil
}
