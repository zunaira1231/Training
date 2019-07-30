package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"vault"
	"vault/pb"
)

func main() {
	//ports are endpoint
	var (
		httpAddr = flag.String("http", ":3030", "http listen address")
		gRPCAddr = flag.String("grpc", ":3031", "gRPC listen address")
	)
	flag.Parse()
	//non nil empty context(no cancelllation or deadline specified and contain no value  )
	//base context of all services
	ctx := context.Background()
	srv := vault.NewService()
	//zero buffered channel for errors
	errChan := make(chan error)

	//trap termination signals and send error down errChan

	go func() {
		c := make(chan os.Signal, 1)
		//notify us when their is a SIGINT and SIGTERM signal
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		//send signal as string down the channel resulting in program termination
		errChan <- fmt.Errorf("%s", <-c)
		//if there is a error program will terminate
	}()

//MakeHashEndpoint return hashResponse(hash)
	hashEndpoint := vault.MakeHashEndpoint(srv)
	validateEndpoint := vault.MakeValidateEndpoint(srv)

	endpoints := vault.Endpoints{
		HashEndpoint:     hashEndpoint,
		ValidateEndpoint: validateEndpoint,
	}

	// HTTP transport
	go func() {
		log.Println("http:", *httpAddr)
		handler := vault.NewHTTPServer(ctx, endpoints)
		errChan <- http.ListenAndServe(*httpAddr, handler)
	}()

	// gRPC transport
	go func() {
		//create low level TCP network listener and serve gRPC over that
		listener, err := net.Listen("tcp", *gRPCAddr)
		if err != nil {
			errChan <- err
			return
		}
		log.Println("grpc:", *gRPCAddr)
		handler := vault.NewGRPCServer(ctx, endpoints)
		gRPCServer := grpc.NewServer()
		pb.RegisterVaultServer(gRPCServer, handler)
		errChan <- gRPCServer.Serve(listener)
	}()

	log.Fatalln(<-errChan)
}
