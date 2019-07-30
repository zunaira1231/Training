package vault

import (
	"encoding/json"
	"errors"
	"net/http"
	//"context"
	"github.com/go-kit/kit/endpoint"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
)

// Service provides password hashing capabilities.
type Service interface {
	//hash need only password and return hash or error
	//validate need hash and password and return true/false or error
	Hash(ctx context.Context, password string) (string, error)
	Validate(ctx context.Context, password, hash string) (bool, error)
}

// NewService makes a new Service.
//Service constructor
//constructor is a global function that return usable instance of struct
func NewService() Service {
	return vaultService{}
}

type vaultService struct{}
//it take password and return hash
func (vaultService) Hash(ctx context.Context, password string) (string, error) {
	//create hash
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
//take password and hash and return true or false for validation
//return nil in error when password validate
func (vaultService) Validate(ctx context.Context, password, hash string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false, nil
	}
	return true, nil
}
//for request and response in and out of services
//add struct for each type of msg that service will accept or return

type hashRequest struct {
	Password string `json:"password"`
}
type hashResponse struct {
	Hash string `json:"hash"`
	Err  string `json:"err,omitempty"`
}
type validateRequest struct {
	Password string `json:"password"`
	Hash     string `json:"hash"`
}
type validateResponse struct {
	Valid bool   `json:"valid"`
	Err   string `json:"err,omitempty"`
}

//for decoding json body of http.request to service.go

func decodeHashRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req hashRequest
	//unmarshal json to into HashRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func decodeValidateRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req validateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}
//encode hash and validate response for sending

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

//handle incoming hash request (call hash service method and build and return appropiate hash response request)
//WE never return a error from endpoint function
//take services as argument so we can send in any service

func MakeHashEndpoint(srv Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		//specify that request should be of type hashRequest
		//req recieve hashRequest strct
		req := request.(hashRequest)
		//send context and password that we got from hashRequest as argument to hash
		//v recience hash
		v, err := srv.Hash(ctx, req.Password)
		if err != nil {
			return hashResponse{v, err.Error()}, nil
		}
		return hashResponse{v, ""}, nil
	}
}

func MakeValidateEndpoint(srv Service) endpoint.Endpoint {

	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(validateRequest)
		v, err := srv.Validate(ctx, req.Password, req.Hash)
		if err != nil {
			return validateResponse{false, err.Error()}, nil
		}
		return validateResponse{v, ""}, nil
	}
}

// Endpoints represents all endpoints for the vault Service. this is because the error returned from endpoint is consided
//transport error so we will wrap error in endpoint struct
type Endpoints struct {
	HashEndpoint     endpoint.Endpoint
	ValidateEndpoint endpoint.Endpoint
}



// Hash uses the HashEndpoint to hash a password.
//build a request object,make a request,and parse the resulting response object into the normal argument to be returned

func (e Endpoints) Hash(ctx context.Context, password string) (string, error) {
	req := hashRequest{Password: password}
	resp, err := e.HashEndpoint(ctx, req)
	if err != nil {
		return "", err
	}
	hashResp := resp.(hashResponse)
	if hashResp.Err != "" {
		return "", errors.New(hashResp.Err)
	}
	return hashResp.Hash, nil
}

// Validate uses the ValidateEndpoint to validate a password and hash pair.
func (e Endpoints) Validate(ctx context.Context, password,
	hash string) (bool, error) {
	req := validateRequest{Password: password, Hash: hash}
	resp, err := e.ValidateEndpoint(ctx, req)
	if err != nil {
		return false, err
	}
	validateResp := resp.(validateResponse)
	if validateResp.Err != "" {
		return false, errors.New(validateResp.Err)
	}
	return validateResp.Valid, nil
}
