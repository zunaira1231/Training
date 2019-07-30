package vault


import (
	"net/http"
	htransport "github.com/go-kit/kit/transport/http"
	"golang.org/x/net/context"
)

// NewHTTPServer makes a new Vault HTTP service.
func NewHTTPServer(ctx context.Context, endpoints Endpoints) http.Handler {
	//create new server
	m := http.NewServeMux()
	//add handler
	//htransport.NewServer give handler to each endpoint
	//encode/decode to marshal/unmarshal json
	m.Handle("/hash", htransport.NewServer(endpoints.HashEndpoint,decodeHashRequest,
		encodeResponse,
	))
	m.Handle("/validate", htransport.NewServer( endpoints.ValidateEndpoint, decodeValidateRequest,
		encodeResponse,
	))
	return m
}