package main

import (
	"log"
	"net/http"

	"github.com/JohanSJA/talks/lyricplayer"
	"github.com/gengo/grpc-gateway/runtime"
	"golang.org/x/net/context"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	err := player.RegisterPlayerHandlerFromEndpoint(ctx, mux, ":2015")
	if err != nil {
		log.Fatal("Couldn't register: ", err)
	}

	http.ListenAndServe(":2017", mux)
}
