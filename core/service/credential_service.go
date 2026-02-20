package service

import (
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"sync"

	"GitSyncer/core/crypto"
	"GitSyncer/core/models"
	"GitSyncer/core/store"
)

const (
	settingMasterPasswordHash = "master_password_hash"
	settingMasterPasswordSalt = "master_password_salt"
)

var (
	ErrLocked             = errors.New("credential service is locked")
	ErrAlreadySetup       = errors.New("master password is already configured")
	ErrNotSetup           = errors.New("master password is not configured")
	ErrInvalidPassword    = errors.New("invalid master password")
	ErrPasswordsDontMatch = errors.New("old password is incorrect")
)

// CredentialService provides encryption-aware credential operations.
type CredentialService struct {
	db            *sql.DB
	credStore     *store.CredentialStore
	settingStore  *store.SettingStore
	mu            sync.RWMutex
	derivedKey    []byte
	locked        bool
}

// NewCredentialService creates a new CredentialService.
func NewCredentialService(db *sql.DB, credStore *store.CredentialStore, settingStore *store.SettingStore) *CredentialService {
	return &CredentialService{
		db:           db,
		credStore:    credStore,
		settingStore: settingStore,
		locked:       true,
	}
}

// IsSetup checks whether a master password has been configured.
func (s *CredentialService) IsSetup() (bool, error) {
	return s.settingStore.Exists(settingMasterPasswordHash)
}

// SetupMasterPassword configures the master password for first-time use.
func (s *CredentialService) SetupMasterPassword(password string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	setup, err := s.IsSetup()
	if err != nil {
		return err
	}

	if setup {
		return ErrAlreadySetup
	}

	salt, err := crypto.GenerateSalt()
	if err != nil {
		return fmt.Errorf("CredentialService.SetupMasterPassword: %w", err)
	}

	hash := crypto.HashPassword(password, salt)

	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("CredentialService.SetupMasterPassword: begin tx: %w", err)
	}

	defer tx.Rollback() //nolint:errcheck

	if err := s.settingStore.SetInTx(tx, settingMasterPasswordHash, base64.StdEncoding.EncodeToString(hash)); err != nil {
		return fmt.Errorf("CredentialService.SetupMasterPassword: store hash: %w", err)
	}

	if err := s.settingStore.SetInTx(tx, settingMasterPasswordSalt, base64.StdEncoding.EncodeToString(salt)); err != nil {
		return fmt.Errorf("CredentialService.SetupMasterPassword: store salt: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("CredentialService.SetupMasterPassword: commit: %w", err)
	}

	s.derivedKey = crypto.DeriveKey(password, salt)
	s.locked = false

	return nil
}

// Unlock verifies the master password and derives the encryption key.
func (s *CredentialService) Unlock(password string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	salt, hash, err := s.loadMasterPasswordData()
	if err != nil {
		return err
	}

	if !crypto.VerifyPassword(password, salt, hash) {
		return ErrInvalidPassword
	}

	s.derivedKey = crypto.DeriveKey(password, salt)
	s.locked = false

	return nil
}

// Lock zeroes the encryption key from memory and marks the service as locked.
func (s *CredentialService) Lock() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.derivedKey != nil {
		crypto.ZeroBytes(s.derivedKey)
		s.derivedKey = nil
	}

	s.locked = true
}

// IsLocked returns whether the service is currently locked.
func (s *CredentialService) IsLocked() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.locked
}

// Store encrypts and stores a credential.
func (s *CredentialService) Store(cred *models.Credential) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.locked {
		return ErrLocked
	}

	salt, err := crypto.GenerateSalt()
	if err != nil {
		return fmt.Errorf("CredentialService.Store: %w", err)
	}

	key := crypto.DeriveKey(string(s.derivedKey), salt)
	defer crypto.ZeroBytes(key)

	ciphertext, err := crypto.Encrypt([]byte(cred.AuthData), key)
	if err != nil {
		return fmt.Errorf("CredentialService.Store: encrypt: %w", err)
	}

	cred.AuthData = base64.StdEncoding.EncodeToString(ciphertext)
	cred.Salt = salt

	return s.credStore.Create(cred)
}

