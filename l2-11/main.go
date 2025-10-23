package main

import (
	"calendar-server/internal/handler"
	"calendar-server/internal/repo"
	"calendar-server/internal/service"
	"calendar-server/internal/utils"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log"
	"net/http"
)

func main() {
	connString := "postgres://postgres:postgres@localhost:5432/calendar"
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		log.Fatal("Ошибка при подключении к БД", err)
	}

	calendarRepository := repo.New(conn)
	calendarService := service.New(calendarRepository)
	calendarUtils := utils.New()
	calendarHandler := handler.New(calendarService, calendarUtils)

	http.HandleFunc("/create_event", calendarHandler.CreateEvent)

	fmt.Println("Сервер запущен")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
}
