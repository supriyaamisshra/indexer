package main

import (
	"fmt"

	"github.com/cyberconnecthq/indexer/fetcher"
)

const (
	address = "0xd8da6bf26964af9d7eed9e03e53415d37aa96045" // vitalik.eth
	//address = "0x983110309620d911731ac0932219af06091b6744" // brantly.eth
	//address = "0x8d07D225a769b7Af3A923481E1FdF49180e6A265" // test address with a verified sybil twitter account
)

func main() {
	f := fetcher.NewFetcher()

	ids, err := f.FetchIdentity(address)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", ids)

	conn, err := f.FetchConnections(address)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", conn)
}
