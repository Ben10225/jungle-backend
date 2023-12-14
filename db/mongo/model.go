package mongo

import (
	"context"
	"jungle-proj/db/configs"
	"jungle-proj/structs"
	"log"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "available")

func GetDateData(c *gin.Context, thisMonth, nextMonth string) (structs.TimeTable, error) {
	// func GetDateData(c *gin.Context, yymm string) (structs.TimeTable, error) {
	data := structs.Available{
		Yymm: thisMonth,
	}

	query := bson.M{
		"yymm": bson.M{"$in": []string{thisMonth, nextMonth}},
	}

	// cursor, err := userCollection.Find(c, bson.M{"yymm": yymm})
	cursor, err := userCollection.Find(c, query)
	if err != nil {
		log.Fatal(err)
	}

	defer cursor.Close(context.TODO())

	var result []structs.Available
	var newReturnSlice []structs.Available
	var boolSlice [][]bool

	for cursor.Next(context.TODO()) {
		err := cursor.Decode(&data)
		if err != nil {
			log.Fatal(err)
		}

		result = append(result, data)

		myslice := append([]bool{}, data.SureTimeArray...)
		boolSlice = append(boolSlice, myslice)
	}

	for i := range result {
		newReturnSlice = append(newReturnSlice, structs.Available{Date: result[i].Date, Yymm: result[i].Yymm, SureTimeArray: boolSlice[i]})
	}

	var timeTable structs.TimeTable

	for i, v := range newReturnSlice {
		if v.Yymm == thisMonth {
			timeTable.ThisMonth = append(timeTable.ThisMonth, newReturnSlice[i])
		} else {
			timeTable.NextMonth = append(timeTable.NextMonth, newReturnSlice[i])
		}
	}
	return timeTable, nil
}

// func GetDateData(c *gin.Context, yymm string) (*structs.Available, error) {
// 	data := structs.Available{
// 		Yymm: yymm,
// 	}
// 	err := userCollection.FindOne(c, bson.M{"yymm": yymm}).Decode(&data)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &data, nil
// }
