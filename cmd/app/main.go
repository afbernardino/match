package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"match/cmd/pkg/controller/partners"
	"match/cmd/pkg/repository"

	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(postgres.Open(getDBDSN()), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	repo := repository.NewDatabase(db)

	partnersHandler := partners.NewHandler(repo)
	registerPartnersHandler(r, partnersHandler)

	p := getOSEnv("APP_PORT")
	s := http.Server{
		Addr:         fmt.Sprintf(":%s", p),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 90 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      r,
	}

	err = s.ListenAndServe()
	if err != nil {
		if err != http.ErrServerClosed {
			log.Println(err)
		}
	}
}

func getDBDSN() string {
	host := getOSEnv("PSQL_HOST")
	port := getOSEnv("PSQL_PORT")
	user := getOSEnv("PSQL_USER")
	password := getOSEnv("PSQL_PASSWORD")
	dbName := getOSEnv("PSQL_DB_NAME")

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbName)
}

func getOSEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("please provide and env variable for the key '%s'", key)
	}
	return v
}

func registerPartnersHandler(router *mux.Router, handler partners.Handler) {
	router.HandleFunc("/partners/match", handler.GetMatches).Methods(http.MethodPost)
	router.HandleFunc("/partners/{id:[0-9]+}", handler.GetPartnerById).Methods(http.MethodGet)
}
