package P2P_Net

import (
	"context"
	"fmt"
	"github.com/libp2p/go-libp2p"
	peerstore "github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/protocol/ping"
	multiaddr "github.com/multiformats/go-multiaddr"
	"os"
	"os/signal"
	"syscall"
)

func P2pPing() {

	node, err := libp2p.New(libp2p.ListenAddrStrings("/ip4/127.0.0.1/tcp/2000"),
		libp2p.Ping(false),
	)
	if err != nil {
		panic(err)
	}

	defer node.Close()
	fmt.Println("Listen address:", node.Addrs())

	pingService := &ping.PingService{Host: node}
	node.SetStreamHandler(ping.ID, pingService.PingHandler)

	peerInfo := peerstore.AddrInfo{
		ID:    node.ID(),
		Addrs: node.Addrs(),
	}
	addrs, err := peerstore.AddrInfoToP2pAddrs(&peerInfo)
	fmt.Println("libp2p node address:", addrs[0])

	if len(os.Args) > 1 {
		addr, err := multiaddr.NewMultiaddr(os.Args[1])
		if err != nil {
			panic(err)
		}
		peer, err := peerstore.AddrInfoFromP2pAddr(addr)
		if err != nil {
			panic(err)
		}
		fmt.Println("sending 5 ping messages to ", addr)
		ch := pingService.Ping(context.Background(), peer.ID)
		for i := 0; i < 5; i++ {
			res := <-ch
			fmt.Println("got ping response!", "RTT:", res.RTT)
		}
	} else {
		//wait for cmd to shut dwon
		sigCh := make(chan os.Signal)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		fmt.Println("Received signal , shutting down...")
	}
}
