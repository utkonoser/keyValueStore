package main

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
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
	err = Put(key, string(value))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logger.WritePut(key, string(value))

	w.WriteHeader(http.StatusCreated)

	log.Printf("PUT key:%s value :%s\n", key, string(value))
}

func keyValueGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	value, err := Get(key)
	if errors.Is(ErrorNoSuchKey, err) {
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

	err := Delete(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logger.WriteDelete(key)
	w.WriteHeader(http.StatusOK)

	log.Printf("DELETE key=%s\n", key)
}

var logger TransactionLogger

func initTransactionLog() error {
	var err error
	logger, err = NewPostgresTransactionLogger(PostgresDBParams{
		dbName:   "postgres",
		host:     "postgresKV",
		user:     "postgres",
		password: "1234qwer",
	})
	// logger, err = NewFileTransactionLogger("tmp/transaction.log")
	if err != nil {
		return fmt.Errorf("failed to create event logger: %w", err)
	}
	events, errors2 := logger.ReadEvents()
	count, e, ok := 0, Event{}, true

	for ok && err == nil {
		select {
		case err, ok = <-errors2:
		case e, ok = <-events:
			switch e.EventType {
			case EventDelete:
				err = Delete(e.Key)
				count++
			case EventPut:
				err = Put(e.Key, e.Value)
				count++
			}
		}
	}
	log.Printf("%d events replayed\n", count)
	logger.Run()

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
	log.Fatal(http.ListenAndServe("localhost:8080", r))
}
