package database

import (
	"database/sql"
	"fmt"

	// Import PQ for the sql package
	_ "github.com/lib/pq"

	"nvmtech.nl/homemanager/tools/log"
)

var Database *DB

// DB interface wrapper for sql/pq
type DB struct {
	// Credentials
	Username string
	Password string
	Name     string

	// Connection settings
	URL  string
	Port int

	// Internal settings
	connection *sql.DB
}

// NewDB is a constructor for the database
func NewDB(username, password, name, url string, port int) *DB {
	db := DB{
		username,
		password,
		name,
		url,
		port,
		nil,
	}

	return &db
}

// Connect opens a connection to the dabase
func (db *DB) Connect() {
	connectionString := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", db.Username, db.Password, db.Name)
	var err error

	// Connect
	db.connection, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal("Database", "Failed to connect to the postgres database")
	}

	// Ping the connection
	err = db.connection.Ping()
	if err != nil {
		log.Fatal("Database", "Failed to ping the database server")
	}

	log.Info("Database", "Succesfully connected to PostgresDB")
}

// CreateTable creates a single table in the database
func (db *DB) CreateTable(query string) error {
	_, err := db.connection.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

// CreateTables creates multiple tables based in the array of strings provided
func (db *DB) CreateTables(queries []string) error {
	tx, err := db.connection.Begin()
	if err != nil {
		log.Warn("Database", "Failed to create a transaction, "+err.Error())
	}

	for _, query := range queries {
		_, err := tx.Exec(query)
		if err != nil {
			log.Err("Database", err.Error())
			continue
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal("Database", err.Error())
	}

	return nil
}

// Exec is short for db.connection.Exec()
func (db *DB) Exec(query string) {
	db.connection.Exec(query)
}

// ExecBatch executes a query for every value in a BatchQuery struct
// it all happens in a transaction
func (db *DB) ExecBatch(batch BatchQuery) {
	tx, err := db.connection.Begin()
	if err != nil {
		log.Err("Database", err.Error())
	}
	for _, vars := range batch.Values {
		_, err = tx.Exec(batch.Query, vars...)
		if err != nil {
			log.Warn("Database", err.Error()+"\n"+batch.Query)
		}
	}
	err = tx.Commit()
	if err != nil {
		log.Err("Database", err.Error())
	}
}

// Query is short for db.connection.query
func (db *DB) Query(sql string, vals ...interface{}) *sql.Rows {
	rows, err := db.connection.Query(sql, vals...)
	if err != nil {
		log.Warn("Database", "Error while executing query "+err.Error())
	}

	return rows
}

// BatchQuery is a struct representing a query and a list of interfaces
// as context values
type BatchQuery struct {
	Query  string
	Values [][]interface{}
}

// AddValues adds a undetermined list of interfaces to the values
func (bq *BatchQuery) AddValues(vals ...interface{}) {
	bq.Values = append(bq.Values, vals)
}
