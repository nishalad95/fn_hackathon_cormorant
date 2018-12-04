package invoker

import (
	"context"
	"fmt"
	"net/http"

	"github.com/oracle/oci-go-sdk/common"
	"github.com/sirupsen/logrus"
)

// FunctionInvoker is a http Handler which calls an OCI function, handling request signing, etc.
type FunctionInvoker struct {
	client    common.BaseClient
	invokeURL string
}

// NewFunctionInvoker creates a function invoker for a given function using auth config from a configuration provider
func NewFunctionInvoker(cfg common.ConfigurationProvider, appShortCode, region, fnID string) (*FunctionInvoker, error) {
	client, err := common.NewClientWithConfig(cfg)
	if err != nil {
		return nil, err
	}
	client.Host = fmt.Sprintf("%s.%s.functions.oci.oraclecorp.com", appShortCode, region)
	return &FunctionInvoker{
		client:    client,
		invokeURL: fmt.Sprintf("%s.%s.fuctions.oci.oraclecloud.com/invoke/%s", appShortCode, region, fnID),
	}, nil
}

func (fi FunctionInvoker) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	req, err := http.NewRequest(http.MethodPost, fi.invokeURL, r.Body)
	if err != nil {
		logrus.WithError(err).Error("Error creating functions request")
	}
	fi.client.Call(context.TODO(), req)
}

var _ http.Handler = &FunctionInvoker{}
