package migration

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", "user=postgres, dbname=test, sslmode=disable")
	if err != nil {
		log.Fatal("Error: The data source arguments are not valid")
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("Error: could not establish connection with db")
	}

	_, err = db.Exec("CREATE TABLE pages ( Id integer, Title text, Body text, User_id integer)")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("CREATE TABLE users ( Id integer, Username text,  Password_digest integer, Session_token integer)")
	if err != nil {
		panic(err)
	}
}
