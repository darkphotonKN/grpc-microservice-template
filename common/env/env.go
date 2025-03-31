package commonenv

import (
	"fmt"
	"os"
)

func EnvString(key, fallback string) string {

	env := os.Getenv(key)

	fmt.Printf("env key %s retrieved env: %s\n", key, env)

	if env == "" {
		return fallback
	}

	return env
}
