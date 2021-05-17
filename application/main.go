package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func main() {
	config := SetupConfiguration()

	server := CreateServer(config)

	server.AddRoute("/", handleIndex)
	server.AddRoute("/audio", handleStream)

	UnhandledErrors(server.Run())
}

func handleIndex(a *ServerApplication) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fp, err := os.OpenFile("/application/templates/index.html", os.O_RDONLY, 0755)

		if err != nil {
			log.Println(err)
			w.WriteHeader(404)
			fmt.Fprintln(w, err.Error())
			return
		}

		defer fp.Close()

		stats, err := fp.Stat()

		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			fmt.Fprintln(w, err.Error())
			return
		}

		var page []byte = make([]byte, stats.Size())

		_, err = fp.Read(page)

		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			fmt.Fprintln(w, err.Error())
			return
		}

		w.Header().Add("Content-Type", "text/html")
		w.Write(page)
	}
}

func handleStream(a *ServerApplication) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fp, err := os.OpenFile("/application/data/testfile.mp3", os.O_RDONLY, 0755)

		if err != nil {
			log.Println(err)
			w.WriteHeader(400)
			fmt.Fprintln(w, err.Error())
			return
		}

		defer fp.Close()

		flusher, ok := w.(http.Flusher)

		if !ok {
			log.Println(err)
			w.WriteHeader(500)
			fmt.Fprintln(w, "Serverr Error Encountered")
			return
		}

		stats, err := fp.Stat()

		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			fmt.Fprintln(w, err.Error())
			return
		}

		var chunkSize int64 = a.config.ChunkSize

		var b []byte = make([]byte, chunkSize)

		// Read from Range
		rangeHeader := strings.Split(r.Header.Get("Range"), "=")
		var rbytes []string = make([]string, 2)

		if len(rangeHeader) != 2 {
			w.WriteHeader(400)
			fmt.Fprintln(w, "Range Header missing")
			return
		}

		rbytes = strings.Split(rangeHeader[1], "-")

		startByte, err := strconv.Atoi(rbytes[0])

		if err != nil {
			startByte = 0
		}

		log.Printf("INFO: Start Byte %d", startByte)
		bytes, err := fp.ReadAt(b, int64(startByte))

		log.Printf("%d bytes read", bytes)

		if err != nil {
			if !errors.Is(err, io.EOF) {
				log.Println(err)
				w.WriteHeader(500)
				fmt.Fprintln(w, err.Error())
				return
			}
		}

		w.Header().Add("Cache-Control", "no-cache")
		w.Header().Add("Accept-Ranges", "bytes")
		w.Header().Set("Content-Length", strconv.Itoa(bytes))
		w.Header().Set("Content-Range", fmt.Sprintf("bytes %s-%s/%s", strconv.Itoa(startByte), strconv.Itoa(startByte+bytes), strconv.Itoa(int(stats.Size()))))
		w.WriteHeader(206)
		w.Write(b[:bytes-1])
		flusher.Flush()
		log.Print("INFO: Chunk Transferred successfully")
	}
}
