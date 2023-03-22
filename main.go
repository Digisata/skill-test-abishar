package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Sale struct {
	ID        int     `json:"id"`
	Customer  string  `json:"customer"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
	Timestamp string  `json:"timestamp"`
}

type RequestBody struct {
	RequestID int    `json:"request_id"`
	Data      []Sale `json:"data"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUsername, dbPassword, dbName)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()

	http.HandleFunc("/sales", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var requestBody RequestBody
		err := json.NewDecoder(r.Body).Decode(&requestBody)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		errChan := make(chan error, 1)

		batchSum := 10
		dataLength := len(requestBody.Data)
		div := dataLength / batchSum
		mod := dataLength % batchSum
		wg := sync.WaitGroup{}
		tx, _ := db.Begin()

		if dataLength != 0 {
			for i := 0; i < batchSum; i++ {
				salesData := requestBody.Data[div*i : div*(i+1)]
				if i == batchSum-1 && mod != 0 {
					salesData = requestBody.Data[div*i : div*(i+1)+mod]
				}

				if div == 0 {
					salesData = requestBody.Data
					i = batchSum - 1
				}

				wg.Add(1)
				go insertIntoSales(salesData, tx, &wg, errChan)
			}
		}

		wg.Wait()
		select {
		case err := <-errChan:
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Sales data fail: %s", err.Error())
		default:
			tx.Commit()
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, "Sales data received successfully")
		}
	})

	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}

func insertIntoSales(data []Sale, tx *sql.Tx, wg *sync.WaitGroup, ch chan<- error) {
	defer wg.Done()
	for _, sale := range data {
		query := `INSERT INTO sales (id, customer, quantity, price, timestamp) VALUES ($1, $2, $3, $4, $5)`
		_, err := tx.Exec(query, sale.ID, sale.Customer, sale.Quantity, sale.Price, sale.Timestamp)
		if err != nil {
			tx.Rollback()
			ch <- err
		}
	}
}
