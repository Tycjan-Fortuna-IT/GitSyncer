package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"GitSyncer/core/database"
	"GitSyncer/core/models"
	"GitSyncer/core/service"
	"GitSyncer/core/store"
)

type App struct {
	ctx context.Context
	db  *sql.DB

	Providers    *store.ProviderStore
	Repositories *store.RepositoryStore
	Credentials  *service.CredentialService
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

	credStore := store.NewCredentialStore(db)
	settingStore := store.NewSettingStore(db)
	a.Credentials = service.NewCredentialService(db, credStore, settingStore)
}

func (a *App) shutdown(ctx context.Context) {
	if a.Credentials != nil {
		a.Credentials.Lock()
	}

	if err := database.Close(a.db); err != nil {
		log.Printf("error closing database: %v", err)
	}
}

func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// IsMasterPasswordSetup returns whether a master password has been configured.
func (a *App) IsMasterPasswordSetup() (bool, error) {
	return a.Credentials.IsSetup()
}

// SetupMasterPassword configures the master password for first-time use.
func (a *App) SetupMasterPassword(password string) error {
	return a.Credentials.SetupMasterPassword(password)
}

// UnlockVault verifies the master password and unlocks credential access.
func (a *App) UnlockVault(password string) error {
	return a.Credentials.Unlock(password)
}

// LockVault locks the credential service, zeroing the encryption key.
func (a *App) LockVault() {
	a.Credentials.Lock()
}

// IsVaultLocked returns whether the credential service is currently locked.
func (a *App) IsVaultLocked() bool {
	return a.Credentials.IsLocked()
}

// ChangeMasterPassword re-encrypts all credentials with a new master password.
func (a *App) ChangeMasterPassword(oldPassword, newPassword string) error {
	return a.Credentials.ChangeMasterPassword(oldPassword, newPassword)
}

// StoreCredential encrypts and stores a new credential.
func (a *App) StoreCredential(providerID int64, label, authType, authData string) (int64, error) {
	cred := &models.Credential{
		ProviderID: providerID,
		Label:      label,
		AuthType:   authType,
		AuthData:   authData,
	}

	if err := a.Credentials.Store(cred); err != nil {
		return 0, err
	}

	return cred.ID, nil
}

// GetCredential retrieves and decrypts a credential by ID.
func (a *App) GetCredential(id int64) (*models.Credential, error) {
	return a.Credentials.GetByID(id)
}

// GetCredentialsByProvider retrieves and decrypts all credentials for a provider.
func (a *App) GetCredentialsByProvider(providerID int64) ([]models.Credential, error) {
	return a.Credentials.GetByProviderID(providerID)
}

// ListCredentials retrieves and decrypts all credentials.
func (a *App) ListCredentials() ([]models.Credential, error) {
	return a.Credentials.List()
}

// UpdateCredential re-encrypts and updates a credential.
func (a *App) UpdateCredential(id, providerID int64, label, authType, authData string) error {
	cred := &models.Credential{
		ID:         id,
		ProviderID: providerID,
		Label:      label,
		AuthType:   authType,
		AuthData:   authData,
	}

	return a.Credentials.Update(cred)
}

// DeleteCredential removes a credential by ID.
func (a *App) DeleteCredential(id int64) error {
	return a.Credentials.Delete(id)
}
