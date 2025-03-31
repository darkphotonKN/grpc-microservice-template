package commonenv

import "os"

func EnvString(key, fallback string) string {

	env := os.Getenv(key)

	if env == "" {
		return fallback
	}

	return env
}
