package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv  string
	AppPort string

	DBURL string

	JWTAccessSecret  string
	JWTRefreshSecret string
	AccessTTL        time.Duration
	RefreshTTL       time.Duration

	CSRfSecret string
	BcryptCost int
}

func Load() *Config {
	_ = godotenv.Load() // loads .env if present (safe in prod: no-op)

	get := func(k string) string {
		v := os.Getenv(k)
		if v == "" {
			log.Fatalf("missing required env: %s", k)
		}
		return v
	}

	accMin, _ := strconv.Atoi(get("JWT_ACCESS_TTL_MIN"))
	refH, _ := strconv.Atoi(get("JWT_REFRESH_TTL_H"))
	bc, _ := strconv.Atoi(get("BCRYPT_COST"))

	return &Config{
		AppEnv:           os.Getenv("APP_ENV"),
		AppPort:          os.Getenv("APP_PORT"),
		DBURL:            get("DB_URL"),
		JWTAccessSecret:  get("JWT_ACCESS_SECRET"),
		JWTRefreshSecret: get("JWT_REFRESH_SECRET"),
		AccessTTL:        time.Duration(accMin) * time.Minute,
		RefreshTTL:       time.Duration(refH) * time.Hour,
		CSRfSecret:       get("CSRF_SECRET"),
		BcryptCost:       bc,
	}
}
