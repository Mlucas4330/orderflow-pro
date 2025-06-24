package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/mlucas4330/orderflow-pro/internal/db"
)

func main() {
	ctx := context.Background()

	dsn := fmt.Sprintf("postgres://%s:%s@%s:5432/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"),
	)

	conn, err := db.New(ctx, dsn)
	if err != nil {
		log.Fatalf("Falha ao conectar com o banco de dados: %v", err)
	}
	defer conn.Close(ctx)

	err = conn.Ping(ctx)
	if err != nil {
		log.Fatalf("Não foi possível fazer o ping no banco de dados: %v", err)
	}

	log.Println("Conexão com o banco de dados estabelecida com sucesso!")
}
