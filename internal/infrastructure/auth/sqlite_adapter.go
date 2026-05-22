package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/meysam81/go-auth/storage"
	"github.com/opsagent/opsagent/internal/domain"
	"github.com/opsagent/opsagent/internal/infrastructure/storage/sqlite"
)

// ErrNotFound is returned when a user or credential is not found.
var ErrNotFound = errors.New("not found")

// SQLiteUserStore adapts sqlite.Store to go-auth's UserStore interface.
type SQLiteUserStore struct {
	store *sqlite.Store
}

// NewSQLiteUserStore creates a new SQLite-backed user store.
func NewSQLiteUserStore(store *sqlite.Store) *SQLiteUserStore {
	return &SQLiteUserStore{store: store}
}

// storageToDomain converts go-auth User to domain.User
func storageToDomain(u *storage.User) *domain.User {
	if u == nil {
		return nil
	}
	return &domain.User{
		ID:           0, // will be set by GetUserByUsername
		Username:     u.Username,
		PasswordHash: "", // not in go-auth storage.User
		Role:         "admin",
		CreatedAt:    u.CreatedAt.Format(time.RFC3339),
	}
}

// CreateUser is not used in opsAgent - users are created via setup
func (s *SQLiteUserStore) CreateUser(ctx context.Context, user *storage.User) error {
	return errors.New("not implemented: use setup flow instead")
}

// GetUserByID retrieves a user by ID (converts string ID back to int64)
func (s *SQLiteUserStore) GetUserByID(ctx context.Context, id string) (*storage.User, error) {
	// Convert string ID to int64 (go-auth stores ID as string)
	var idInt int64
	_, err := fmt.Sscanf(id, "%d", &idInt)
	if err != nil {
		return nil, ErrNotFound
	}

	u, err := s.store.GetUserByID(ctx, idInt)
	if err != nil || u == nil {
		return nil, ErrNotFound
	}
	return &storage.User{
		ID:            fmt.Sprintf("%d", u.ID),
		Username:      u.Username,
		Email:         u.Username + "@localhost",
		Name:          u.Username,
		Provider:      "local",
		EmailVerified: true,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}, nil
}

// GetUserByEmail retrieves a user by email (not supported in opsAgent)
func (s *SQLiteUserStore) GetUserByEmail(ctx context.Context, email string) (*storage.User, error) {
	return nil, ErrNotFound
}

// GetUserByUsername retrieves a user by username
func (s *SQLiteUserStore) GetUserByUsername(ctx context.Context, username string) (*storage.User, error) {
	u, err := s.store.GetUserByUsername(ctx, username)
	if err != nil || u == nil {
		return nil, ErrNotFound
	}
	return &storage.User{
		ID:            fmt.Sprintf("%d", u.ID),
		Username:      u.Username,
		Email:         u.Username + "@localhost",
		Name:          u.Username,
		Provider:      "local",
		EmailVerified: true,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}, nil
}

// UpdateUser updates user info (not used)
func (s *SQLiteUserStore) UpdateUser(ctx context.Context, user *storage.User) error {
	return errors.New("not implemented")
}

// DeleteUser removes a user (not used)
func (s *SQLiteUserStore) DeleteUser(ctx context.Context, id string) error {
	return errors.New("not implemented")
}

// SQLiteCredentialStore adapts sqlite.Store to go-auth's CredentialStore interface.
type SQLiteCredentialStore struct {
	store *sqlite.Store
}

// NewSQLiteCredentialStore creates a new SQLite-backed credential store.
func NewSQLiteCredentialStore(store *sqlite.Store) *SQLiteCredentialStore {
	return &SQLiteCredentialStore{store: store}
}

// StorePasswordHash stores a password hash for a user
func (s *SQLiteCredentialStore) StorePasswordHash(ctx context.Context, userID string, hash []byte) error {
	// Not used - password hashes are stored during user creation
	return nil
}

// GetPasswordHash retrieves the password hash for a user
func (s *SQLiteCredentialStore) GetPasswordHash(ctx context.Context, userID string) ([]byte, error) {
	// Convert string ID back to int64
	var userIDInt int64
	_, err := fmt.Sscanf(userID, "%d", &userIDInt)
	if err != nil {
		return nil, ErrNotFound
	}

	u, err := s.store.GetUserByID(ctx, userIDInt)
	if err != nil || u == nil {
		return nil, ErrNotFound
	}
	return []byte(u.PasswordHash), nil
}

// StoreWebAuthnCredential is not implemented
func (s *SQLiteCredentialStore) StoreWebAuthnCredential(ctx context.Context, userID string, credential *storage.WebAuthnCredential) error {
	return errors.New("not implemented")
}

// GetWebAuthnCredentials is not implemented
func (s *SQLiteCredentialStore) GetWebAuthnCredentials(ctx context.Context, userID string) ([]*storage.WebAuthnCredential, error) {
	return nil, errors.New("not implemented")
}

// UpdateWebAuthnCredential is not implemented
func (s *SQLiteCredentialStore) UpdateWebAuthnCredential(ctx context.Context, credential *storage.WebAuthnCredential) error {
	return errors.New("not implemented")
}

// DeleteWebAuthnCredential is not implemented
func (s *SQLiteCredentialStore) DeleteWebAuthnCredential(ctx context.Context, credentialID []byte) error {
	return errors.New("not implemented")
}

// StorePasswordResetToken is not implemented
func (s *SQLiteCredentialStore) StorePasswordResetToken(ctx context.Context, userID string, token string, expiresAt time.Time) error {
	return errors.New("not implemented")
}

// ValidatePasswordResetToken is not implemented
func (s *SQLiteCredentialStore) ValidatePasswordResetToken(ctx context.Context, token string) (string, error) {
	return "", errors.New("not implemented")
}

// DeletePasswordResetToken is not implemented
func (s *SQLiteCredentialStore) DeletePasswordResetToken(ctx context.Context, token string) error {
	return errors.New("not implemented")
}

// StoreEmailVerificationToken is not implemented
func (s *SQLiteCredentialStore) StoreEmailVerificationToken(ctx context.Context, userID string, token string, expiresAt time.Time) error {
	return errors.New("not implemented")
}

// ValidateEmailVerificationToken is not implemented
func (s *SQLiteCredentialStore) ValidateEmailVerificationToken(ctx context.Context, token string) (string, error) {
	return "", errors.New("not implemented")
}

// DeleteEmailVerificationToken is not implemented
func (s *SQLiteCredentialStore) DeleteEmailVerificationToken(ctx context.Context, token string) error {
	return errors.New("not implemented")
}

// StoreTOTPSecret is not implemented
func (s *SQLiteCredentialStore) StoreTOTPSecret(ctx context.Context, userID string, secret string, backupCodes []string) error {
	return errors.New("not implemented")
}

// GetTOTPSecret is not implemented
func (s *SQLiteCredentialStore) GetTOTPSecret(ctx context.Context, userID string) (string, []string, error) {
	return "", nil, errors.New("not implemented")
}

// DeleteTOTPSecret is not implemented
func (s *SQLiteCredentialStore) DeleteTOTPSecret(ctx context.Context, userID string) error {
	return errors.New("not implemented")
}

// UseBackupCode is not implemented
func (s *SQLiteCredentialStore) UseBackupCode(ctx context.Context, userID string, code string) error {
	return errors.New("not implemented")
}
