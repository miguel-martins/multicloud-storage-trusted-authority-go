package main

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/miguel-martins/multicloud-storage-trusted-authority-go/internal/handlers"
	"github.com/miguel-martins/multicloud-storage-trusted-authority-go/internal/repositories"
	"github.com/miguel-martins/multicloud-storage-trusted-authority-go/internal/util"
)

type Application struct {
	Server *http.Server
	DB     *sql.DB
}

func NewApplication(serverCert, serverKey string) (*Application, error) {
	serverTLSConfig, err := loadServerTLSConfig(serverCert, serverKey)
	if err != nil {
		return nil, err
	}

	server := &http.Server{
		Addr:      ":443",
		TLSConfig: serverTLSConfig,
	}

	db, err := util.InitDB()
	if err != nil {
		return nil, err
	}

	return &Application{
		Server: server,
		DB:     db,
	}, nil
}

func (app *Application) Start() error {

	keyRepository := repositories.NewKeyRepository(app.DB)

	http.HandleFunc("/gpk", handlers.GenerateKeysHandler(keyRepository))
	log.Println("Starting Multicloud Storage Trusted Authority service...")
	log.Println(`
	    __  ___      ____  _      __                __   _____ __
	   /  |/  /_  __/ / /_(_)____/ /___  __  ______/ /  / ___// /_____  _________  ____  ___
	  / /|_/ / / / / / __/ / ___/ / __ \/ / / / __  /   \__ \/ __/ __ \/ ___/ __ \/ __ \/ _ \
	 / /  / / /_/ / / /_/ / /__/ / /_/ / /_/ / /_/ /   ___/ / /_/ /_/ / /  / /_/ / /_/ /  __/
	/_/  /_/\__,_/_/\__/_/\___/_/\____/\__,_/\__,_/   /____/\__/\____/_/   \__,_/\__, /\___/
	                                                                            /____/		`)
	log.Println("Server started on port 443")
	return app.Server.ListenAndServeTLS("", "")
}

func loadServerTLSConfig(certFile, keyFile string) (*tls.Config, error) {
	server, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	cert, err := os.ReadFile("ca.crt")
	if err != nil {
		log.Fatalf("Failed to read certificate file: %v", err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(cert)

	return &tls.Config{
		Certificates: []tls.Certificate{server},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    caCertPool,
	}, nil
}

func main() {
	app, err := NewApplication("server.crt", "server.key")
	if err != nil {
		log.Fatal("Error initializing application: ", err)
	}

	if err := app.Start(); err != nil {
		log.Fatal("Error starting server:", err)
	}
}
