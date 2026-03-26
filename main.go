package main

import "database/sql"

func main() {

	connStr := "user=postgres password=dbpass dbname=calculator sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

}
