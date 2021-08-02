package models

import (
	"regexp"

	"github.com/mhogar/migrationrunner"
)

const (
	ValidateMigrationValid            = 0x0
	ValidateMigrationInvalidTimestamp = 0x1
)

// Migration represents the migration model
type Migration struct {
	Timestamp string
}

type MigrationCRUD interface {
	migrationrunner.MigrationCRUD

	// GetMigrationByTimestamp fetches the migration with the matching timestamp
	// If no migrations are found, returns nil migration. Also returns any errors
	GetMigrationByTimestamp(timestamp string) (*Migration, error)
}

func CreateNewMigration(timestamp string) *Migration {
	return &Migration{
		Timestamp: timestamp,
	}
}

// Validate validates the migration is a valid migration model
// Returns an int indicating which fields are invalid
func (m Migration) Validate() int {
	code := ValidateMigrationValid

	matched, _ := regexp.MatchString(`^\d{14}$`, m.Timestamp)
	if !matched {
		code |= ValidateMigrationInvalidTimestamp
	}

	return code
}
