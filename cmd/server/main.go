package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"tempcards/internal/share"
	"tempcards/internal/storage"
	"tempcards/internal/web"
)

func main() {
	path := flag.String("db", "cards.db", "database path")
	addr := flag.String("addr", ":8080", "listen address")
	flag.Parse()
	s, e := storage.Open(*path)
	if e != nil {
		log.Fatal(e)
	}
	defer s.Close()
	api := web.New(share.NewService(s))
	if e := http.ListenAndServe(*addr, api.Routes()); e != nil && !os.IsTimeout(e) {
		log.Fatal(e)
	}
}
