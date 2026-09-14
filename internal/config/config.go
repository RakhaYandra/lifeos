package config

import "os"

type Config struct {
	Port        string
	DBPath      string
	JWTSecret   string
	FrontendURL string
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func Load() Config {
	return Config{
		Port:        getenv("PORT", "8080"),
		DBPath:      getenv("DB_PATH", "lifeos.db"),
		JWTSecret:   getenv("JWT_SECRET", "dev-secret-ganti-di-prod"),
		FrontendURL: getenv("FRONTEND_URL", "http://localhost:5173"),
	}
}
