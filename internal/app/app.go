package app

import (
	"github.com/edaywalid/reverse-proxy/internal/cache"
	"github.com/edaywalid/reverse-proxy/internal/config"
	"github.com/edaywalid/reverse-proxy/internal/handlers/proxy"
)

type App struct {
	Cache    *Cache
	Services *Services
	Config   *config.Config
	Handlers *Handlers
}

type (
	Cache struct {
		Cache *cache.Cache
	}
	Services struct {
		CacheService *cache.CacheService
	}
	Handlers struct {
		ProxyHandler *proxy.ProxyHandler
	}
)

func New(filePath string) (*App, error) {
	config, err := config.LoadConfig(filePath)
	if err != nil {
		return nil, err
	}

	app := &App{
		Config: config,
	}
	app.init()
	return app, nil
}

func (a *App) initCache() {
	a.Cache = &Cache{
		Cache: cache.NewCache(),
	}
}

func (a *App) initServices() {
	a.Services = &Services{
		CacheService: cache.NewCacheService(a.Cache.Cache),
	}
}

func (a *App) initHandlers() {
	a.Handlers = &Handlers{
		ProxyHandler: proxy.NewProxyHandler(a.Services.CacheService, a.Config),
	}
}

func (a *App) init() {
	a.initCache()
	a.initServices()
	a.initHandlers()

}
