package main

import (
	"database/sql"
	"encoding/base32"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
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
	InsertIntoDB(encoded[:8], longURL)
	//	mapping[encoded[:8]] = longURL
	w.Write([]byte(encoded[:8]))

}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	shortURL := r.URL.Path[1:]
	//longURL, ok := mapping[shortURL]
	longURL, ok := QueryFromDB(shortURL)
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
func tryDBconnect() {
	db, err := sql.Open("sqlite3", "./main.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	//TODO: short urls can clash, so come up with a better db design
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS urls(
        short_url VARCHAR NOT NULL PRIMARY KEY ,
				long_url VARCHAR NOT NULL 

    );
	`

	result, err := db.Exec(sqlStmt)
	fmt.Println(result.RowsAffected())
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Table 'urls' created successfully!")

}
func InsertIntoDB(short_url string, long_url string) {
	db, err := sql.Open("sqlite3", "./main.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	sqlStmt := `
	INSERT INTO urls(short_url, long_url) values($1,$2)
	`

	_, err = db.Exec(sqlStmt, short_url, long_url)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("inserted values successfully")

}
func QueryFromDB(short_url string) (string, bool) {
	db, err := sql.Open("sqlite3", "./main.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	sqlStmt := `
	SELECT long_url from urls where short_url=$1;
	`

	rows, err := db.Query(sqlStmt, short_url)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var temp string
		err = rows.Scan(&temp)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("long_url: %s", temp)
		return temp, true
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	return "kahitri gandlay", false
}
func main() {
	tryDBconnect()
	//	InsertIntoDB("something", "somethinglong")
	//	QueryFromDB("something")
	http.HandleFunc("/shorten", generateshortString)
	//	http.HandleFunc("/custom", customShortener)
	http.HandleFunc("/", redirectHandler)
	fmt.Println("Server running on port 8080")
	http.ListenAndServe(":8080", nil)
}
