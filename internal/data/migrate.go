package data

import (
	"errors"
	"fmt"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func (d *Data) migrate() error {
	run := d.conf.Database.Migrate.Run
	only := d.conf.Database.Migrate.Only
	source := d.conf.Database.Migrate.Source

	if !run {
		return nil
	}

	db, err := d.db.DB()
	if err != nil {
		return err
	}

	var driver database.Driver
	switch strings.ToLower(d.conf.Database.Driver) {
	case "mysql":
		driver, err = mysql.WithInstance(db, &mysql.Config{})
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported database driver: %s", d.conf.Database.Driver)
	}

	m, err := migrate.NewWithDatabaseInstance(
		source,
		d.conf.Database.Driver,
		driver,
	)
	if err != nil {
		return err
	}
	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	if only {
		return errors.New("migrate only")
	}
	return nil
}
