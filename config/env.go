package config

import "os"

// Env constants
const (
	Stage Env42 = "ENV_STAGE"
	Dev   Env42 = "ENV_DEV"
	Prod  Env42 = "ENV_PROD"
)

type Env42 string

var env Env42

func init() {
	switch os.Getenv("ENV_42") {
	case string(Stage):
		env = Stage
	case string(Prod):
		env = Prod
	default:
		env = Dev
	}
}

// CurrentEnv represents the current environment
func CurrentEnv() Env42 {
	return env
}
