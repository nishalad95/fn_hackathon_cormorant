package main

import (
	"bufio"
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"net/url"
	"os"
	"strings"

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

	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "&text=") {
			line = line[6:]
			u, err := url.QueryUnescape(line)
			if err != nil {
				fdk.SetHeader(out, "Content-Type", "text/plain")
				_, err := out.Write([]byte("Your message is not parseable"))
				if err != nil {
					panic(err)
				}
				return
			}
			p.Year = u[0:4]
			p.Fact = u[5:]
		}
	}
	if len(p.Fact) == 0 {
		fdk.SetHeader(out, "Content-Type", "text/plain")
		_, err := out.Write([]byte("Your message is empty"))
		if err != nil {
			panic(err)
		}
		return
	}

	err := insertFactoid(p.Year, p.Fact, out)
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

	fdk.SetHeader(out, "Content-Type", "text/plain")
	_, err = out.Write([]byte("Thanks! That's an interesting fact"))
	if err != nil {
		return err
	}
	return nil
}
