package proxy

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/edaywalid/reverse-proxy/internal/cache"
	"github.com/edaywalid/reverse-proxy/internal/config"
	"github.com/rs/zerolog/log"
)

type ProxyHandler struct {
	cacheService *cache.CacheService
	config       *config.Config
}

func NewProxyHandler(cacheService *cache.CacheService, config *config.Config) *ProxyHandler {
	ph := &ProxyHandler{cacheService: cacheService, config: config}
	ph.InitServers()
	return ph
}

func (p *ProxyHandler) HttpRedirect(w http.ResponseWriter, r *http.Request) {
	host, _, _ := net.SplitHostPort(r.Host)
	u := r.URL
	port := fmt.Sprintf("%d", p.config.HTTPSPort)
	u.Host = net.JoinHostPort(host, port)
	u.Scheme = "https"
	log.Info().Msgf("Redirecting to %s", u.String())
	http.Redirect(w, r, u.String(), http.StatusMovedPermanently)
}

func (p *ProxyHandler) Handler(w http.ResponseWriter, r *http.Request) {
	nextServer := p.getNextServer()

	if item, ok := p.cacheService.Get(r.RequestURI); ok {
		log.Info().Msgf("Cache hit for %s", r.RequestURI)

		for key, values := range item.Headers {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}

		w.WriteHeader(item.Status)
		io.Copy(w, bytes.NewReader(item.Body))
		return
	}
	log.Info().Msgf("Cache miss for %s", r.RequestURI)

	proxyReq, err := http.NewRequest(r.Method, nextServer+r.RequestURI, r.Body)

	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	log.Info().Msgf("Proxying request to %s", nextServer+r.RequestURI)
	proxyReq.Header = r.Header

	resp, err := http.DefaultClient.Do(proxyReq)
	if err != nil {
		http.Error(w, "Error contacting backend", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	p.cacheService.Set(r.RequestURI, resp)
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
