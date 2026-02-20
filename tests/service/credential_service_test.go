package service_test

import (
	"testing"

	"GitSyncer/core/database"
	"GitSyncer/core/models"
	"GitSyncer/core/service"
	"GitSyncer/core/store"
)

// setupTestDB creates an in-memory SQLite database with migrations applied.
func setupTestDB(t *testing.T) (*service.CredentialService, *store.ProviderStore) {
	t.Helper()

	db, err := database.Open(":memory:")
	if err != nil {
		t.Fatalf("open test db: %v", err)
	}

	if err := database.RunMigrations(db); err != nil {
		t.Fatalf("run migrations: %v", err)
	}

	t.Cleanup(func() { db.Close() })

	providerStore := store.NewProviderStore(db)
	credStore := store.NewCredentialStore(db)
	settingStore := store.NewSettingStore(db)
	credService := service.NewCredentialService(db, credStore, settingStore)

	return credService, providerStore
}

// createTestProvider inserts a provider and returns its ID.
func createTestProvider(t *testing.T, ps *store.ProviderStore) int64 {
	t.Helper()

	p := &models.Provider{
		Name:    "test-provider",
		Type:    "github",
		BaseURL: "https://github.com",
	}

	if err := ps.Create(p); err != nil {
		t.Fatalf("create test provider: %v", err)
	}

	return p.ID
}

func TestSetupAndUnlock(t *testing.T) {
	svc, _ := setupTestDB(t)

	// Not set up initially
	setup, err := svc.IsSetup()
	if err != nil {
		t.Fatalf("IsSetup() error: %v", err)
	}

	if setup {
		t.Fatal("IsSetup() = true before setup")
	}

	// Should be locked initially
	if !svc.IsLocked() {
		t.Fatal("IsLocked() = false before setup")
	}

	// Setup master password
	if err := svc.SetupMasterPassword("my-master-password"); err != nil {
		t.Fatalf("SetupMasterPassword() error: %v", err)
	}

	// Now set up and unlocked
	setup, err = svc.IsSetup()
	if err != nil {
		t.Fatalf("IsSetup() after setup error: %v", err)
	}

	if !setup {
		t.Fatal("IsSetup() = false after setup")
	}

	if svc.IsLocked() {
		t.Fatal("IsLocked() = true after setup")
	}

	// Lock
	svc.Lock()

	if !svc.IsLocked() {
		t.Fatal("IsLocked() = false after Lock()")
	}

	// Unlock with correct password
	if err := svc.Unlock("my-master-password"); err != nil {
		t.Fatalf("Unlock() error: %v", err)
	}

	if svc.IsLocked() {
		t.Fatal("IsLocked() = true after Unlock()")
	}
}

func TestUnlockWithWrongPassword(t *testing.T) {
	svc, _ := setupTestDB(t)

	if err := svc.SetupMasterPassword("correct-password"); err != nil {
		t.Fatalf("SetupMasterPassword() error: %v", err)
	}

	svc.Lock()

	err := svc.Unlock("wrong-password")
	if err == nil {
		t.Fatal("Unlock() with wrong password should return error")
	}

	if err != service.ErrInvalidPassword {
		t.Fatalf("Unlock() error = %v, want ErrInvalidPassword", err)
	}

	if !svc.IsLocked() {
		t.Fatal("IsLocked() = false after failed Unlock()")
	}
}

func TestSetupMasterPasswordDuplicate(t *testing.T) {
	svc, _ := setupTestDB(t)

	if err := svc.SetupMasterPassword("password1"); err != nil {
		t.Fatalf("first SetupMasterPassword() error: %v", err)
	}

	err := svc.SetupMasterPassword("password2")
	if err == nil {
		t.Fatal("second SetupMasterPassword() should return error")
	}

	if err != service.ErrAlreadySetup {
		t.Fatalf("second SetupMasterPassword() error = %v, want ErrAlreadySetup", err)
	}
}

func TestStoreAndRetrieveCredential(t *testing.T) {
	svc, ps := setupTestDB(t)
	providerID := createTestProvider(t, ps)

	if err := svc.SetupMasterPassword("master-pass"); err != nil {
		t.Fatalf("SetupMasterPassword() error: %v", err)
	}

	// Store a credential
	cred := &models.Credential{
		ProviderID: providerID,
		Label:      "my-token",
		AuthType:   "token",
		AuthData:   "ghp_secret123abc",
	}

	if err := svc.Store(cred); err != nil {
		t.Fatalf("Store() error: %v", err)
	}

	if cred.ID == 0 {
		t.Fatal("Store() did not set credential ID")
	}

	// Retrieve by ID - should get back plaintext
	got, err := svc.GetByID(cred.ID)
	if err != nil {
		t.Fatalf("GetByID() error: %v", err)
	}

	if got.AuthData != "ghp_secret123abc" {
		t.Fatalf("GetByID().AuthData = %q, want %q", got.AuthData, "ghp_secret123abc")
	}

	if got.Label != "my-token" {
		t.Fatalf("GetByID().Label = %q, want %q", got.Label, "my-token")
	}
}

