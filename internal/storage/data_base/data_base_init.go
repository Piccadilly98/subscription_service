package data_base

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type DataBase struct {
	db               *sql.DB
	dbCriticalLogger *log.Logger
}

func NewDB(connectStr string) (*DataBase, error) {

	dataBase := &DataBase{
		dbCriticalLogger: log.New(os.Stderr, "[DB CONNECTION ERROR] ", log.Ldate|log.Ltime),
	}
	db, err := sql.Open("postgres", connectStr)
	if err != nil {
		dataBase.dbCriticalLogger.Printf("CRITICAL: error in connect dataBase: %s\n", err.Error())
		return nil, fmt.Errorf("error in connect dataBase: %s\n", err.Error())
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		dataBase.dbCriticalLogger.Printf("CRITICAL: error in ping dataBase: %s\n", err.Error())
		return nil, fmt.Errorf("error in ping dataBase: %s\n", err.Error())
	}
	log.Println("Database connected successfully")
	dataBase.db = db
	return dataBase, nil
}

func (db *DataBase) Close() error {
	return db.db.Close()
}

func (db *DataBase) PingWithTimeout(duration time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()
	return db.db.PingContext(ctx)
}

func (db *DataBase) PingWithCtx(ctx context.Context) error {
	return db.db.PingContext(ctx)
}

func (db *DataBase) Begin() (*sql.Tx, error) {
	return db.db.Begin()
}
