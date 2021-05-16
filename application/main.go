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
	fp, err := os.OpenFile("/application/data/testfile.mp3", os.O_RDONLY, 0755)

	if err != nil {
		log.Println(err)
		fmt.Fprintf(w, err.Error())
		return
	}

	defer fp.Close()

	stats, err := fp.Stat()

	if err != nil {
		log.Println(err)
		fmt.Fprintf(w, err.Error())
		return
	}

	var b []byte = make([]byte, stats.Size())

	bytes, err := fp.Read(b)

	log.Printf("%d bytes read", bytes)

	if err != nil {
		log.Println(err)
		fmt.Fprintf(w, err.Error())
		return
	}

	w.Write(b)
}
