package main

import (
	"log"
	"net/http"
	handler "yandex-practicum-go-shortener/internal/short/handlers"
	"yandex-practicum-go-shortener/internal/short/repo"
	"yandex-practicum-go-shortener/internal/short/usecase"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
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
