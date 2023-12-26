package mongo

import (
	"context"
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
