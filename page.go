package main

import (
	"io/ioutil"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Read the contents of the file.
		content, err := ioutil.ReadFile("output.txt")
		if err != nil {
			http.Error(w, "Could not read file", http.StatusInternalServerError)
			return
		}

		// Write the contents of the file to the response.
		w.Write(content)
	})

	// Start the web server.
	http.ListenAndServe(":8080", nil)
}
