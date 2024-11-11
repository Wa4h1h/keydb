package config

import "github.com/Wa4h1h/memdb/internal/utils"

type ServerConfig struct {
	Port                        string
	LogLevel                    string
	TCPReadTimeout              uint
	BackOffLimit                int
	TTLBackgroundWorkerInterval int
}

func LoadServerConfig() *ServerConfig {
	srvConfig := new(ServerConfig)

	srvConfig.Port = utils.GetEnv[string]("8000", false, "PORT")
	srvConfig.LogLevel = utils.GetEnv[string]("debug", false, "LOG_LEVEL")
	srvConfig.TCPReadTimeout = utils.GetEnv[uint]("10", false, "TCP_READ_TIMEOUT")
	srvConfig.BackOffLimit = utils.GetEnv[int]("5", false, "Back_OFF_Limit")
	srvConfig.TTLBackgroundWorkerInterval = utils.GetEnv[int]("5",
		false, "TTL_BACKGROUNO_WORKDER_INTERVAL")

	return srvConfig
}
