package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/ory/dockertest/v3"
	"go.uber.org/zap/zapcore"
	gMysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

type CloseFunc func() error

// RunInTx begin transaction from given database and execute f.
func RunInTx(ctx context.Context, db *gorm.DB, opts *sql.TxOptions, f func(txDb *gorm.DB) error) error {
	tx := db.WithContext(ctx).Begin(opts)
	if tx.Error != nil {
		return fmt.Errorf("start tx: %v", tx.Error)
	}

	if err := f(tx); err != nil {
		if err1 := tx.Rollback().Error; err1 != nil {
			return fmt.Errorf("rollback tx: %v (original error: %v)", err1, err)
		}
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("commit tx: %v", err)
	}
	return nil
}

// NewTestDatabase start a mysql docker container and returns gorm.DB
func NewTestDatabase(tb testing.TB, migration bool) (*gorm.DB, CloseFunc) {
	tb.Helper()
	var db *sql.DB

	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		tb.Fatalf("Failed to connect to docker: %v", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.Run("mysql", "8.0.17", []string{"MYSQL_ROOT_PASSWORD=secret"})
	if err != nil {
		tb.Fatalf("Failed to not start resource: %v", err)
	}
	err = resource.Expire(60 * 5)

	dcn := fmt.Sprintf("root:secret@(localhost:%s)/mysql?charset=utf8&parseTime=true&multiStatements=true", resource.GetPort("3306/tcp"))
	if err := pool.Retry(func() error {
		var err error
		db, err = sql.Open("mysql", dcn)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("Failed to connect to docker: %v", err)
	}

	gdb, err := gorm.Open(gMysql.New(gMysql.Config{
		Conn: db,
	}), &gorm.Config{
		Logger: NewLogger(time.Second, true, zapcore.FatalLevel),
	})
	if err != nil {
		log.Fatalf("Failed to create a new gorm.DB: %s", err)
	}

	closeFn := func() error {
		_ = db.Close()
		if err := pool.Purge(resource); err != nil {
			log.Fatalf("Failed to purge resource: %s", err)
			return err
		}
		return nil
	}

	if !migration {
		return gdb, closeFn
	}
	err = migrateDB(dcn, "")
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	return gdb, closeFn
}

func DeleteRecordAll(_ testing.TB, db *gorm.DB, tableWhereClauses []string) error {
	if len(tableWhereClauses)%2 != 0 {
		return errors.New("must exist table and where clause")
	}

	for i := 0; i < len(tableWhereClauses)-1; i += 2 {
		rowDB, err := db.DB()
		if err != nil {
			return err
		}
		query := fmt.Sprintf("DELETE FROM %s WHERE %s", tableWhereClauses[i], tableWhereClauses[i+1])
		_, err = rowDB.Exec(query)
		if err != nil {
			return err
		}
	}
	return nil
}

func migrateDB(dcn string, dir string) error {
	db, err := sql.Open("mysql", dcn)
	if err != nil {
		return fmt.Errorf("failed create connect database: %w", err)
	}
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		return fmt.Errorf("failed to mysql instance: %w", err)
	}
	if dir == "" {
		dir = migrationDir()
	}
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", dir),
		"mysql",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to new database instance: %w", err)
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed run migrate: %w", err)
	}
	sourceErr, dbErr := m.Close()
	if sourceErr != nil {
		return fmt.Errorf("failed close source: %w", sourceErr)
	}
	if dbErr != nil {
		return fmt.Errorf("failed close db: %w", dbErr)
	}
	return nil
}

func migrationDir() string {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		return ""
	}
	return filepath.Join(filepath.Dir(filename), "../../migrations")
}
