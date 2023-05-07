package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type Order struct {
	ID           string  `json:"id"`
	Status       string  `json:"status"`
	Items        []Item  `json:"items"`
	Total        float64 `json:"total"`
	CurrencyUnit string  `json:"currencyUnit"`
}

type Item struct {
	ID          string  `json:"id"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
}

func main() {
	fmt.Println("Beginning of the main function")
	db, err := sql.Open("mysql", "root:jattnjuliet2@tcp(localhost:3306)/order_management")
	fmt.Println("After mysql in the main function")
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/orders/add", func(w http.ResponseWriter, r *http.Request) {
		addOrderHandler(w, r, db)
	})

	http.HandleFunc("/orders/update", func(w http.ResponseWriter, r *http.Request) {
		updateOrderStatusHandler(w, r, db)
	})

	http.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		fetchOrdersHandler(w, r, db)
	})

	fmt.Println("Before localhost call in the main function")
	log.Fatal(http.ListenAndServe(":8080", nil))
	fmt.Println("After localhost call in the main function")
}

func updateOrderStatusHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	id := r.URL.Path[len("/orders/"):]

	if r.Method != "PUT" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	status := r.URL.Query().Get("status")

	_, err := db.Exec("UPDATE orders SET status = ? WHERE id = ?", status, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func addOrderHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var order Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("Order", order)
	itemsJSON, err := json.Marshal(order.Items)
	if err != nil {
		log.Fatal(err)
		return
	}

	_, err = db.Exec("INSERT INTO orders (id, status, items, total, currencyUnit) VALUES (?, ?, ?, ?, ?)", order.ID, order.Status, itemsJSON, order.Total, order.CurrencyUnit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func fetchOrdersHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	status := r.URL.Query().Get("status")
	currency := r.URL.Query().Get("currency")
	sortBy := r.URL.Query().Get("sort_by")
	if sortBy == "" {
		sortBy = "id"
	}

	query := "SELECT * FROM orders WHERE 1 = 1"
	args := make([]interface{}, 0)

	if status != "" {
		query += " AND status = ?"
		args = append(args, status)
	}
	if currency != "" {
		query += " AND currencyUnit = ?"
		args = append(args, currency)
	}

	query += " ORDER BY " + sortBy

	rows, err := db.Query(query, args...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	orders := make([]Order, 0)
	for rows.Next() {
		var order Order
		var itemsJSON []byte

		err := rows.Scan(&order.ID, &order.Status, &itemsJSON, &order.Total, &order.CurrencyUnit)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = json.Unmarshal(itemsJSON, &order.Items)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}
