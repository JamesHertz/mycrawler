package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-kad-dht/dual"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"

	cidlib "github.com/ipfs/go-cid"
)


func newNode(bootstrapAddrs ...string) (host.Host, *dual.DHT, error) {

	bootstrapPeers := []peer.AddrInfo{}

	for _, pa := range bootstrapAddrs {
		pi, err := peer.AddrInfoFromString(pa)
		if err != nil {
			return nil, nil, err
		}
		bootstrapPeers = append(bootstrapPeers, *pi)
	}

	node, err := libp2p.New(
		libp2p.ListenAddrStrings("/ip4/127.0.0.1/tcp/0"),
	)

	if err != nil {
		return nil, nil, err
	}

	ctx := context.Background()

	// connect to bootstrap nodes
	for _, pi := range bootstrapPeers {
		node.Connect(ctx, pi)
	}

	kademlia, err := dual.New(
		ctx, node,
	)

	if err != nil {
		return nil, nil, err
	}

	if err := kademlia.Bootstrap(ctx); err != nil {
		return nil, nil, err
	}

	return node, kademlia, nil
}

func main(){
    var (
        cid string
        peers string
    )

    flag.StringVar(&cid, "cid", "", "the cid to resolve")
    flag.StringVar(&peers, "peers", "", "peers to connect in the format: multaddr-1,...,multiaddr-n")

    flag.Parse()
    if cid == "" || peers == ""{
        fmt.Println("ERROR: cid and peers are required")
        flag.Usage()
        os.Exit(1)
    }

	fmt.Println("--------------------------------------")
    fmt.Printf("cid=%s\npeers=%s\n", cid, peers)
	fmt.Println("--------------------------------------")
	fmt.Println()

	key, err := cidlib.Decode(cid) 
	if err != nil {
		fmt.Printf("Converting cid: %v\n", err)
		os.Exit(1)
	}

    h, dht, err := newNode(strings.Split(peers, ":")...)
    if err != nil {
        panic( err )
    }

	defer h.Close()

	for p := range dht.FindProvidersAsync(context.Background(), key, 0) {
		fmt.Printf("\033[92m+provider: %v\n\033[39m", p)
	}

	fmt.Println("\n+------------------------------------+")
    fmt.Println("|           ended look up            |")
	fmt.Println("+------------------------------------+")
}
