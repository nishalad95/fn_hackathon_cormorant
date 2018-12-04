package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"os"

	fdk "github.com/fnproject/fdk-go"
	_ "github.com/go-sql-driver/mysql"
)

type Payload struct {
	Year string `json:"year"`
	Fact string `json:"fact"`
}

var db *sql.DB

func main() {

	fdk.Handle(fdk.HandlerFunc(myHandler))

}

func myHandler(ctx context.Context, in io.Reader, out io.Writer) {
	p := Payload{}
	err := json.NewDecoder(in).Decode(&p)
	if err != nil {
		json.NewEncoder(out).Encode(err.Error())
	}
	err = insertFactoid(p.Year, p.Fact, out)
	if err != nil {
		json.NewEncoder(out).Encode(err.Error())
	}

}

func insertFactoid(year string, fact string, out io.Writer) error {

	user := os.Getenv("user")
	password := os.Getenv("password")
	ip := os.Getenv("ip")
	port := os.Getenv("port")
	dbName := os.Getenv("db")
	dataSource := user + ":" + password + "@tcp(" + ip + ":" + port + ")/" + dbName

	var err error
	db, err = sql.Open("mysql", dataSource)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sqlStatement := `INSERT INTO facts (year, factoid) VALUES (?, ?)`
	_, err = db.Exec(sqlStatement, year, fact)
	if err != nil {
		return err
	}
	msg := struct {
		Msg string `json:"message"`
	}{
		Msg: fmt.Sprintf("Thanks! That's an interesting fact"),
	}
	json.NewEncoder(out).Encode(&msg)
	return nil
}
