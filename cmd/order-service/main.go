package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mlucas4330/orderflow-pro/internal/db"
	"github.com/mlucas4330/orderflow-pro/internal/handler"
)

func main() {
	router := gin.Default()

	ctx := context.Background()

	dsn := fmt.Sprintf("postgres://%s:%s@%s:5432/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"),
	)

	dbpool, err := db.New(ctx, dsn)
	if err != nil {
		log.Fatalf("Falha ao conectar com o banco de dados: %v", err)
	}
	defer dbpool.Close()

	healthHandler := handler.New(dbpool)

	router.GET("/ping", healthHandler.Check)
	router.Run()
}
