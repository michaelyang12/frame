package main

import (
	"bytes"
	"log"
	"net/http"
	"strconv"

	"github.com/h2non/bimg"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "Image Processing API is running"}`))
	})

	http.HandleFunc("/resize", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		file, _, err := r.FormFile("image")
		if err != nil {
			http.Error(w, "Failed to read image", http.StatusBadRequest)
			return
		}
		defer file.Close()

		buffer := new(bytes.Buffer)
		if _, err := buffer.ReadFrom(file); err != nil {
			http.Error(w, "Failed to read file", http.StatusInternalServerError)
			return
		}

		width, _ := strconv.Atoi(r.URL.Query().Get("width"))
		height, _ := strconv.Atoi(r.URL.Query().Get("height"))
		if width == 0 {
			width = 800
		}
		if height == 0 {
			height = 600
		}

		newImage, err := bimg.NewImage(buffer.Bytes()).Resize(width, height)
		if err != nil {
			http.Error(w, "Failed to process image", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "image/jpeg")
		w.Write(newImage)
	})

	log.Println("Server running on port 8080")
	http.ListenAndServe(":8080", nil)
}