// GetByID retrieves and decrypts a credential by its ID.
func (s *CredentialService) GetByID(id int64) (*models.Credential, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.locked {
		return nil, ErrLocked
	}

	cred, err := s.credStore.GetByID(id)
	if err != nil {
		return nil, err
	}

	if err := s.decryptCredential(cred); err != nil {
		return nil, fmt.Errorf("CredentialService.GetByID(%d): %w", id, err)
	}

	return cred, nil
}

// GetByProviderID retrieves and decrypts all credentials for a provider.
func (s *CredentialService) GetByProviderID(providerID int64) ([]models.Credential, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.locked {
		return nil, ErrLocked
	}

	creds, err := s.credStore.GetByProviderID(providerID)
	if err != nil {
		return nil, err
	}

	for i := range creds {
		if err := s.decryptCredential(&creds[i]); err != nil {
			return nil, fmt.Errorf("CredentialService.GetByProviderID(%d): credential %d: %w", providerID, creds[i].ID, err)
		}
	}

	return creds, nil
}

// List retrieves and decrypts all credentials.
func (s *CredentialService) List() ([]models.Credential, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.locked {
		return nil, ErrLocked
	}

	creds, err := s.credStore.List()
	if err != nil {
		return nil, err
	}

	for i := range creds {
		if err := s.decryptCredential(&creds[i]); err != nil {
			return nil, fmt.Errorf("CredentialService.List: credential %d: %w", creds[i].ID, err)
		}
	}

	return creds, nil
}

// Update re-encrypts and updates a credential.
func (s *CredentialService) Update(cred *models.Credential) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.locked {
		return ErrLocked
	}

	salt, err := crypto.GenerateSalt()
	if err != nil {
		return fmt.Errorf("CredentialService.Update: %w", err)
	}

	key := crypto.DeriveKey(string(s.derivedKey), salt)
	defer crypto.ZeroBytes(key)

	ciphertext, err := crypto.Encrypt([]byte(cred.AuthData), key)
	if err != nil {
		return fmt.Errorf("CredentialService.Update: encrypt: %w", err)
	}

	cred.AuthData = base64.StdEncoding.EncodeToString(ciphertext)
	cred.Salt = salt

	return s.credStore.Update(cred)
}

// Delete removes a credential by ID.
func (s *CredentialService) Delete(id int64) error {
	return s.credStore.Delete(id)
}

