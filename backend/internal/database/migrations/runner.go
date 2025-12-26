package migrations

import (
	"log"

	"gorm.io/gorm"
)

type Migration struct {
	ID   string
	Name string
	Up   func(*gorm.DB) error
}

type MigrationRunner struct {
	db         *gorm.DB
	migrations []Migration
}

func NewMigrationRunner(db *gorm.DB) *MigrationRunner {
	return &MigrationRunner{
		db:         db,
		migrations: []Migration{},
	}
}

func (r *MigrationRunner) AddMigration(id, name string, up func(*gorm.DB) error) {
	r.migrations = append(r.migrations, Migration{
		ID:   id,
		Name: name,
		Up:   up,
	})
}

func (r *MigrationRunner) Run() error {
	for _, migration := range r.migrations {
		log.Printf("Running migration: %s - %s", migration.ID, migration.Name)
		if err := migration.Up(r.db); err != nil {
			return err
		}
		log.Printf("✅ Migration %s completed", migration.ID)
	}
	return nil
}
