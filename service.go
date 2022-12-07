package main

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"io"
	"keyValueStore/api"
	"keyValueStore/logger"
	"log"
	"net/http"
	"os"
)

func keyValuePutHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	value, err := io.ReadAll(r.Body)
	defer func() {
		err = r.Body.Close()
		if err != nil {
			log.Printf(err.Error())
		}
	}()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = api.Put(key, string(value))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logger2.WritePut(key, string(value))

	w.WriteHeader(http.StatusCreated)

	log.Printf("PUT key:%s value :%s\n", key, string(value))
}

func keyValueGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	value, err := api.Get(key)
	if errors.Is(api.ErrorNoSuchKey, err) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write([]byte(value))
	if err != nil {
		log.Printf(err.Error())
		return
	}

	log.Printf("GET key=%s\n", key)
}

func keyValueDeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	err := api.Delete(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logger2.WriteDelete(key)
	w.WriteHeader(http.StatusOK)

	log.Printf("DELETE key=%s\n", key)
}

var logger2 logger.TransactionLogger

func initTransactionLog() error {
	var err error

	err = godotenv.Load()
	if err != nil {
		return fmt.Errorf("error loading .env file")
	}
	dbName := os.Getenv("DB_NAME")
	host := os.Getenv("HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")

	logger2, err = logger.NewPostgresTransactionLogger(logger.PostgresDBParams{
		DbName:   dbName,
		Host:     host,
		User:     user,
		Password: password,
	})
	// logger, err = NewFileTransactionLogger("tmp/transaction.log")
	if err != nil {
		return fmt.Errorf("failed to create event logger: %w", err)
	}
	events, errors2 := logger2.ReadEvents()
	count, e, ok := 0, logger.Event{}, true

	for ok && err == nil {
		select {
		case err, ok = <-errors2:
		case e, ok = <-events:
			switch e.EventType {
			case logger.EventDelete:
				err = api.Delete(e.Key)
				count++
			case logger.EventPut:
				err = api.Put(e.Key, e.Value)
				count++
			}
		}
	}
	log.Printf("%d events replayed\n", count)
	logger2.Run()

	return err
}

func main() {
	err := initTransactionLog()
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	r := mux.NewRouter()
	r.HandleFunc("/v1/{key}", keyValuePutHandler).Methods("PUT")
	r.HandleFunc("/v1/{key}", keyValueGetHandler).Methods("GET")
	r.HandleFunc("/v1/{key}", keyValueDeleteHandler).Methods("DELETE")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", r))
}
