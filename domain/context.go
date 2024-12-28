package domain

import (
	"flag"
)

var ServerEnv string

const SERVER_ENV_LOCAL = "local"
const SERVER_ENV_TEST_1 = "test1"
const SERVER_ENV_PROD = "prod"

func FlagInit() {
	serverEnvFlag := flag.String("server_env", "test1", "development environment")
	flag.Parse()
	ServerEnv = *serverEnvFlag
}
