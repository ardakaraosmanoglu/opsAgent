package app

import (
    "context"
    "fmt"
    "log"
    "net/http"

    "github.com/opsagent/opsagent/internal/config"
    "github.com/opsagent/opsagent/internal/infrastructure/storage/sqlite"
    "github.com/opsagent/opsagent/internal/web"
)

type App struct {
    cfg    *config.Config
    server *http.Server
    store  *sqlite.Store
}

func New(cfg *config.Config) (*App, error) {
    store, err := sqlite.Open(cfg.Database.Path)
    if err != nil {
        return nil, fmt.Errorf("opening sqlite store: %w", err)
    }

    if err := store.Migrate(); err != nil {
        return nil, fmt.Errorf("running migrations: %w", err)
    }

    router := web.NewRouter(cfg, store)

    addr := fmt.Sprintf("%s:%d", cfg.Server.BindAddress, cfg.Server.Port)
    server := &http.Server{
        Addr:    addr,
        Handler: router,
    }

    return &App{
        cfg:    cfg,
        server: server,
        store:  store,
    }, nil
}

func (a *App) Run() error {
    ctx := context.Background()

    go func() {
        if err := StartScheduler(ctx, a.cfg, a.store); err != nil {
            log.Printf("scheduler stopped: %v", err)
        }
    }()

    log.Printf("OpsAgent dashboard listening on %s", a.server.Addr)
    return a.server.ListenAndServe()
}
