package keygen

import (
	"log"

	"github.com/fentec-project/gofe/abe"
	"github.com/miguel-martins/multicloud-storage-trusted-authority-go/internal/models"
	repository "github.com/miguel-martins/multicloud-storage-trusted-authority-go/internal/repositories"
)

var (
	globalPublicKey *abe.FAMEPubKey
	masterSecretKey *abe.FAMESecKey
)

func GetGlobalPublicKey() *abe.FAMEPubKey {
	return globalPublicKey
}

func InitKeys(keyRepository *repository.KeyRepository) error {
	log.Printf("Initializing keys...")

	if keysExist(keyRepository) {
		log.Printf("Keys already exist in the database. Skipping key generation.")
		pk, sk, err := getKeys(keyRepository)
		if err != nil {
			return err
		}

		globalPublicKey = pk
		masterSecretKey = sk

		return nil
	}

	var err error
	globalPublicKey, masterSecretKey, err = generateKeys()
	if err != nil {
		log.Printf("Error generating keys")
		return err
	}

	err = storeKeys(keyRepository, globalPublicKey, masterSecretKey)
	if err != nil {
		log.Printf("Error storing keys in the database")
		return err
	}

	return nil
}

func keysExist(repository *repository.KeyRepository) bool {
	hasKeys, err := repository.HasKeys()
	if err != nil {
		return false
	}
	return hasKeys
}

func getKeys(repository *repository.KeyRepository) (*abe.FAMEPubKey, *abe.FAMESecKey, error) {
	keys, err := repository.Get()
	if err != nil {
		return nil, nil, err
	}

	pk, sk, err := keys.CreateFameKeysFrom()
	if err != nil {
		return nil, nil, err
	}
	return pk, sk, nil
}

func storeKeys(repository *repository.KeyRepository, publicKey *abe.FAMEPubKey, secretKey *abe.FAMESecKey) error {
	keys, err := models.CreateKeysFrom(publicKey, secretKey)
	if err != nil {
		return err
	}

	err = repository.Save(keys)
	if err != nil {
		return err
	}

	return nil
}

func generateKeys() (*abe.FAMEPubKey, *abe.FAMESecKey, error) {
	return abe.NewFAME().GenerateMasterKeys()
}
