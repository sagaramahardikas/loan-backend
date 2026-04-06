package main

import (
	"fmt"
	"net/http"
	"os"

	"example.com/loan/module/user/config"
)

func main() {
	mux := http.NewServeMux()
	err := RegisterGatewayHandler(mux)
	if err != nil {
		fmt.Printf("failed to register handler: %v\n", err)
		return
	}

	err = http.ListenAndServe("0.0.0.0:8888", mux)
	if err != nil {
		fmt.Printf("failed to start server: %v\n", err)
		os.Exit(1)
	}
}

func RegisterGatewayHandler(mux *http.ServeMux) error {
	err := config.RegisterUserGatewayHandler(mux)

	return err
}
