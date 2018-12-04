package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"

	fdk "github.com/fnproject/fdk-go"
	_ "github.com/go-sql-driver/mysql"
)

type Payload struct {
	Message string `json:"year"`
}

var db *sql.DB

func main() {

	fdk.Handle(fdk.HandlerFunc(myHandler))

}

func myHandler(ctx context.Context, in io.Reader, out io.Writer) {
	p := Payload{}
	err := json.NewDecoder(in).Decode(&p)
	if err != nil {
		respondWithFactoid(err.Error(), out)
	}
	year := retrieveYear(p)

	factoid := retrieveFactoid(year)
	respondWithFactoid(factoid, out)
}

func retrieveYear(p Payload) string {
	re := regexp.MustCompile("[0-9]+")
	year := re.FindAllString(p.Message, -1)

	if year == nil {
		return "could not get year"
	}
	return year[0]
}

func retrieveFactoid(year string) string {

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

	var factoid string
	sqlStatement := `SELECT factoid FROM facts WHERE year=?`
	row := db.QueryRow(sqlStatement, year)
	err = row.Scan(&factoid)
	if err != nil {
		if err == sql.ErrNoRows {
			return "Hmmm... It looks like I couldn't find a fact for that year, use the '/teach' command to teach me something historical that happened in " + year + "?"
		}
		return err.Error()
	}
	return factoid
}

func respondWithFactoid(fact string, out io.Writer) {
	msg := struct {
		Msg string `json:"factoid"`
	}{
		Msg: fmt.Sprintf("%s", fact),
	}
	json.NewEncoder(out).Encode(&msg)
}
