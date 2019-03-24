package main

import (
	"fmt"
	"github.com/op/go-logging"
	"os"
	"testing"
)

func failTestOnError(t *testing.T, err error) {
	if err != nil {
		t.Errorf("failed test: %v", err)
	}
}

func getPublicTables() ([]string, error) {
	statement := `
    SELECT tablename
    FROM pg_tables
    WHERE schemaname='public';
  `
	rows, err := db.Raw(statement).Rows()
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot query public tables: %v", err)
		return []string{}, err
	}

	defer rows.Close()

	tables := []string{}
	for rows.Next() {
		var table string
		if err = rows.Scan(&table); err != nil {
			fmt.Fprintf(os.Stderr, "cannot scan table name: %v", err)
			return []string{}, err
		}

		tables = append(tables, table)
	}

	return tables, nil
}

func dropTables() error {
	tables, err := getPublicTables()
	if err != nil {
		return err
	}

	for _, table := range tables {
		if table != "migrations" {
			err = db.Exec(fmt.Sprintf("DROP TABLE %s;", table)).Error
			if err != nil {
				fmt.Fprintf(os.Stderr, "cannot drop %s: %v", table, err)
				return err
			}
		}
	}

	return nil
}

func setupTestDatabase() error {
	if os.Getenv("DATABASE_URL") != "" {
		// set it manually because we need to use commento_test, not commento, by mistake
		os.Setenv("DATABASE_URL", os.Getenv("DATABASE_URL"))
	} else {
		os.Setenv("DATABASE_URL", "commento:123@/commento?charset=utf8mb4&parseTime=True&loc=Local")
	}

	if err := dbConnect(0); err != nil {
		return err
	}

	if err := dropTables(); err != nil {
		return err
	}
	return nil
}

func clearTables() error {
	tables, err := getPublicTables()
	if err != nil {
		return err
	}

	for _, table := range tables {
		err = db.Exec(fmt.Sprintf("DELETE FROM %s;", table)).Error
		if err != nil {
			fmt.Fprintf(os.Stderr, "cannot clear %s: %v", table, err)
			return err
		}
	}

	return nil
}

var setupComplete bool

func setupTestEnv() error {
	if !setupComplete {
		setupComplete = true

		if err := loggerCreate(); err != nil {
			return err
		}

		// Print messages to console only if verbose. Sounds like a good idea to
		// keep the console clean on `go test`.
		if !testing.Verbose() {
			logging.SetLevel(logging.CRITICAL, "")
		}

		if err := setupTestDatabase(); err != nil {
			return err
		}

		if err := markdownRendererCreate(); err != nil {
			return err
		}
	}

	if err := clearTables(); err != nil {
		return err
	}

	return nil
}
