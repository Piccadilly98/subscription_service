package data_base

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type DataBase struct {
	db *sql.DB
}

func NewDB(connectStr string) (*DataBase, error) {
	db, err := sql.Open("postgres", connectStr)
	if err != nil {
		return nil, fmt.Errorf("error in connect dataBase: %s\n", err.Error())
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("error in ping dataBase: %s\n", err.Error())
	}
	log.Println("Database connected successfully")

	dataBase := &DataBase{
		db: db,
	}
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
