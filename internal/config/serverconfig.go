package config

import "github.com/Wa4h1h/memdb/internal/utils"

type ServerConfig struct {
	Port           string
	LogLevel       string
	TCPReadTimeout uint
}

func LoadServerConfig() *ServerConfig {
	srvConfig := new(ServerConfig)

	srvConfig.Port = utils.GetEnv[string]("8000", false, "PORT")
	srvConfig.LogLevel = utils.GetEnv[string]("debug", false, "LOG_LEVEL")
	srvConfig.TCPReadTimeout = utils.GetEnv[uint]("1", false, "TCP_READ_TIMEOUT")

	return srvConfig
}
