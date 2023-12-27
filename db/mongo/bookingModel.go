package mongo

import (
	"context"
	"fmt"
	"jungle-proj/db/configs"
	"jungle-proj/structs"
	"log"
	"sync"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var reserveCollection *mongo.Collection = configs.GetCollection(configs.DB, "reserve")

func CreateReserveData(c *gin.Context, wg *sync.WaitGroup, addData structs.ReserveData) error {
	defer wg.Done()
	_, err := reserveCollection.InsertOne(c, structs.ReserveData{
		Titles: addData.Titles,
		Detail: addData.Detail,
		Hour:   addData.Hour,
		Yymm:   addData.Yymm,
		Date:   addData.Date,
		User:   addData.User,
	})
	if err != nil {
		return err
	}

	return nil
}

func RealAllBooking(c *gin.Context, wg *sync.WaitGroup, thisMonth, nextMonth, theMonthAfterNext string) (structs.Booking, error) {
	defer wg.Done()
	data := structs.ReserveData{
		Yymm: thisMonth,
	}

	query := bson.M{
		"yymm": bson.M{"$in": []string{thisMonth, nextMonth, theMonthAfterNext}},
	}

	cursor, err := reserveCollection.Find(c, query)
	if err != nil {
		log.Fatal(err)
	}

	defer cursor.Close(context.TODO())

	var result []structs.ReserveData
	var newReturnSlice []structs.ReserveData
	var titleSlice [][]string

	for cursor.Next(context.TODO()) {
		err := cursor.Decode(&data)
		if err != nil {
			log.Fatal(err)
		}

		result = append(result, data)

		myslice := append([]string{}, data.Titles...)
		titleSlice = append(titleSlice, myslice)
	}

	for i := range result {
		newReturnSlice = append(newReturnSlice, structs.ReserveData{Date: result[i].Date, Yymm: result[i].Yymm, Titles: titleSlice[i], Detail: result[i].Detail, Hour: result[i].Hour, User: result[i].User})
	}

	var booking structs.Booking

	for i, v := range newReturnSlice {
		if v.Yymm == thisMonth {
			booking.ThisMonth = append(booking.ThisMonth, newReturnSlice[i])
		} else if v.Yymm == nextMonth {
			booking.NextMonth = append(booking.NextMonth, newReturnSlice[i])
		} else if v.Yymm == theMonthAfterNext {
			booking.TheMonthAfterNext = append(booking.TheMonthAfterNext, newReturnSlice[i])
		}
	}

	return booking, nil
}

func UpdateBookingStateData(c *gin.Context, yymm, date string, hourIndex, newValue int) error {
	filter := bson.D{{Key: "yymm", Value: yymm}, {Key: "date", Value: date}, {Key: "hour.index", Value: hourIndex}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "detail.state", Value: newValue}}}}

	_, err := reserveCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}

	return err
}

func DeleteAllReserveData(c *gin.Context) error {
	filter := bson.D{}
	_, err := reserveCollection.DeleteMany(context.TODO(), filter)

	if err != nil {
		return fmt.Errorf("err: %v", err)
	}
	return err
}

func CreateTestReserveData(c *gin.Context) error {
	var err error
	createData := []structs.ReserveData{
		{Yymm: "2023-12", Date: "14", Titles: []string{"修眉毛", "修鬍子", "修手毛", "修腳毛", "修私密處"}, User: structs.BookingUser{Name: "彭彭", Phone: "0912345678"}, Hour: structs.BookingHour{Index: 4, Whole: 3}, Detail: structs.BookingDetails{Time: "2 小時 30 分鐘", Cost: "7500", State: 0}},
		{Yymm: "2023-12", Date: "14", Titles: []string{"修眉毛"}, User: structs.BookingUser{Name: "彭彭", Phone: "0912345678"}, Hour: structs.BookingHour{Index: 10, Whole: 1}, Detail: structs.BookingDetails{Time: "20 分鐘", Cost: "1000", State: 0}},
		{Yymm: "2023-12", Date: "14", Titles: []string{"修手毛", "修私密處"}, User: structs.BookingUser{Name: "彭彭", Phone: "0912345678"}, Hour: structs.BookingHour{Index: 8, Whole: 2}, Detail: structs.BookingDetails{Time: "1 小時 30 分鐘", Cost: "4500", State: 0}},
		{Yymm: "2023-12", Date: "14", Titles: []string{"修眉毛"}, User: structs.BookingUser{Name: "彭彭", Phone: "0912345678"}, Hour: structs.BookingHour{Index: 0, Whole: 1}, Detail: structs.BookingDetails{Time: "20 分鐘", Cost: "1000", State: 0}},
	}

	for _, v := range createData {
		_, err := reserveCollection.InsertOne(c, structs.ReserveData{
			Yymm:   v.Yymm,
			Date:   v.Date,
			Titles: v.Titles,
			Detail: v.Detail,
			Hour:   v.Hour,
			User:   v.User,
		})
		if err != nil {
			return fmt.Errorf("err: %v", err)
		}
	}
	return err
}
