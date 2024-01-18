package config

import (
	"fmt"
	"os"
)

type DatabaseConfiguration struct {
	Driver   string
	Dbname   string
	Username string
	Password string
	Host     string
	Port     string
	LogMode  bool
}

func DbConfiguration() (string, string) {
	masterDBName := os.Getenv("MASTER_DB_NAME")
	masterDBUser := os.Getenv("MASTER_DB_USER")
	masterDBPassword := os.Getenv("MASTER_DB_PASSWORD")
	masterDBHost := os.Getenv("MASTER_DB_HOST")
	masterDBPort := os.Getenv("MASTER_DB_PORT")
	masterDBSslMode := os.Getenv("MASTER_SSL_MODE")

	replicaDBName := os.Getenv("REPLICA_DB_NAME")
	replicaDBUser := os.Getenv("REPLICA_DB_USER")
	replicaDBPassword := os.Getenv("REPLICA_DB_PASSWORD")
	replicaDBHost := os.Getenv("REPLICA_DB_HOST")
	replicaDBPort := os.Getenv("REPLICA_DB_PORT")
	replicaDBSslMode := os.Getenv("REPLICA_SSL_MODE")

	masterDBDSN := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		masterDBHost, masterDBUser, masterDBPassword, masterDBName, masterDBPort, masterDBSslMode,
	)

	replicaDBDSN := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		replicaDBHost, replicaDBUser, replicaDBPassword, replicaDBName, replicaDBPort, replicaDBSslMode,
	)
	return masterDBDSN, replicaDBDSN
}
