package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"github.com/BagasDhitya/owasp-secure-todo/internal/models"
	"github.com/BagasDhitya/owasp-secure-todo/internal/repo"
	"github.com/BagasDhitya/owasp-secure-todo/internal/security"
	"github.com/BagasDhitya/owasp-secure-todo/internal/validators"
)

type AuthHandler struct {
	Log            *zap.Logger
	Users          *repo.UserRepo
	DBRefreshStore RefreshStore // see below
	AccessSecret   []byte
	RefreshSecret  []byte
	AccessTTL      time.Duration
	RefreshTTL     time.Duration
	CSRFCookieTTL  time.Duration
	CSRFSecret     string
	BcryptCost     int
}

type RefreshStore interface {
	Save(ctx context.Context, userID int64, jti string, rawToken string, exp time.Time, ua, ip string) error
	Revoke(ctx context.Context, jti string) error
	FindValid(ctx context.Context, userID int64, jti, rawToken string) (bool, error)
}

// Register
func (h *AuthHandler) Register(c *gin.Context) {
	var dto validators.RegisterDTO
	if err := c.ShouldBindJSON(&dto); err != nil || validators.V.Struct(dto) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(dto.Password), h.BcryptCost)
	u := &models.User{Username: dto.Username, Email: dto.Email, PasswordHash: string(hash)}
	if err := h.Users.Create(c, u); err != nil {
		h.Log.Warn("register failed", zap.Error(err))
		c.JSON(http.StatusConflict, gin.H{"error": "username or email already exists"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "registered"})
}

// Login
func (h *AuthHandler) Login(c *gin.Context) {
	var dto validators.LoginDTO
	if err := c.ShouldBindJSON(&dto); err != nil || validators.V.Struct(dto) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid credentials"})
		return
	}
	u, _ := h.Users.ByEmail(c, dto.Email)
	if u == nil || bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(dto.Password)) != nil {
		h.Log.Info("failed login", zap.String("email", dto.Email))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	issueTokens(h, c, u)
}

// Refresh
func (h *AuthHandler) Refresh(c *gin.Context) {
	refreshCookie, err := c.Request.Cookie("refresh_token")
	if err != nil || refreshCookie.Value == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing refresh"})
		return
	}
	// parse jwt
	claims := &security.Claims{}
	tok, err := jwt.ParseWithClaims(refreshCookie.Value, claims, func(t *jwt.Token) (interface{}, error) {
		return h.RefreshSecret, nil
	}, jwt.WithValidMethods([]string{"HS256"}))
	if err != nil || !tok.Valid || claims.ExpiresAt == nil || time.Now().After(claims.ExpiresAt.Time) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh"})
		return
	}
	ok, _ := h.DBRefreshStore.FindValid(c, claims.UserID, claims.ID, refreshCookie.Value)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh"})
		return
	}
	// rotate: revoke old & issue new
	_ = h.DBRefreshStore.Revoke(c, claims.ID)
	u := &models.User{ID: claims.UserID, Email: claims.Email}
	issueTokens(h, c, u)
}

// Logout
func (h *AuthHandler) Logout(c *gin.Context) {
	refreshCookie, _ := c.Request.Cookie("refresh_token")
	if refreshCookie != nil && refreshCookie.Value != "" {
		// best effort revoke
		claims := &security.Claims{}
		if tok, err := jwt.ParseWithClaims(refreshCookie.Value, claims, func(t *jwt.Token) (interface{}, error) {
			return h.RefreshSecret, nil
		}, jwt.WithValidMethods([]string{"HS256"})); err == nil && tok.Valid {
			_ = h.DBRefreshStore.Revoke(c, claims.ID)
		}
	}
	clearAuthCookies(c)
	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}

func issueTokens(h *AuthHandler, c *gin.Context, u *models.User) {
	now := time.Now()
	accessExp := now.Add(h.AccessTTL)
	refreshExp := now.Add(h.RefreshTTL)
	accessClaims := security.Claims{
		UserID: u.ID, Email: u.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(accessExp),
		},
	}
	refreshClaims := accessClaims
	refreshClaims.ID = security.RandJTI() // implement helper
	refreshClaims.ExpiresAt = jwt.NewNumericDate(refreshExp)

	acc, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(h.AccessSecret)
	ref, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(h.RefreshSecret)

	// persist hashed refresh
	ua := c.Request.UserAgent()
	ip := c.ClientIP()
	_ = h.DBRefreshStore.Save(c, u.ID, refreshClaims.ID, ref, refreshExp, ua, ip)

	// httpOnly secure cookies
	c.SetSameSite(http.SameSiteStrictMode)
	setCookie(c, "access_token", acc, accessExp, true)
	setCookie(c, "refresh_token", ref, refreshExp, true)

	// CSRF cookie (non-HttpOnly)
	csrfVal, csrfExp := "csrf_"+security.RandJTI(), accessExp // simple random token per session
	setCookie(c, "csrf", csrfVal, csrfExp, false)

	h.Log.Info("login success", zap.Int64("uid", u.ID), zap.String("ip", ip))
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func setCookie(c *gin.Context, name, val string, exp time.Time, httpOnly bool) {
	secure := c.Request.TLS != nil // set true behind TLS (prod should be true)
	c.SetCookie(name, val, int(time.Until(exp).Seconds()), "/", "", secure, httpOnly)
}

func clearAuthCookies(c *gin.Context) {
	for _, n := range []string{"access_token", "refresh_token", "csrf"} {
		c.SetCookie(n, "", -1, "/", "", true, n != "csrf")
	}
}
