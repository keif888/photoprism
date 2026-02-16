package dsn

import (
	"fmt"
	"os"
)

// PhotoPrismTestToDriverDSN function to transform PHOTOPRISM_TEST_DSN environment variables to driver and dsn strings
func PhotoPrismTestToDriverDSN(dbn int) (driver string, dsn string) {
	dsnName := os.Getenv("PHOTOPRISM_TEST_DSN_NAME")
	switch dsnName {
	case "sqlite":
		driver = "sqlite"
		dsn = os.Getenv("PHOTOPRISM_TEST_DSN_SQLITE")
	case "sqlitefile":
		driver = "sqlite"
		dsn = os.Getenv("PHOTOPRISM_TEST_DSN_SQLITEFILE")
		if dbn > 0 {
			d := Parse(dsn)
			d.Driver = driver
			d.Name = fmt.Sprintf("%s_%02d.db", d.Name, dbn)
			dsn = d.ToString()
		}
	case "mariadb":
		driver = "mysql"
		dsn = os.Getenv("PHOTOPRISM_TEST_DSN_MARIADB")
		if dbn > 0 {
			d := Parse(dsn)
			d.Driver = driver
			d.Name = fmt.Sprintf("%s_%02d", d.Name, dbn)
			dsn = d.ToString()
		}
	case "mysql8":
		driver = "mysql"
		dsn = os.Getenv("PHOTOPRISM_TEST_DSN_MYSQL8")
		if dbn > 0 {
			d := Parse(dsn)
			d.Driver = driver
			d.Name = fmt.Sprintf("%s_%02d", d.Name, dbn)
			dsn = d.ToString()
		}
	case "postgres":
		driver = "postgres"
		dsn = os.Getenv("PHOTOPRISM_TEST_DSN_POSTGRES")
		if dbn > 0 {
			d := Parse(dsn)
			d.Driver = driver
			d.Name = fmt.Sprintf("%s_%02d", d.Name, dbn)
			dsn = d.ToString()
		}
	default:
		driver = "sqlite"
		dsn = ""
	}

	return driver, dsn
}

// SetDSNToEnv pushes the dsn to the appropriate environment variable for testing.
func SetDSNToEnv(dsn string) {
	dsnName := os.Getenv("PHOTOPRISM_TEST_DSN_NAME")
	switch dsnName {
	case "sqlite":
		os.Setenv("PHOTOPRISM_TEST_DSN_SQLITE", dsn)
	case "sqlitefile":
		os.Setenv("PHOTOPRISM_TEST_DSN_SQLITEFILE", dsn)
	case "mariadb":
		os.Setenv("PHOTOPRISM_TEST_DSN_MARIADB", dsn)
	case "mysql8":
		os.Setenv("PHOTOPRISM_TEST_DSN_MYSQL8", dsn)
	case "postgres":
		os.Setenv("PHOTOPRISM_TEST_DSN_POSTGRES", dsn)
	}
}
