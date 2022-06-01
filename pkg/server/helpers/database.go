package helpers

import (
	"fmt"
	"github.com/Datariah/psychic-guacamole/internal"
	"github.com/Datariah/psychic-guacamole/internal/secrets"
	"github.com/onrik/gorm-logrus"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func GetDSN(username, password, dbName, host, port string) string {
	defaultDbName := "psychic-guacamole"
	defaultHost := "127.0.0.1"
	defaultPort := "3306"

	if username == "" {
		log.Fatal("database username cannot be empty")
	}

	if password == "" {
		log.Fatal("database password cannot be empty")
	}

	if dbName == "" {
		log.Warnf("database name not specified, defaulting to \"%s\"", defaultDbName)
		dbName = defaultDbName
	}

	if host == "" {
		log.Warnf("database host not specified, defaulting to \"%s\"", defaultHost)
		host = defaultHost
	}

	if port == "" {
		log.Warnf("database port not specified, defaulting to \"%s\"", defaultPort)
		port = defaultPort
	}

	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		username, password, host, port, dbName)
}

func InitializeDBConnection() *gorm.DB {
	dbConfig, err := secrets.GetSecretValues("psychic-guacamole-rds", internal.AwsRegion)
	if err != nil {
		log.Panicf("error while retrieving DSN config from SecretsManager: %v", err)
	}

	dsn := GetDSN(
		(*dbConfig)["username"].(string), (*dbConfig)["password"].(string), (*dbConfig)["dbname"].(string), (*dbConfig)["host"].(string), fmt.Sprintf("%.0f", (*dbConfig)["port"]),
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gorm_logrus.New(),
	})
	if err != nil {
		log.Panicf("error while opening database connection: %v", err)
	}

	return db
}
