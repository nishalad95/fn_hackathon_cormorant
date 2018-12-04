package main

import (
	"bufio"
	"context"
	"database/sql"
	"io"
	"os"
	"regexp"
	"strings"

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

	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "&text=") {
			p.Message = line[6:]
		}
	}
	if len(p.Message) == 0 {
		respondWithFactoid("Your message is empty", out)
	}

	year := retrieveYear(p)

	factoid := retrieveFactoid(year)
	respondWithFactoid(factoid, out)
}

func retrieveYear(p Payload) string {
	re := regexp.MustCompile("[0-9]+")
	year := re.FindAllString(p.Message, -1)

	if year == nil {
		return "Hmmm... It looks like you didn't enter a year, please try again :-)"
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

	fdk.SetHeader(out, "Content-Type", "text/plain")
	_, err := out.Write([]byte(fact))
	if err != nil {
		panic(err)
	}

}