// ChangeMasterPassword re-encrypts all credentials with a new master password.
func (s *CredentialService) ChangeMasterPassword(oldPassword, newPassword string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Verify old password
	oldSalt, oldHash, err := s.loadMasterPasswordData()
	if err != nil {
		return err
	}

	if !crypto.VerifyPassword(oldPassword, oldSalt, oldHash) {
		return ErrPasswordsDontMatch
	}

	// Load all credentials (encrypted)
	creds, err := s.credStore.List()
	if err != nil {
		return fmt.Errorf("CredentialService.ChangeMasterPassword: list: %w", err)
	}

	// Decrypt all credentials with old key
	oldDerivedKey := s.derivedKey
	plaintexts := make([]string, len(creds))

	for i := range creds {
		key := crypto.DeriveKey(string(oldDerivedKey), creds[i].Salt)

		ciphertext, err := base64.StdEncoding.DecodeString(creds[i].AuthData)
		if err != nil {
			crypto.ZeroBytes(key)

			return fmt.Errorf("CredentialService.ChangeMasterPassword: decode credential %d: %w", creds[i].ID, err)
		}

		plaintext, err := crypto.Decrypt(ciphertext, key)
		crypto.ZeroBytes(key)

		if err != nil {
			return fmt.Errorf("CredentialService.ChangeMasterPassword: decrypt credential %d: %w", creds[i].ID, err)
		}

		plaintexts[i] = string(plaintext)
	}

	// Generate new master password hash
	newMasterSalt, err := crypto.GenerateSalt()
	if err != nil {
		return fmt.Errorf("CredentialService.ChangeMasterPassword: generate master salt: %w", err)
	}

	newMasterHash := crypto.HashPassword(newPassword, newMasterSalt)
	newDerivedKey := crypto.DeriveKey(newPassword, newMasterSalt)

	// Re-encrypt all credentials in a transaction
	tx, err := s.db.Begin()
	if err != nil {
		crypto.ZeroBytes(newDerivedKey)

		return fmt.Errorf("CredentialService.ChangeMasterPassword: begin tx: %w", err)
	}

	defer tx.Rollback() //nolint:errcheck

	// Update master password hash and salt via SettingStore
	if err := s.settingStore.SetInTx(tx, settingMasterPasswordHash, base64.StdEncoding.EncodeToString(newMasterHash)); err != nil {
		crypto.ZeroBytes(newDerivedKey)

		return fmt.Errorf("CredentialService.ChangeMasterPassword: update hash: %w", err)
	}

	if err := s.settingStore.SetInTx(tx, settingMasterPasswordSalt, base64.StdEncoding.EncodeToString(newMasterSalt)); err != nil {
		crypto.ZeroBytes(newDerivedKey)

		return fmt.Errorf("CredentialService.ChangeMasterPassword: update salt: %w", err)
	}

	// Re-encrypt each credential via CredentialStore
	for i := range creds {
		newSalt, err := crypto.GenerateSalt()
		if err != nil {
			crypto.ZeroBytes(newDerivedKey)

			return fmt.Errorf("CredentialService.ChangeMasterPassword: generate salt for credential %d: %w", creds[i].ID, err)
		}

		key := crypto.DeriveKey(string(newDerivedKey), newSalt)

		ciphertext, err := crypto.Encrypt([]byte(plaintexts[i]), key)
		crypto.ZeroBytes(key)

		if err != nil {
			crypto.ZeroBytes(newDerivedKey)

			return fmt.Errorf("CredentialService.ChangeMasterPassword: encrypt credential %d: %w", creds[i].ID, err)
		}

		if err := s.credStore.UpdateEncryptedInTx(tx, creds[i].ID, base64.StdEncoding.EncodeToString(ciphertext), newSalt); err != nil {
			crypto.ZeroBytes(newDerivedKey)

			return fmt.Errorf("CredentialService.ChangeMasterPassword: update credential %d: %w", creds[i].ID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		crypto.ZeroBytes(newDerivedKey)

		return fmt.Errorf("CredentialService.ChangeMasterPassword: commit: %w", err)
	}

	// Swap the derived key
	crypto.ZeroBytes(s.derivedKey)
	s.derivedKey = newDerivedKey

	return nil
}

// decryptCredential decrypts a credential's AuthData in place.
func (s *CredentialService) decryptCredential(cred *models.Credential) error {
	if cred.AuthData == "" {
		return nil
	}

	ciphertext, err := base64.StdEncoding.DecodeString(cred.AuthData)
	if err != nil {
		return fmt.Errorf("decode auth_data: %w", err)
	}

	key := crypto.DeriveKey(string(s.derivedKey), cred.Salt)
	defer crypto.ZeroBytes(key)

	plaintext, err := crypto.Decrypt(ciphertext, key)
	if err != nil {
		return fmt.Errorf("decrypt auth_data: %w", err)
	}

	cred.AuthData = string(plaintext)

	return nil
}

// loadMasterPasswordData reads the master password hash and salt from settings.
func (s *CredentialService) loadMasterPasswordData() (salt, hash []byte, err error) {
	hashB64, err := s.settingStore.Get(settingMasterPasswordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, ErrNotSetup
		}

		return nil, nil, fmt.Errorf("CredentialService: load master hash: %w", err)
	}

	saltB64, err := s.settingStore.Get(settingMasterPasswordSalt)
	if err != nil {
		return nil, nil, fmt.Errorf("CredentialService: load master salt: %w", err)
	}

	hash, err = base64.StdEncoding.DecodeString(hashB64)
	if err != nil {
		return nil, nil, fmt.Errorf("CredentialService: decode master hash: %w", err)
	}

	salt, err = base64.StdEncoding.DecodeString(saltB64)
	if err != nil {
		return nil, nil, fmt.Errorf("CredentialService: decode master salt: %w", err)
	}

	return salt, hash, nil
}
