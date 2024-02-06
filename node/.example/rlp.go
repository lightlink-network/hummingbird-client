package main

import (
	"bytes"
	"fmt"
	"hummingbird/node"
	"hummingbird/utils"
)

func main() {
	n, log := setup()

	log.Info("Fetch Bundle from L2")
	bundle, err := fetchBundleFromL2(n)
	panicIfError(err, "failed to fetch bundle from L2")

	log.Info("Test 1.A Encode bundle to RLP")
	enc1, err := bundle.EncodeRLP()
	panicIfError(err, "failed to encode bundle to RLP")

	log.Info("Test 1.B Decode bundle from RLP")
	dec1 := &node.Bundle{}
	panicIfError(dec1.DecodeRLP(enc1), "failed to decode bundle from RLP")

	log.Info("Convert Bundle to Blob")
	blob, err := bundle.Blob(n.Namespace())
	panicIfError(err, "failed to get blob from bundle")

	log.Info("Convert Blob to Shares")
	s, err := utils.BlobToShares(blob)
	panicIfError(err, "failed to get shares from blob")

	shareData := []byte{}
	for _, _share := range s {
		d, _ := _share.RawData()
		shareData = append(shareData, d...)
	}

	log.Info("Byte sizes", "bundle", len(enc1), "shares", len(shareData), "shares_rlp", rlpNextItemSize(shareData))

	log.Info("Test 2.Compare share data to RLP encoded bundle")
	if bytes.Compare(shareData[:len(enc1)], enc1) != 0 {
		log.Info("Shares", "len", len(shareData), "data", fmt.Sprintf("%x...%x", shareData[:10], shareData[len(shareData)-10:]))
		log.Info("Expect", "len", len(enc1), "data", fmt.Sprintf("%x...%x", enc1[:10], enc1[len(enc1)-10:]))

		log.Info("Shares", "len", len(enc1), "data", fmt.Sprintf("%x...%x", shareData[:10], shareData[len(enc1)-10:]))
		panic("share data and RLP encoded bundle do not match")
	}

}

func fetchBundleFromL2(n *node.Node) (*node.Bundle, error) {
	blocks, err := n.LightLink.GetBlocks(62207259, 62207259+10)
	if err != nil {
		return nil, err
	}

	return &node.Bundle{
		Blocks: blocks,
	}, nil
}

func rlpNextItemSize(data []byte) int {
	if len(data) == 0 {
		return -1
	}

	prefix := data[0]

	switch {
	case prefix <= 0x7f:
		// Single byte
		return 1

	case prefix <= 0xb7:
		// Short string
		return int(prefix - 0x80 + 1)

	case prefix <= 0xbf:
		// Long string
		lengthSize := int(prefix - 0xb7)
		if len(data) < lengthSize+1 {
			return -1
		}
		length := int(data[1])
		for i := 2; i < lengthSize+1; i++ {
			length = (length << 8) + int(data[i])
		}
		return length + lengthSize + 1

	case prefix <= 0xf7:
		// Short list
		return int(prefix - 0xc0 + 1)

	case prefix <= 0xff:
		// Long list
		lengthSize := int(prefix - 0xf7)
		if len(data) < lengthSize+1 {
			return -1
		}
		length := int(data[1])
		for i := 2; i < lengthSize+1; i++ {
			length = (length << 8) + int(data[i])
		}
		return length + lengthSize + 1

	default:
		return -1
	}
}
