package repositories

import (
	"context"
	"database/sql"

	"github.com/miguel-martins/multicloud-storage-trusted-authority-go/internal/models"
)

type KeyRepository struct {
	DB *sql.DB
}

func NewKeyRepository(db *sql.DB) *KeyRepository {
	return &KeyRepository{DB: db}
}

func (kr *KeyRepository) Save(keys *models.Keys) error {
	_, err := kr.DB.ExecContext(context.Background(), "INSERT INTO global_keys (public_key, private_key) VALUES ($1, $2)", keys.GlobalPublicKey, keys.MasterSecretKey)
	if err != nil {
		return err
	}
	return nil
}

func (kr *KeyRepository) Get() (*models.Keys, error) {
	var keys *models.Keys
	err := kr.DB.QueryRowContext(context.Background(), "SELECT * FROM global_keys limit 1").Scan(&keys)
	if err != nil {
		return nil, err
	}
	return keys, nil
}

func (kr *KeyRepository) HasKeys() (bool, error) {
	err := kr.DB.QueryRowContext(context.Background(), "SELECT count(*) FROM global_keys having count(*) = 1").Scan()
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil // No keys
		}
		return false, err
	}

	return true, nil
}
