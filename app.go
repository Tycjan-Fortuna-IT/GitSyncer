package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"GitSyncer/core/database"
	"GitSyncer/core/store"
)

type App struct {
	ctx context.Context
	db  *sql.DB

	Providers    *store.ProviderStore
	Repositories *store.RepositoryStore
	Credentials  *store.CredentialStore
}

func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	dbPath, err := database.DefaultDBPath()
	if err != nil {
		log.Fatalf("failed to resolve database path: %v", err)
	}

	db, err := database.Open(dbPath)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	a.db = db

	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	a.Providers = store.NewProviderStore(db)
	a.Repositories = store.NewRepositoryStore(db)
	a.Credentials = store.NewCredentialStore(db)
}

func (a *App) shutdown(ctx context.Context) {
	if err := database.Close(a.db); err != nil {
		log.Printf("error closing database: %v", err)
	}
}

func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
