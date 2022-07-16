package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Dakak-Takto/yandex-practicum-go-shortener/internal/storage/inmem"
	"github.com/gorilla/securecookie"
)

func Example_postHandler() {
	const (
		addr                 string = "localhost:55555"
		baseURL              string = `http://` + addr
		secureCookieHashKey  string = `secret`
		secureCookieBlockKey string = `secret`
	)

	//init storage
	store, err := inmem.New()
	if err != nil {
		log.Fatal(err)
	}

	//init secureCookie
	securecookie := securecookie.New([]byte(secureCookieHashKey), []byte(secureCookieBlockKey))

	//create new app
	app := New(
		WithStorage(store),
		WithBaseURL(baseURL),
		WithAddr(addr),
		WithSecureCookie(securecookie),
	)

	//run app
	go app.Run()

	// payload for sending to server
	requestPayload := struct {
		URL string `json:"url"`
	}{
		URL: `http://original.url/example/?foo=bar`,
	}

	// make json
	requestPayloadJSON, err := json.Marshal(&requestPayload)
	if err != nil {
		log.Fatal(err)
	}

	// init http client
	client := http.Client{}

	// make request
	req, err := http.NewRequest(http.MethodPost, baseURL+`/api/shorten`, bytes.NewBuffer(requestPayloadJSON))
	if err != nil {
		log.Fatal(err)
	}

	// send request
	response, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		log.Fatal("success code must be 201")
	}

	var responsePayload struct {
		Result string `json:"result"`
	}
	if err := json.NewDecoder(response.Body).Decode(&responsePayload); err != nil {
		log.Fatal(err)
	}

	fmt.Println(responsePayload)
}
