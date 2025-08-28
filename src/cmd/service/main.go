package main

import (
	"fmt"
	"myapp/src/internal/database"
	"myapp/src/internal/httpHandlers"
	"myapp/src/internal/memory"
	"myapp/src/internal/order_model"
	"net/http"
	"os"
	"time"

	//"myapp/src/testutils"
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
    "myapp/src/internal/kafka"

	//"gorm.io/gorm/clause"
	"log"
)

func main() {
    memory:= createMemory()
    go kafka_consumer.RunConsumer(memory)
    handler:= handlers.NewOrderHandler(memory)
    r := mux.NewRouter()
    r.HandleFunc("/orders/{id}",handler.HandleOrder).Methods("GET")
    r.HandleFunc("/orders", handler.ShowOrderSearchForm).Methods("GET")
    err:= http.ListenAndServe(":8081",r)
    if err != nil {
        log.Fatalf("%v",err)
    }
}


func createMemory() *memory.Memory{
    time.Sleep(5*time.Second)
    dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&model.Order{},&model.Item{})
    repo:= database.NewOrderRepository(db)
    memory:= memory.NewMemory(repo)
    log.Println("Database is ready")
    return memory
}
