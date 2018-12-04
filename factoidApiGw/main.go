package main

import (
    "fmt"
    "log"
    "net/http"
)

func teach(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Teach command received\n")
}

func tell(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Tell command received\n")
}

func main() {
    http.HandleFunc("/teach", teach)
    http.HandleFunc("/tell", tell)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
