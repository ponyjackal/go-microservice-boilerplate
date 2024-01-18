package config

import (
	"fmt"
	"os"

	"github.com/ponyjackal/go-microservice-boilerplate/pkg/logger"
)

type ServerConfiguration struct {
	Port                 string
	Secret               string
	LimitCountPerRequest int64
}

func ServerConfig() string {
	appServer := fmt.Sprintf("%s:%s", os.Getenv("SERVER_HOST"), os.Getenv("SERVER_PORT"))
	logger.Infof("Server Running at : %s", appServer)
	return appServer
}
