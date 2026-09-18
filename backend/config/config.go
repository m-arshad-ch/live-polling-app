package config

import "os"

type Config struct {
	Port        string
	MongoURI    string
	MongoDB     string
	RedisURL    string
	JWTSecret   string
	FrontendURL string
}

func Load() Config {
	return Config{
		Port:        os.Getenv("PORT"),
		MongoURI:    get("MONGO_URI", "mongodb://localhost:27017"),
		MongoDB:     get("MONGO_DB", "livepoll"),
		RedisURL:    get("REDIS_URL", "redis://localhost:6379"),
		JWTSecret:   get("JWT_SECRET", "dev-secret-change-me"),
		FrontendURL: get("FRONTEND_URL", "http://localhost:5173"),
	}
}

func get(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
