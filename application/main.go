package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	PORT := ":" + os.Getenv("PORT")
	http.HandleFunc("/", handleIndex)
	log.Println("Server Listening on :", PORT)
	if err := http.ListenAndServe(PORT, nil); err != nil {
		log.Fatal(err)
	}
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from Go!")
}
