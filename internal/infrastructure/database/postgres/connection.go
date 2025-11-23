package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func ConnectionDatabase() *sql.DB {
	// TODO : load envs with specific file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error is occurred  on .env file please check")
	}

	host := os.Getenv("DB_HOST")
	port, _ := strconv.Atoi(os.Getenv("DB_PORT"))
	user := os.Getenv("DB_USER")
	dbname := os.Getenv("DB_NAME")
	pass := os.Getenv("DB_PASSWORD")

	postgreslSetup := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		host, port, user, dbname, pass)
	fmt.Printf("url montada %v", postgreslSetup)
	db, err := sql.Open("postgres", postgreslSetup)
	if err != nil {
		fmt.Println("There is an error while connecting to the database ", err)
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}
