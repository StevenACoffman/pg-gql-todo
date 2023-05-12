package sqldb

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	pgxmigrate "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/mitchellh/mapstructure"

	"github.com/StevenACoffman/pg-gql-todo/assets"
	"github.com/StevenACoffman/pg-gql-todo/generated/gql/model"
	"github.com/StevenACoffman/pg-gql-todo/generated/todosql"
)

func NewDBPool(dbInfo *DBInfo, automigrate bool) (*sql.DB, error) {
	connString := dbInfo.ConnectionString()
	c, err := pgx.ParseConfig(connString)
	if err != nil {
		return nil, err
	}

	// Register registers a DriverConfig and
	// obtains a connection String for use with sql.Open.
	registeredConnString := stdlib.RegisterConnConfig(c)

	// opening a driver typically will not attempt to connect to the database.
	// any parse or other error here does not require a Close
	pool, err := sql.Open("pgx", registeredConnString)
	if err != nil {
		return nil, err
	}

	// Reasonable defaults here should be tuned by observing application
	// PostgreSQL maxes at 500 open connections, so 20 app instances
	// may consume all available.
	pool.SetMaxOpenConns(25)
	pool.SetMaxIdleConns(25)
	pool.SetConnMaxIdleTime(5 * time.Minute)
	pool.SetConnMaxLifetime(2 * time.Hour)
	if automigrate {
		err = MigrateDB(pool, dbInfo)
		if err != nil {
			return nil, err
		}
	}

	return pool, nil
}

func MigrateDB(stdpool *sql.DB, dbInfo *DBInfo) error {
	iofsDriver, err := iofs.New(assets.EmbeddedFiles, "migrations")
	if err != nil {
		return err
	}
	defer iofsDriver.Close()

	migrateDriver, err := pgxmigrate.WithInstance(stdpool, &pgxmigrate.Config{
		DatabaseName: dbInfo.DBName,
		SchemaName:   dbInfo.DBSchema,
	})
	if err != nil {
		return err
	}
	logName := fmt.Sprintf("%s.%s", dbInfo.DBName, dbInfo.DBSchema)
	migrator, err := migrate.NewWithInstance("iofs", iofsDriver, logName, migrateDriver)
	if err != nil {
		return err
	}
	migrator.Log = &Logger{}

	err = migrator.Up()
	switch {
	case errors.Is(err, migrate.ErrNoChange):
		break
	case err != nil:
		return err
	}

	return nil
}

type DBInfo struct {
	DBUser   string `mapstructure:"user,omitempty"`
	DBPass   string `mapstructure:"password,omitempty"`
	DBHost   string `mapstructure:"host,omitempty"`
	DBPort   string `mapstructure:"port,omitempty"`
	DBName   string `mapstructure:"dbname,omitempty"`
	DBSchema string `mapstructure:"search_path,omitempty"`
}

func NewDBInfo(
	user string,
	password string,
	host string,
	dbName string,
	dbSchema string,
) *DBInfo {
	dbInfo := &DBInfo{
		DBUser:   user,
		DBPass:   password,
		DBHost:   host,
		DBPort:   "5432",
		DBName:   dbName,
		DBSchema: dbSchema,
	}
	return dbInfo
}

func (dbInfo *DBInfo) ConnectionString() string {
	x := map[string]string{}

	err := mapstructure.Decode(dbInfo, &x)
	if err != nil {
		return ""
	}

	pairs := []string{"sslmode=disable"}
	for k, v := range x {
		pairs = append(pairs, fmt.Sprint(k, "=", v))
	}
	return strings.Join(pairs, " ")
}

type Logger struct{}

func (l *Logger) Printf(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}

func (l *Logger) Verbose() bool {
	return true
}

func ConvertSQLtoGQLTodo(sqltodo *todosql.Todo) *model.Todo {
	return &model.Todo{
		ID:   sqltodo.ID.String(),
		Text: sqltodo.Description,
		Done: sqltodo.Done,
	}
}
