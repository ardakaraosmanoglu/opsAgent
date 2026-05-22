package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/meysam81/go-auth/auth/jwt"
	"github.com/meysam81/go-auth/middleware"
	"github.com/meysam81/go-auth/storage"
	"github.com/opsagent/opsagent/internal/config"
	"github.com/opsagent/opsagent/internal/infrastructure/storage/sqlite"
	"golang.org/x/crypto/bcrypt"
)

// JWTAuth handles JWT authentication for opsAgent.
type JWTAuth struct {
	tokenManager *jwt.TokenManager
	userStore   *SQLiteUserStore
	credStore   *SQLiteCredentialStore
	store       *sqlite.Store
	config      *config.Config
	middleware  *middleware.JWTMiddleware
}

// NewJWTAuth creates a new JWT authentication handler.
func NewJWTAuth(cfg *config.Config, store *sqlite.Store) (*JWTAuth, error) {
	userStore := NewSQLiteUserStore(store)
	credStore := NewSQLiteCredentialStore(store)

	// Signing key - in production load from config or settings
	// For MVP, use a fixed key that persists across restarts
	signingKey := []byte("opsagent-jwt-signing-key-32bytes!")

	tokenManager, err := jwt.NewTokenManager(jwt.Config{
		UserStore:       userStore,
		TokenStore:      nil, // We'll handle refresh tokens differently for now
		SigningKey:      signingKey,
		AccessTokenTTL:  15 * time.Minute,
		RefreshTokenTTL: 7 * 24 * time.Hour,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create token manager: %w", err)
	}

	jwtMiddleware := middleware.NewJWTMiddleware(middleware.JWTConfig{
		TokenManager: tokenManager,
	})

	return &JWTAuth{
		tokenManager: tokenManager,
		userStore:   userStore,
		credStore:   credStore,
		store:       store,
		config:      cfg,
		middleware:  jwtMiddleware,
	}, nil
}

// Authenticator interface for login
type Authenticator interface {
	Authenticate(ctx context.Context, username, password string) (*storage.User, error)
}

// Authenticate verifies username/password and returns the user.
func (a *JWTAuth) Authenticate(ctx context.Context, username, password string) (*storage.User, error) {
	u, err := a.userStore.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	// Get password hash from store
	hash, err := a.credStore.GetPasswordHash(ctx, u.ID)
	if err != nil {
		return nil, err
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword(hash, []byte(password)); err != nil {
		return nil, ErrNotFound
	}

	return u, nil
}

// LoginHandler returns a handler for login requests.
func (a *JWTAuth) LoginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		if req.Username == "" || req.Password == "" {
			http.Error(w, "username and password required", http.StatusBadRequest)
			return
		}

		user, err := a.Authenticate(r.Context(), req.Username, req.Password)
		if err != nil {
			// Log failed attempt
			a.store.InsertAuditLog(r.Context(), "login_failed", req.Username, "Failed login attempt", nil)
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}

		// Generate token pair
		tokenPair, err := a.tokenManager.GenerateTokenPair(r.Context(), user)
		if err != nil {
			http.Error(w, "failed to generate tokens", http.StatusInternalServerError)
			return
		}

		// Log successful login
		a.store.InsertAuditLog(r.Context(), "login", user.Username, "User logged in", nil)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"access_token":  tokenPair.AccessToken,
			"refresh_token": tokenPair.RefreshToken,
			"token_type":   tokenPair.TokenType,
			"expires_in":   tokenPair.ExpiresIn,
			"user": map[string]string{
				"id":       user.ID,
				"username": user.Username,
				"email":    user.Email,
				"name":     user.Name,
			},
		})
	}
}

// RefreshHandler returns a handler for token refresh requests.
func (a *JWTAuth) RefreshHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			RefreshToken string `json:"refresh_token"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		if req.RefreshToken == "" {
			http.Error(w, "refresh_token required", http.StatusBadRequest)
			return
		}

		tokenPair, err := a.tokenManager.RefreshAccessToken(r.Context(), req.RefreshToken)
		if err != nil {
			http.Error(w, "invalid or expired refresh token", http.StatusUnauthorized)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"access_token":  tokenPair.AccessToken,
			"refresh_token": req.RefreshToken, // Return same refresh token
			"token_type":   tokenPair.TokenType,
			"expires_in":   tokenPair.ExpiresIn,
		})
	}
}

// LogoutHandler returns a handler for logout requests.
func (a *JWTAuth) LogoutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Get user from context (set by JWT middleware)
		claims, ok := middleware.GetClaims(r)
		if ok {
			a.store.InsertAuditLog(r.Context(), "logout", claims.UserID, "User logged out", nil)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{"success": true})
	}
}

// MeHandler returns a handler that returns the current user's info.
func (a *JWTAuth) MeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := middleware.GetClaims(r)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		user, err := a.userStore.GetUserByID(r.Context(), claims.UserID)
		if err != nil || user == nil {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"name":     user.Name,
		})
	}
}

// Middleware returns the JWT authentication middleware.
func (a *JWTAuth) Middleware() func(http.Handler) http.Handler {
	return a.middleware.Middleware
}

// ValidateToken validates a token and returns claims.
func (a *JWTAuth) ValidateToken(ctx context.Context, tokenString string) (*jwt.Claims, error) {
	return a.tokenManager.ValidateToken(ctx, tokenString)
}

// GetSigningKey returns the signing key for testing purposes.
func (a *JWTAuth) GetSigningKey() []byte {
	return []byte("opsagent-jwt-signing-key-32bytes!")
}
