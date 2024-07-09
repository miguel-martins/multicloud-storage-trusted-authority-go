package handlers

import (
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/miguel-martins/multicloud-storage-trusted-authority-go/internal/keygen"
	repository "github.com/miguel-martins/multicloud-storage-trusted-authority-go/internal/repositories"
)

type GenerateKeysResponse struct {
	GlobalPublicKey string `json:"global_public_key"`
}

func GenerateKeysHandler(keyRepository *repository.KeyRepository) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		err := keygen.InitKeys(keyRepository)
		if err != nil {
			http.Error(w, "Failed to generate keys", http.StatusInternalServerError)
		}

		globalPublicKey := keygen.GetGlobalPublicKey()
		globalPublicKeyBytes, err := json.Marshal(globalPublicKey)
		if err != nil {
			http.Error(w, "Failed to encode global public key", http.StatusInternalServerError)
			return
		}

		globalPublicKeyBase64 := base64.StdEncoding.EncodeToString(globalPublicKeyBytes)

		response := GenerateKeysResponse{
			GlobalPublicKey: globalPublicKeyBase64,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}
}
