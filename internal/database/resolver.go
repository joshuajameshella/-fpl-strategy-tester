package database

import (
	"database/sql"
	"fmt"

	"github.com/doug-martin/goqu/v9"

	// Needed to construct SQL queries
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	_ "github.com/go-sql-driver/mysql"
)

// Resolver is the entry-point for accessing the football data
type Resolver struct {
	FPLDB      *sql.DB
	sqlBuilder *goqu.Database
}

// NewResolver creates and returns an empty Resolver
func NewResolver() *Resolver {
	return &Resolver{}
}

// ResolveFPLDB returns or initiates a new database connection
func (r *Resolver) ResolveFPLDB() *sql.DB {
	if r.FPLDB == nil {
		fmt.Printf("Initiating a new database connection to : %v\n", dbSchemaName)
		databaseLogin := fmt.Sprintf("%v:%v@tcp(%v)/%v", dbUsername, dbPassword, dbAddress, dbSchemaName)
		conn, err := sql.Open("mysql", databaseLogin)
		if err != nil {
			fmt.Printf("Error while resolving database: %v\n", err)
			return nil
		}
		r.FPLDB = conn
	}
	return r.FPLDB
}

// ResolveMySQLQueryBuilder creates a new goqu-based query builder, using the mysql dialect.
func (r *Resolver) ResolveMySQLQueryBuilder() *goqu.Database {
	if r.sqlBuilder == nil {
		r.sqlBuilder = goqu.New("mysql", nil)
	}
	return r.sqlBuilder
}
