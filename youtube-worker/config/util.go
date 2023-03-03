package config

import "os"

func getEnvOr(key string, def string) string {
	env, exist := os.LookupEnv(key)
	if exist {
		return env
	} else {
		return def
	}
}