func TestStoreAndListCredentials(t *testing.T) {
	svc, ps := setupTestDB(t)
	providerID := createTestProvider(t, ps)

	if err := svc.SetupMasterPassword("master-pass"); err != nil {
		t.Fatalf("SetupMasterPassword() error: %v", err)
	}

	secrets := []string{"token-one", "token-two", "token-three"}
	for i, secret := range secrets {
		cred := &models.Credential{
			ProviderID: providerID,
			Label:      "cred-" + secret,
			AuthType:   "token",
			AuthData:   secret,
		}

		if err := svc.Store(cred); err != nil {
			t.Fatalf("Store() credential %d error: %v", i, err)
		}
	}

	creds, err := svc.List()
	if err != nil {
		t.Fatalf("List() error: %v", err)
	}

	if len(creds) != 3 {
		t.Fatalf("List() returned %d credentials, want 3", len(creds))
	}

	for i, c := range creds {
		if c.AuthData != secrets[i] {
			t.Fatalf("List()[%d].AuthData = %q, want %q", i, c.AuthData, secrets[i])
		}
	}
}

func TestGetByProviderID(t *testing.T) {
	svc, ps := setupTestDB(t)
	providerID := createTestProvider(t, ps)

	if err := svc.SetupMasterPassword("master-pass"); err != nil {
		t.Fatalf("SetupMasterPassword() error: %v", err)
	}

	cred := &models.Credential{
		ProviderID: providerID,
		Label:      "provider-cred",
		AuthType:   "token",
		AuthData:   "secret-for-provider",
	}

	if err := svc.Store(cred); err != nil {
		t.Fatalf("Store() error: %v", err)
	}

	creds, err := svc.GetByProviderID(providerID)
	if err != nil {
		t.Fatalf("GetByProviderID() error: %v", err)
	}

	if len(creds) != 1 {
		t.Fatalf("GetByProviderID() returned %d credentials, want 1", len(creds))
	}

	if creds[0].AuthData != "secret-for-provider" {
		t.Fatalf("GetByProviderID()[0].AuthData = %q, want %q", creds[0].AuthData, "secret-for-provider")
	}
}

func TestUpdateCredential(t *testing.T) {
	svc, ps := setupTestDB(t)
	providerID := createTestProvider(t, ps)

	if err := svc.SetupMasterPassword("master-pass"); err != nil {
		t.Fatalf("SetupMasterPassword() error: %v", err)
	}

	cred := &models.Credential{
		ProviderID: providerID,
		Label:      "updatable",
		AuthType:   "token",
		AuthData:   "original-secret",
	}

	if err := svc.Store(cred); err != nil {
		t.Fatalf("Store() error: %v", err)
	}

	// Update with new value
	cred.AuthData = "updated-secret"

	if err := svc.Update(cred); err != nil {
		t.Fatalf("Update() error: %v", err)
	}

	got, err := svc.GetByID(cred.ID)
	if err != nil {
		t.Fatalf("GetByID() after update error: %v", err)
	}

	if got.AuthData != "updated-secret" {
		t.Fatalf("GetByID().AuthData after update = %q, want %q", got.AuthData, "updated-secret")
	}
}

func TestDeleteCredential(t *testing.T) {
	svc, ps := setupTestDB(t)
	providerID := createTestProvider(t, ps)

	if err := svc.SetupMasterPassword("master-pass"); err != nil {
		t.Fatalf("SetupMasterPassword() error: %v", err)
	}

	cred := &models.Credential{
		ProviderID: providerID,
		Label:      "deletable",
		AuthType:   "token",
		AuthData:   "delete-me",
	}

	if err := svc.Store(cred); err != nil {
		t.Fatalf("Store() error: %v", err)
	}

	if err := svc.Delete(cred.ID); err != nil {
		t.Fatalf("Delete() error: %v", err)
	}

	_, err := svc.GetByID(cred.ID)
	if err == nil {
		t.Fatal("GetByID() after Delete() should return error")
	}
}

