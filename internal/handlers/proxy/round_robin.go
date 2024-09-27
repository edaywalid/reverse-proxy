package proxy

import (
	"fmt"
	"sync"

	"github.com/edaywalid/reverse-proxy/internal/config"
	"github.com/rs/zerolog/log"
)

var (
	servers     []config.ServerConfig
	serverIndex = 0
	serverLock  = sync.Mutex{}
)

func (p *ProxyHandler) InitServers() {
	log.Info().Msg("Initializing servers")
	servers = p.config.Servers
	log.Info().Msgf("Servers: %v", servers)
}

func (p *ProxyHandler) getNextServer() string {
	serverLock.Lock()
	defer serverLock.Unlock()

	nextServer := servers[serverIndex]
	serverIndex = (serverIndex + 1) % len(servers)
	return fmt.Sprintf("%s", nextServer.URL)
}
