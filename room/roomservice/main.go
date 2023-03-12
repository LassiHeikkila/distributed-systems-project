package main

import (
	"flag"

	"github.com/LassiHeikkila/flmnchll/account/accountservice/accountclient"
)

func main() {
	var (
		peerServerAddr     string
		accountServiceAddr string
		dbPath             string
		allowedCORSOrigins string
		httpPort           uint
	)

	flag.StringVar(&peerServerAddr, "peerjs", "", "Address of PeerJS server")
	flag.StringVar(&accountServiceAddr, "accountServiceAddr", "", "Address of account-service")
	flag.StringVar(&dbPath, "db", "content.db", "Path to database file")
	flag.StringVar(&allowedCORSOrigins, "cors", "*", "Comma-separated list of accepted CORS origins")
	flag.UintVar(&httpPort, "httpPort", 8080, "HTTP port")

	if peerServerAddr != "" {
		RegisterPeerServer(peerServerAddr)
	}

	if accountServiceAddr != "" {
		accountclient.SetAccountServiceAddr(accountServiceAddr)
	}
}
