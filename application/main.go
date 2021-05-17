package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	config := SetupConfiguration()

	server := CreateServer(config)

	server.AddRoute("/", handleIndex)
	server.AddRoute("/audio", handleStream)

	UnhandledErrors(server.Run())
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	page := `
	<!doctype html>
	<html>
		<head>
			<meta charset="utf-8" />
			<title>Golang Audio Streaming Experiment</title>
		</head>
		<body>
			<audio controls preload="auto">
				<source src="/audio" type="audio/mpeg" />
			</audio>
		</body>
	</html>
	`

	w.Header().Add("Content-Type", "text/html")
	fmt.Fprint(w, page)
}

func handleStream(w http.ResponseWriter, r *http.Request) {
	fp, err := os.OpenFile("/application/data/testfile.mp3", os.O_RDONLY, 0755)

	if err != nil {
		log.Println(err)
		fmt.Fprintln(w, err.Error())
		return
	}

	defer fp.Close()

	stats, err := fp.Stat()

	if err != nil {
		log.Println(err)
		fmt.Fprintln(w, err.Error())
		return
	}

	var chunkSize int64 = 1024 * 1024
	flusher, ok := w.(http.Flusher)

	if !ok {
		fmt.Fprintln(w, "Failed to initialize chunked stream")
		return
	}

	w.Header().Add("Tranfer-Encoding", "chunked")
	w.Header().Add("Cache-Control", "no-store")

	var b []byte = make([]byte, chunkSize)

	for i := int64(0); i < stats.Size(); i = i + chunkSize {
		bytes, err := fp.Read(b)

		log.Printf("%d bytes read", bytes)

		if err != nil {
			log.Println(err)
			fmt.Fprintln(w, err.Error())
			return
		}

		w.Write(b)
		flusher.Flush()
	}
}
