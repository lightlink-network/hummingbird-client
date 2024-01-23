package main

import (
	"bytes"
	"crypto/ecdsa"
	"fmt"
	"hummingbird/cli/cmd"
	"hummingbird/config"
	"hummingbird/node"
	"hummingbird/utils"
	"log/slog"
	"os"

	"github.com/celestiaorg/celestia-node/share"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/spf13/viper"
)

func main() {
	viper.Set("config-path", "./ganache")
	viper.Set("log-level", "debug")
	cfg := config.Load()
	log := cmd.ConsoleLogger()
	ethKey := getEthKey()

	blockStart := 63922950
	blockEnd := 63922952

	n, err := node.NewFromConfig(cfg, log, ethKey)
	if err != nil {
		panic(err)
	}

	// 0. Get Blocks
	blocks, err := n.LightLink.GetBlocks(uint64(blockStart), uint64(blockEnd))
	if err != nil {
		log.Error("Failed to get blocks", "start", blockStart, "end", blockEnd, "err", err)
		panic(err)
	}
	log.Info("Got blocks", "start", blockStart, "end", blockEnd, "count", len(blocks))

	// print block hashes
	for _, b := range blocks {
		log.Info("Block", "index", b.NumberU64(), "hash", b.Hash().Hex(), "stateroot", b.Root().Hex())
	}

	// get nsShares
	nsShares := getShares(log, n, "0x2e6dcb22143fe94e1641cc65f6f5573064bcc5fdacaf9b115ea529ff5761e0d4")
	_shares := utils.NSSharesToShares(nsShares)

	// check if data is in the shares
	rawData := []byte{}
	for _, s := range _shares {
		raw, _ := s.RawData()
		rawData = append(rawData, raw...)
	}

	fmt.Println("Data hash", crypto.Keccak256Hash(rawData).Hex())

	encH0, _ := rlp.EncodeToBytes(blocks[0].Header())

	fmt.Println("Data found", bytes.Contains(rawData, encH0))

	// fmt.Printf("H0 %x\n\n", string(encH0))
	// fmt.Printf("R0 %x\n\n", string(rawData))

	// 2. Get proof
	//getShareProofs(log, n, "0x2e6dcb22143fe94e1641cc65f6f5573064bcc5fdacaf9b115ea529ff5761e0d4", blocks)
}

func getShareProofs(log *slog.Logger, n *node.Node, txHash string, blocks []*types.Block) {
	b := &node.Bundle{Blocks: blocks}
	pointer, err := b.FindHeaderShares(b.Blocks[1].Hash(), n.Celestia.Namespace())
	if err != nil {
		log.Error("Failed to find header shares", "err", err)
		panic(err)
	}

	celTxHash := common.HexToHash(txHash)
	proof, err := n.Celestia.GetShareProofs(celTxHash[:], pointer)
	if err != nil {
		log.Error("Failed to get share proofs", "err", err)
		panic(err)
	}

	//log.Info("Proof generated", "proof", fmt.Sprintf("%+v", proof))

	// hash the proof datas
	for i, p := range proof.Data {
		log.Info("Proof data", "share", i, "hash", crypto.Keccak256Hash(p).Hex())
		//log.Info("Raw", "share", i, "raw", fmt.Sprintf("%s", p))
	}

}

func getShares(log *slog.Logger, n *node.Node, txHash string) share.NamespacedShares {
	celTxHash := common.HexToHash(txHash)

	shares, err := n.Celestia.GetShares(celTxHash[:], n.Namespace())
	if err != nil {
		log.Error("Failed to get shares", "err", err)
		panic(err)
	}

	// log.Info("Shares found", "shares", fmt.Sprintf("%+v", shares))
	// print share hashes
	for row, rows := range shares {
		for col, s := range rows.Shares {
			log.Info("Share", "row", row, "col", col, "hash", crypto.Keccak256Hash(s).Hex())
		}
	}

	return shares
}

func getEthKey() *ecdsa.PrivateKey {
	key := os.Getenv("ETH_KEY")
	if key == "" {
		panic("env ETH_KEY not set")
	}

	ethKey, err := crypto.ToECDSA(hexutil.MustDecode(key))
	if err != nil {
		panic("Failed to decode ETH_KEY: " + err.Error())
	}

	return ethKey
}
