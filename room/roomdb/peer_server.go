package roomdb

import (
	"errors"
)

type PeerServer struct {
	Addr string
}

func (s *PeerServer) Address() string {
	return s.Addr
}

var (
	registeredPeerServerAddr string = ""
)

func GetAvailablePeerServers() []PeerServer {
	// TODO: This should *somehow* ask k8s about available peer servers

	if registeredPeerServerAddr != "" {
		return []PeerServer{
			{
				Addr: registeredPeerServerAddr,
			},
		}
	}

	return nil
}

func SelectAvailablePeerServer() (*PeerServer, error) {
	availableServers := GetAvailablePeerServers()
	if len(availableServers) < 1 {
		return nil, errors.New("no peer server available")
	}

	s := selectServer(availableServers)
	return &s, nil
}

func selectServer(candidates []PeerServer) PeerServer {
	// TODO: actually implement load-balancing algorithm
	return candidates[0]
}

func RegisterPeerServer(addr string) {
	registeredPeerServerAddr = addr
}
