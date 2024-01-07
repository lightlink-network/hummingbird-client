package node

type Node struct {
	Ethereum
	Celestia
	LightLink

	Store KVStore
}
