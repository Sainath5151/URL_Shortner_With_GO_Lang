package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type URL struct {
	ID           string    `json:"id"`
	OriginalURL  string    `json:"original_url"`
	ShortUrl     string    `json:"short_url"`
	Creationdate time.Time `json:"creation_date"`
}

var urlDB = make(map[string]URL)

func generateShortUrl(OriginalURL string) string {
	hasher := md5.New()
	hasher.Write([]byte(OriginalURL)) //Converting URl into bytes
	fmt.Println("hasher", hasher)
	data := hasher.Sum(nil)
	fmt.Println("hasher data", data)
	hash := hex.EncodeToString(data)
	fmt.Println("Hashed String", hash)
	fmt.Println("Final String", hash[:8])
	return hash[:8]

}
func createUrl(originalURL string) string {
	shortUrl := generateShortUrl(originalURL)
	id := shortUrl //Using short url as id for simplicity

	urlDB[id] = URL{
		ID:           id,
		OriginalURL:  originalURL,
		ShortUrl:     shortUrl,
		Creationdate: time.Now(),
	}
	return shortUrl
}

func getUrl(id string) (URL, error) {
	url, ok := urlDB[id]
	if !ok {
		return URL{}, errors.New("URL Not found")
	}
	return url, nil
}

func RootPageUrl(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello, world!")
}
func shortUrlHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(w, r)

	if r.Method == "OPTIONS" {
		return
	}
	var data struct {
		URL string `json:"url"`
	}
	err := json.NewDecoder(r.Body).Decode(&data)

	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	shortURL_ := createUrl(data.URL)
	// fmt.Fprintf(w, shortURL)
	response := struct {
		ShortURL string `json:"url"`
	}{ShortURL: shortURL_}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		fmt.Println("error faced during sending url", err)
	}
}

func redirectUrlHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/redirect/"):]
	url, err := getUrl(id)
	if err != nil {
		http.Error(w, "url not found", http.StatusNotFound)
		return
	}
	http.Redirect(w, r, url.OriginalURL, http.StatusFound)
}

func enableCors(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
}

func main() {

	http.HandleFunc("/shorten", shortUrlHandler)
	http.HandleFunc("/redirect/", redirectUrlHandler) // note trailing slash!
	http.HandleFunc("/", RootPageUrl)

	// Start the HTTP server on port 8080
	const port = 3000
	fmt.Println("The server is started on port ", port)

	err := http.ListenAndServe(":3000", nil)

	if err != nil {
		fmt.Println("Error Occured in while connecting server", err)
	}
}
