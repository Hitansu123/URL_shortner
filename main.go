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
	OriginalURl  string    `json:"OriginalURl"`
	NewURL       string    `json:"NewURL"`
	CreationDate time.Time `json:"CreationDate"`
}

var urlDB = make(map[string]URL)

func generateUrl(OriginalURl string) string {
	start := md5.New()
	start.Write([]byte(OriginalURl))
	date := start.Sum(nil)
	hash := hex.EncodeToString(date)
	return hash[:8]
}
func createURl(OriginalURl string) string {
	shortUrl := generateUrl(OriginalURl)
	urlDB[OriginalURl] = URL{
		OriginalURl:  OriginalURl,
		NewURL:       shortUrl,
		CreationDate: time.Now(),
	}
	return shortUrl
}
func getURl(OriginalURl string) (URL, error) {
	data, ok := urlDB[OriginalURl]
	if !ok {
		return URL{}, errors.New("URl not found")
	}
	return data, nil
}
func shortUrlHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		URL string `json:"url"`
	}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invlid request", http.StatusBadRequest)
		return
	}
	shortUrl_ := createURl(data.URL)
	responce := struct {
		Shorturl string `json:"shortUrl"`
	}{Shorturl: shortUrl_}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responce)
}
func redirectURlHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/redirect/"):]
	url, err := getURl(id)
	if err != nil {
		http.Error(w, "Invalid", http.StatusBadRequest)
	}
	http.Redirect(w, r, url.OriginalURl, http.StatusFound)
}

func main() {
	fmt.Println("starting the server")
	fmt.Println(generateUrl("https://github.com/Hitansu123"))

	http.HandleFunc("/shorter", shortUrlHandler)
	fmt.Println("starting the server.........")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		fmt.Println("Error in starting the server")
	}

}
