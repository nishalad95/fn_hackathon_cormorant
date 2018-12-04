package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/oracle/oci-go-sdk/common"
	"github.com/sirupsen/logrus"

	"github.com/nishalad95/fn_hackathon_cormorant/factoidApiGw/invoker"
)

func teach(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Teach command received\n")
}

func tell(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Tell command received\n")
}

func main() {
	cfg, err := common.ConfigurationProviderFromFile("/hackathon/ociconfig", "")
	if err != nil {
		logrus.WithError(err).Error("Error creating configuration provider")
	}

	teacher, err := invoker.NewFunctionInvoker(cfg, "somdzwpae4a", "us-phoenix-1", "ocid1.fnfunc.oc1.phx.aaaaaaaaabzqhfe6z5vy4tjrm5oyfyoksonjmkl6dodhxyyrwi4uqkmzlrlq")
	if err != nil {
		logrus.WithError(err).Error("Error creating teach function handler")
	}
	http.Handle("/teach", teacher)

	teller, err := invoker.NewFunctionInvoker(cfg, "somdzwpae4a", "us-phoenix-1", "ocid1.fnfunc.oc1.phx.aaaaaaaaad3y62jgjpcqk47bzzcvm37qbfgjv6247nyc7jrt7x4z7enphdgq")
	if err != nil {
		logrus.WithError(err).Error("Error creating tell function handler")
	}
	http.Handle("/tell", teller)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
