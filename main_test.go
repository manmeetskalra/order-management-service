package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAddOrderHandler(t *testing.T) {
	order := Order{
		ID:           "1",
		Status:       "new",
		Items:        []Item{{ID: "1", Description: "item 1", Price: 1.0, Quantity: 1}},
		Total:        1.0,
		CurrencyUnit: "USD",
	}
	body, err := json.Marshal(order)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/orders/add", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	db, err := sql.Open("mysql", "root:jattnjuliet2@tcp(localhost:3306)/order_management")
	if err != nil {
		t.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	addOrderHandler(rr, req, db)
	fmt.Printf("items: %+v\n", order.Items)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"status":"new","items":[{"id":"1","description":"item 1","price":1,"quantity":1}],"total":1,"currencyUnit":"USD"}`
	if rr.Body.String() != expected {
		t.Errorf("unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestUpdateOrderStatusHandler(t *testing.T) {
	db, err := sql.Open("mysql", "root:jattnjuliet2@tcp(localhost:3306)/order_management")
	if err != nil {
		t.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	order := Order{
		ID:           "2",
		Status:       "new",
		Items:        []Item{{ID: "1", Description: "item 1", Price: 1.0, Quantity: 1}},
		Total:        1.0,
		CurrencyUnit: "USD",
	}

	itemsJSON, _ := json.Marshal(order.Items)
	_, err = db.Exec("INSERT INTO orders (id, status, items, total, currencyUnit) VALUES (?, ?, ?, ?, ?)", order.ID, order.Status, itemsJSON, order.Total, order.CurrencyUnit)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("PUT", "/orders/"+order.ID+"?status=processed", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	updateOrderStatusHandler(rr, req, db)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("wrong status code: got %v want %v", status, http.StatusOK)
	}

	var updatedStatus string
	row := db.QueryRow("SELECT status FROM orders WHERE id = ?", order.ID)
	err = row.Scan(&updatedStatus)
	if err != nil {
		t.Fatal(err)
	}
	if updatedStatus != "processed" {
		t.Errorf("order status not updated in the database: got %v want processed", updatedStatus)
	}
}

func TestFetchOrdersHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/orders?status=new&currency=USD&sort_by=id", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	db, err := sql.Open("mysql", "root:jattnjuliet2@tcp(localhost:3306)/order_management")
	if err != nil {
		t.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	fetchOrdersHandler(rr, req, db)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `[{"id":"1","status":"new","items":[{"id":"1","description":"item 1","price":1,"quantity":1}],"total":1,"currencyUnit":"USD"}]`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
