package main
import (
	"fmt"
	"log"
	"net/http"
	"restfulAPI2/handler"
)

func main() {
	mux := handler.SetUpRouting()

	fmt.Println("http://localhost:3300")
	log.Fatal(http.ListenAndServe(":3300", mux))
}