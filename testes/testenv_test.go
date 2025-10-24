package testes

import "os"

func mustSetEnvDefault(key, def string) {
	if os.Getenv(key) == "" {
		_ = os.Setenv(key, def)
	}
}
