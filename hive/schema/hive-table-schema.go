package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/beltran/gohive"
	"github.com/rea1shane/gooooo/log"
	"github.com/sirupsen/logrus"
	"os"
)

var (
	logger     *logrus.Logger
	connection *gohive.Connection
	counter    = 1
)

func initHive() (err error) {
	hiveConfig := gohive.NewConnectConfiguration()
	hiveConfig.Username = "hive"
	hiveConfig.Password = "hive"
	connection, err = gohive.Connect("localhost", 10000, "NONE", hiveConfig)
	if err != nil {
		return err
	}
	return
}

// closeHive 关闭 hive
func closeHive() {
	connection.Close()
}

func main() {
	// logger
	logger = log.NewLogger()
	logFile, err := os.OpenFile("schema.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logger.Fatal(err)
	}
	logger.SetOutput(logFile)
	defer logFile.Close()

	// start
	err = initHive()
	if err != nil {
		logger.Fatal(err)
	}
	defer closeHive()
	schemas := getDbs()
	schemasJson, err := json.Marshal(schemas)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Info(string(schemasJson))
}

func getDbs() map[string]map[string]string {
	schemas := make(map[string]map[string]string)
	cursor := connection.Cursor()
	defer cursor.Close()
	cursor.Exec(context.Background(), "SHOW DATABASES")
	for cursor.HasMore(context.Background()) {
		var db string
		cursor.FetchOne(context.Background(), &db)
		schemas[db] = getTables(db)
	}
	return schemas
}

func getTables(db string) map[string]string {
	schemas := make(map[string]string)
	cursor := connection.Cursor()
	defer cursor.Close()
	cursor.Exec(context.Background(), fmt.Sprintf("USE %s", db))
	cursor.Exec(context.Background(), "SHOW TABLES")
	for cursor.HasMore(context.Background()) {
		var table string
		cursor.FetchOne(context.Background(), &table)
		schemas[table] = getSchema(db, table)
	}
	return schemas
}

func getSchema(db, table string) string {
	cursor := connection.Cursor()
	defer cursor.Close()
	cursor.Exec(context.Background(), fmt.Sprintf("DESC FORMATTED %s.%s", db, table))
	logger.Info(fmt.Sprintf("%d: %s.%s", counter, db, table))
	counter++
	for cursor.HasMore(context.Background()) {
		var colName, dataType, comment string
		cursor.FetchOne(context.Background(), &colName, &dataType, &comment)
		if cursor.Err != nil {
			logger.Error(cursor.Err)
			return ""
		}
		if colName == "SerDe Library:      " {
			return dataType
		}
	}
	return ""
}
