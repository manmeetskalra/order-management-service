# Order Management Service
Prerequisites
- Go (version 1.16 or later)
- MySQL server (version 5.7 or later)

## Getting started
Clone the repository and change to the project directory:

`git clone https://github.com/manmeetskalra/order-management-service/tree/master`

`cd order-management-service`

Set up the database by executing the SQL script in database/schema.sql:

`mysql -u <username> -p < database/schema.sql`

Once you are connected to the MySQL server, create a new database to store the order information.

`CREATE DATABASE order_management;`

Switch to the newly created database:

`USE order_management;`

Create a table to store the orders. 
```
CREATE TABLE orders (
  id VARCHAR(255) PRIMARY KEY,
  status VARCHAR(50),
  items JSON,
  total DECIMAL(10,2),
  currencyUnit VARCHAR(10)
);
```

Install the required dependencies:

`go mod tidy`

Build and run the service:

`go run main.go`
