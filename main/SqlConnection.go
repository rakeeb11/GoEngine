package main

import (
	_ "github.com/go-sql-driver/mysql"
	"os"
	"log"
	"net/http"
	"database/sql"
	"fmt"
	"bytes"
)

func main() {
	Init()
}

const (
	SQL_NAME = "CLOUDSQL_CONNECTION_NAME"
	SQL_USER = "CLOUDSQL_USER"
	SQL_PWD = "CLOUDSQL_PASSWORD"
)

func Init()  {
	http.HandleFunc("/", handler)
}

func handler(w http.ResponseWriter, r *http.Request)  {

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	connectionName := mustGetEnv(SQL_NAME)
	user := mustGetEnv(SQL_USER)
	password := mustGetEnv(SQL_PWD)

	w.Header().Set("Content-Type", "text/plain")

	// connect to db
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@cloudsql(%s)/", user, password, connectionName))
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not open db: %v", err), 500)
	}
	defer db.Close()

	// query all database
	rows, err := db.Query("SHOW DATABASES")
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not query db: %v", err), 500)
		return
	}
	defer rows.Close()

	// print database
	buf := bytes.NewBufferString("Databases:\n")
	for rows.Next() {
		var dbName string
		if err := rows.Scan(&dbName); err != nil {
			http.Error(w, fmt.Sprintf("Could not scan result: %v", err), 500)
			return
		}
		fmt.Fprintf(buf, "- %s\n", dbName)
	}
	w.Write(buf.Bytes())
}

func mustGetEnv(key string) (env string) {
	env = os.Getenv(key)
	if env == "" {
		log.Panicf("%s environment variable not set.", key)
	}
	return
}