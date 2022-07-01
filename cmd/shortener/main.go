package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"

	handler "github.com/Dakak-Takto/yandex-practicum-go-shortener/internal/short/handlers"
	"github.com/Dakak-Takto/yandex-practicum-go-shortener/internal/short/repo"
	"github.com/Dakak-Takto/yandex-practicum-go-shortener/internal/short/usecase"
)

func main() {

	storage, err := repo.NewPostgresRepository("postgres://postgres:postgres@localhost/praktikum?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	usecase := usecase.New(storage)

	cookieStore := sessions.NewCookieStore([]byte("secret"))
	handler := handler.New(usecase, cookieStore)

	router := chi.NewMux()
	handler.Register(router)

	http.ListenAndServe("localhost:8080", router)

}
