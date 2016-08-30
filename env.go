package main

import (
	"log"
	"net/url"
	"os"
	"strings"

	_ "github.com/joho/godotenv/autoload"
)

// Handles loading configuration from environment.
//
// sseserver/admin also does some of this, but since that's an indepdendant pkg
// now, leave that be where it is.

// Redis stuff
func envRedis() (host, password string) {
	env := os.Getenv("REDIS_URL")
	if env == "" {
		env = "http://localhost:6379"
	}

	url, err := url.Parse(env)
	if err != nil {
		log.Fatal("Could not parse what you have in the $REDIS_URL variable. Dying.")
	}

	host = url.Host
	if url.User != nil {
		password, _ = url.User.Password()
	}

	return
}

// da port to run on
func envPort() string {
	env := os.Getenv("PORT")
	if env != "" {
		return (":" + env) //golang and its weird port mechanics
	}
	return ":8001"
}

// what is our dev/staging/prod environment
func env() string {
	env := os.Getenv("GO_ENV")
	if env != "" {
		return strings.ToLower(env)
	}
	return "development"
}
func envIsStaging() bool {
	return env() == "staging"
}
func envIsProduction() bool {
	return env() == "production"
}
