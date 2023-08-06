package sqlc

import (
	"database/sql"
	"log"
	"os"
	"simple_bank/db/util"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *Queries
var testdb *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil{
		log.Fatal("can't load config values",err)
	}

	testdb, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Connection to DB can't be established", err)
	}
	testQueries = New(testdb)

	os.Exit(m.Run())
}
