package main

import (
	"encoding/base32"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type longURL struct {
	URL string
}
type customURL struct {
	URL string
}

var mapping = map[string]string{}

func generateshortString(w http.ResponseWriter, r *http.Request) {
	// Write your logic here
	//	w.Write([]byte("Hello, HTTP!"))
	defer r.Body.Close()
	body, err1 := io.ReadAll(r.Body)
	if err1 != nil {
		w.WriteHeader(400)
		w.Write([]byte("INVALID BODY"))
	} else {
		fmt.Println(string(body))
	}

	// parse body into json
	var jsonData longURL
	err2 := json.Unmarshal(body, &jsonData)
	if err2 != nil {
		w.WriteHeader(400)
		w.Write([]byte("INVALID JSON"))
	} else {
		fmt.Println(jsonData)
	}

	longURL := jsonData.URL
	encoded := base32.StdEncoding.EncodeToString([]byte(longURL))
	//	fmt.Println(encoded[:8])
	mapping[encoded[:8]] = longURL
	w.Write([]byte(encoded[:8]))

}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	shortURL := r.URL.Path[1:]
	longURL, ok := mapping[shortURL]
	if ok {
		http.Redirect(w, r, longURL, http.StatusTemporaryRedirect)
	} else {
		http.Error(w, "Not Found", http.StatusNotFound)
	}
}

//	func customShortener(w http.ResponseWriter, r *http.Request) {
//		defer r.Body.Close()
//		body, err1 := io.ReadAll(r.Body)
//		if err1 != nil {
//			w.WriteHeader(400)
//			w.Write([]byte("INVALID BODY"))
//		} else {
//			fmt.Println(string(body))
//		}
//
//		var jsonData customURL
//		err2 := json.Unmarshal(body, &jsonData)
//		if err2 != nil {
//			w.WriteHeader(400)
//			w.Write([]byte("INVALID JSON"))
//		} else {
//			fmt.Println(jsonData)
//		}
//
//		mapping[customURL[:8]] = longURL
//
// }
func main() {
	http.HandleFunc("/shorten", generateshortString)
	//	http.HandleFunc("/custom", customShortener)
	http.HandleFunc("/", redirectHandler)
	fmt.Println("Server running on port 8080")
	http.ListenAndServe(":8080", nil)
}
