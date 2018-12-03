package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"regexp"

	fdk "github.com/fnproject/fdk-go"
)

func main() {
	fdk.Handle(fdk.HandlerFunc(myHandler))

	// json decode the incoming payload that this function receives
	// from json grab the number payload
	// run a sql query using number payload
	// return a response to user

}

type Person struct {
	Name string `json:"name"`
}

type Payload struct {
	Message string //TODO json
}

func myHandler(ctx context.Context, in io.Reader, out io.Writer) {
	p := &Payload{}
	json.NewDecoder(in).Decode(p)
	year := retrieveYear(p)
	factoid := retrieveFactoid(year)
	respondWithFactoid(factoid, out)
}

func decodeJSON(ctx context.Context, in io.Reader) *Payload {
	p := &Payload{}
	json.NewDecoder(in).Decode(p)
	return p
}

// retrieve all numbers from the user message
func retrieveYear(p *Payload) string {
	re := regexp.MustCompile("[0-9]+")
	year := re.FindAllString(p.Message, -1)

	return year[0]
}

func retrieveFactoid(year string) string {
	db, err := sql.Open("mysql", "theUser:thePassword@/theDbName")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var factoid string
	sqlStatement := `SELECT factoid FROM facts WHERE year=$` + year + ``
	row := db.QueryRow(sqlStatement, 1)
	err = row.Scan(&factoid)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Zero rows found")
		} else {
			panic(err)
		}
	}
	return factoid
}

func respondWithFactoid(fact string, out io.Writer) {
	msg := struct {
		Msg string `json:"message"`
	}{
		Msg: fmt.Sprintf("%s", fact),
	}
	json.NewEncoder(out).Encode(&msg)
}
