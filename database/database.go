package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"sort"
)

var Pool *sqlx.DB

func InitDatabase() {
	if Pool != nil {
		slog.Debug("Pool already initialized returning...")
		return
	}

	createDatabasePool()
	createTables()
}

func createDatabasePool() {
	if Pool != nil {
		return
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		viper.GetString("database.username"),
		viper.GetString("database.password"),
		viper.GetString("database.host"),
		viper.GetInt("database.port"),
		viper.GetString("database.name"),
		viper.GetString("database.sslmode"),
	)

	slog.Info(fmt.Sprintf("Database string %s", connStr))

	Pool = sqlx.MustConnect("postgres", connStr)
	Pool.SetMaxIdleConns(1)
	Pool.SetMaxOpenConns(1)
}

func createTables() {
	fileContent, err := os.ReadFile("/src/app/sql/init_database.sql")
	if err != nil {
		slog.Error(fmt.Sprintf("Error reading content of file: %s", err))
		panic(err)
	}
	executeSql(string(fileContent))
}

func executeSql(sqlCommand string) {
	tx := Pool.MustBegin()
	tx.MustExec(sqlCommand)
	err := tx.Commit()
	if err != nil {
		slog.Error(fmt.Sprintf("Error executing sql from file: %s", err))
		panic(err)
	}
}

func getCleanScripts(scriptsDirectory string) []string {
	// Get scripts to clean the data

	filesFound := []string{}
	entries, err := os.ReadDir(scriptsDirectory)
	if err != nil {
		slog.Error(fmt.Sprintf("Error reading directory: %s", err))
		panic(err)
	}

	re := regexp.MustCompile(`\d+\_\w+\.sql$`)

	for _, entry := range entries {
		entryName := entry.Name()

		if entry.IsDir() || !re.MatchString(entryName) {
			continue
		}
		filesFound = append(filesFound, entryName)
	}

	sort.Strings(filesFound)
	return filesFound

}

func CleanDataScripts() {
	scriptFiles := getCleanScripts("/src/app/sql/")
	for _, file := range scriptFiles {
		fileContent, err := os.ReadFile(filepath.Join("/src/app/sql/", file))
		if err != nil {
			slog.Error(fmt.Sprintf("Error reading content of file: %s", err))
			panic(err)
		}
		executeSql(string(fileContent))
	}
}
