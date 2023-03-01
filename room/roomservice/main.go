package main

import (
	"flag"

	"github.com/LassiHeikkila/flmnchll/room/roomdb"
)

func main() {
	var (
		peerServerAddr = ""
	)

	flag.StringVar(&peerServerAddr, "peerjs", "", "Address of PeerJS server")

	if peerServerAddr != "" {
		roomdb.RegisterPeerServer(peerServerAddr)
	}
}
