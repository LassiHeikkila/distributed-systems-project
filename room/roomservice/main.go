package main

import (
	"flag"
)

func main() {
	var (
		peerServerAddr = ""
	)

	flag.StringVar(&peerServerAddr, "peerjs", "", "Address of PeerJS server")

	if peerServerAddr != "" {
		RegisterPeerServer(peerServerAddr)
	}
}