func TestOperationsWhileLocked(t *testing.T) {
	svc, ps := setupTestDB(t)
	providerID := createTestProvider(t, ps)

	if err := svc.SetupMasterPassword("master-pass"); err != nil {
		t.Fatalf("SetupMasterPassword() error: %v", err)
	}

	svc.Lock()

	cred := &models.Credential{
		ProviderID: providerID,
		Label:      "locked-cred",
		AuthType:   "token",
		AuthData:   "should-fail",
	}

	if err := svc.Store(cred); err != service.ErrLocked {
		t.Fatalf("Store() while locked: error = %v, want ErrLocked", err)
	}

	if _, err := svc.GetByID(1); err != service.ErrLocked {
		t.Fatalf("GetByID() while locked: error = %v, want ErrLocked", err)
	}

	if _, err := svc.GetByProviderID(providerID); err != service.ErrLocked {
		t.Fatalf("GetByProviderID() while locked: error = %v, want ErrLocked", err)
	}

	if _, err := svc.List(); err != service.ErrLocked {
		t.Fatalf("List() while locked: error = %v, want ErrLocked", err)
	}

	if err := svc.Update(cred); err != service.ErrLocked {
		t.Fatalf("Update() while locked: error = %v, want ErrLocked", err)
	}
}

func TestChangeMasterPassword(t *testing.T) {
	svc, ps := setupTestDB(t)
	providerID := createTestProvider(t, ps)

	if err := svc.SetupMasterPassword("old-password"); err != nil {
		t.Fatalf("SetupMasterPassword() error: %v", err)
	}

	// Store credentials before changing password
	secrets := []string{"secret-alpha", "secret-beta", "secret-gamma"}
	for _, secret := range secrets {
		cred := &models.Credential{
			ProviderID: providerID,
			Label:      "cred-" + secret,
			AuthType:   "token",
			AuthData:   secret,
		}

		if err := svc.Store(cred); err != nil {
			t.Fatalf("Store() error: %v", err)
		}
	}

	// Change master password
	if err := svc.ChangeMasterPassword("old-password", "new-password"); err != nil {
		t.Fatalf("ChangeMasterPassword() error: %v", err)
	}

	// Verify all credentials are still accessible with new key
	creds, err := svc.List()
	if err != nil {
		t.Fatalf("List() after password change error: %v", err)
	}

	if len(creds) != 3 {
		t.Fatalf("List() returned %d credentials after password change, want 3", len(creds))
	}

	for i, c := range creds {
		if c.AuthData != secrets[i] {
			t.Fatalf("List()[%d].AuthData after password change = %q, want %q", i, c.AuthData, secrets[i])
		}
	}

	// Lock and verify new password works
	svc.Lock()

	if err := svc.Unlock("new-password"); err != nil {
		t.Fatalf("Unlock() with new password error: %v", err)
	}

	// Old password should not work
	svc.Lock()

	if err := svc.Unlock("old-password"); err != service.ErrInvalidPassword {
		t.Fatalf("Unlock() with old password: error = %v, want ErrInvalidPassword", err)
	}
}

func TestChangeMasterPasswordWrongOldPassword(t *testing.T) {
	svc, _ := setupTestDB(t)

	if err := svc.SetupMasterPassword("correct-password"); err != nil {
		t.Fatalf("SetupMasterPassword() error: %v", err)
	}

	err := svc.ChangeMasterPassword("wrong-password", "new-password")
	if err == nil {
		t.Fatal("ChangeMasterPassword() with wrong old password should return error")
	}

	if err != service.ErrPasswordsDontMatch {
		t.Fatalf("ChangeMasterPassword() error = %v, want ErrPasswordsDontMatch", err)
	}
}

func TestLockUnlockCycle(t *testing.T) {
	svc, ps := setupTestDB(t)
	providerID := createTestProvider(t, ps)

	if err := svc.SetupMasterPassword("cycle-password"); err != nil {
		t.Fatalf("SetupMasterPassword() error: %v", err)
	}

	// Store a credential while unlocked
	cred := &models.Credential{
		ProviderID: providerID,
		Label:      "cycle-cred",
		AuthType:   "token",
		AuthData:   "cycle-secret",
	}

	if err := svc.Store(cred); err != nil {
		t.Fatalf("Store() error: %v", err)
	}

	// Lock, unlock, and verify credential is still accessible
	svc.Lock()

	if err := svc.Unlock("cycle-password"); err != nil {
		t.Fatalf("Unlock() error: %v", err)
	}

	got, err := svc.GetByID(cred.ID)
	if err != nil {
		t.Fatalf("GetByID() after lock/unlock cycle error: %v", err)
	}

	if got.AuthData != "cycle-secret" {
		t.Fatalf("GetByID().AuthData after lock/unlock = %q, want %q", got.AuthData, "cycle-secret")
	}
}
