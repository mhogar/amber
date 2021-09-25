package firestoreadapter

import (
	"github.com/mhogar/amber/models"
)

func (crud *FirestoreCRUD) Setup() error {
	return nil
}

func (crud *FirestoreCRUD) CreateMigration(timestamp string) error {
	return nil
}

func (crud *FirestoreCRUD) GetMigrationByTimestamp(timestamp string) (*models.Migration, error) {
	return nil, nil
}

func (crud *FirestoreCRUD) GetLatestTimestamp() (string, bool, error) {
	return "", false, nil
}

func (crud *FirestoreCRUD) DeleteMigrationByTimestamp(timestamp string) error {
	return nil
}
