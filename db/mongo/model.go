package mongo

import (
	"context"
	"fmt"
	"jungle-proj/db/configs"
	"jungle-proj/structs"
	"log"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "available")

func GetWorkData(c *gin.Context, thisMonth, nextMonth, theMonthAfterNext string) (structs.TimeTable, error) {
	data := structs.Available{
		Yymm: thisMonth,
	}

	query := bson.M{
		"yymm": bson.M{"$in": []string{thisMonth, nextMonth, theMonthAfterNext}},
	}

	cursor, err := userCollection.Find(c, query)
	if err != nil {
		log.Fatal(err)
	}

	defer cursor.Close(context.TODO())

	var result []structs.Available
	var newReturnSlice []structs.Available
	var intSlice [][]int

	for cursor.Next(context.TODO()) {
		err := cursor.Decode(&data)
		if err != nil {
			log.Fatal(err)
		}

		result = append(result, data)

		myslice := append([]int{}, data.WorkTime...)
		intSlice = append(intSlice, myslice)
	}

	for i := range result {
		newReturnSlice = append(newReturnSlice, structs.Available{Date: result[i].Date, Yymm: result[i].Yymm, WorkTime: intSlice[i]})
	}

	var timeTable structs.TimeTable

	for i, v := range newReturnSlice {
		if v.Yymm == thisMonth {
			timeTable.ThisMonth = append(timeTable.ThisMonth, newReturnSlice[i])
		} else if v.Yymm == nextMonth {
			timeTable.NextMonth = append(timeTable.NextMonth, newReturnSlice[i])
		} else if v.Yymm == theMonthAfterNext {
			timeTable.TheMonthAfterNext = append(timeTable.TheMonthAfterNext, newReturnSlice[i])
		}
	}
	return timeTable, nil
}

func PostWorkData(c *gin.Context, createData, updateData []structs.Available) error {
	var er error

	if len(createData) > 0 {
		// newCreateDatas := []structs.Available{}
		for _, v := range createData {
			// newCreateDatas = append(newCreateDatas, v)
			_, err := userCollection.InsertOne(c, structs.Available{
				Yymm:     v.Yymm,
				Date:     v.Date,
				WorkTime: v.WorkTime,
			})
			if err != nil {
				er = fmt.Errorf("err: %v", err)
			}
		}
	}

	if len(updateData) > 0 {
		for _, v := range updateData {
			filter := bson.D{{Key: "yymm", Value: v.Yymm}, {Key: "date", Value: v.Date}}
			update := bson.D{{Key: "$set", Value: bson.D{{Key: "workTime", Value: v.WorkTime}}}}

			_, err := userCollection.UpdateOne(context.TODO(), filter, update)
			if err != nil {
				er = fmt.Errorf("err: %v", err)
			}
		}
	}
	return er
}

func DeleteAllAvailableData(c *gin.Context) error {
	filter := bson.D{}
	_, err := userCollection.DeleteMany(context.TODO(), filter)

	if err != nil {
		return fmt.Errorf("err: %v", err)
	}
	return err
}

func CreateTestAvailableData(c *gin.Context) error {
	var err error
	createData := []structs.Available{
		{Yymm: "2023-12", Date: "20", WorkTime: []int{-1, -1, -1, -1, -1, -1, -1, -1, -1, 1, 1, 1}},
		{Yymm: "2023-12", Date: "31", WorkTime: []int{-1, 1, 1, -1, -1, -1, -1, -1, -1, 1, 1, 1}},
		{Yymm: "2024-1", Date: "31", WorkTime: []int{-1, -1, -1, -1, -1, -1, 1, 1, 1, -1, -1, -1}},
	}

	for _, v := range createData {
		_, err := userCollection.InsertOne(c, structs.Available{
			Yymm:     v.Yymm,
			Date:     v.Date,
			WorkTime: v.WorkTime,
		})
		if err != nil {
			return fmt.Errorf("err: %v", err)
		}
	}
	return err
}
