package models

import (
	"encoding/json"

	"github.com/fentec-project/gofe/abe"
)

type Keys struct {
	ID              int
	GlobalPublicKey string
	MasterSecretKey string
}

func CreateKeysFrom(publicKey *abe.FAMEPubKey, secretKey *abe.FAMESecKey) (*Keys, error) {
	pk, err := json.Marshal(publicKey)
	if err != nil {
		return nil, err
	}

	sk, err := json.Marshal(secretKey)
	if err != nil {
		return nil, err
	}

	return &Keys{
		GlobalPublicKey: string(pk),
		MasterSecretKey: string(sk),
	}, nil
}

func (keys *Keys) CreateFameKeysFrom() (*abe.FAMEPubKey, *abe.FAMESecKey, error) {
	var publicKey abe.FAMEPubKey
	err := json.Unmarshal([]byte(keys.GlobalPublicKey), &publicKey)
	if err != nil {
		return nil, nil, err
	}

	var secretKey abe.FAMESecKey
	err = json.Unmarshal([]byte(keys.MasterSecretKey), &secretKey)
	if err != nil {
		return nil, nil, err
	}

	return &publicKey, &secretKey, nil
}
