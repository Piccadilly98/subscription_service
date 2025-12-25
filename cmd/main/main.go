package main

import (
	"fmt"
	"time"

	"github.com/Piccadilly98/subscription_service/internal/models/dto"
	"github.com/Piccadilly98/subscription_service/internal/models/entities_data_base"
)

func main() {
	// db, err := data_base.NewDB("")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// end := "27-03-2030"
	// dt := &dto.RequestCreateSubscriptions{
	// 	UserID:      "60601fee-2bf1-4721-ae6f-7636e79a0cba",
	// 	StartDate:   "26-02-2006",
	// 	EndDate:     &end,
	// 	ServiceName: "ser",
	// 	Price:       100,
	// }
	// en, err := dt.ToEntitie()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// ctx, _ := context.WithCancel(context.Background())
	// fmt.Println(db.CreateNewSubscribe(ctx, en))
	en := &entities_data_base.ReadSubscription{
		StartDate: time.Date(2025, 02, 02, 0, 0, 0, 0, time.UTC),
		EndDate:   nil,
	}

	d := dto.FromEntity(en)
	fmt.Println(d.Status)
}

func GetTimeDate(t time.Time) *time.Time {
	return &t
}
