package migrations

import (
	"dental-marketplace/backend/internal/database/migrations/seeds"
	"log"

	"gorm.io/gorm"
)

func RunAll(db *gorm.DB) error {
	log.Println("🔄 Running migrations...")

	runner := NewMigrationRunner(db)

	// Add migrations in order
	runner.AddMigration("001", "Create Constants Tables", CreateConstantsTables)
	runner.AddMigration("002", "Create Business Tables", CreateBusinessTables)

	// Run migrations
	if err := runner.Run(); err != nil {
		return err
	}

	log.Println("✅ All migrations completed")

	// Seed constants (always runs to ensure they're up to date)
	if err := seeds.SeedConstants(db); err != nil {
		return err
	}

	// Seed sample data (only if database is empty)
	if err := seeds.SeedSampleData(db); err != nil {
		return err
	}

	return nil
}
