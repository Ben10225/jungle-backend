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

var availableCollection *mongo.Collection = configs.GetCollection(configs.DB, "available")

func GetWorkData(c *gin.Context, wg *sync.WaitGroup, thisMonth, nextMonth, theMonthAfterNext string) (structs.TimeTable, error) {
	defer wg.Done()

	data := structs.Available{
		Yymm: thisMonth,
	}

	query := bson.M{
		"yymm": bson.M{"$in": []string{thisMonth, nextMonth, theMonthAfterNext}},
	}
	if theMonthAfterNext == "" {
		query = bson.M{
			"yymm": bson.M{"$in": []string{thisMonth, nextMonth}},
		}
	}

	cursor, err := availableCollection.Find(c, query)
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
		for _, v := range createData {
			_, err := availableCollection.InsertOne(c, structs.Available{
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

			_, err := availableCollection.UpdateOne(context.TODO(), filter, update)
			if err != nil {
				er = fmt.Errorf("err: %v", err)
			}
		}
	}
	return er
}

func DeleteAllAvailableData(c *gin.Context) error {
	filter := bson.D{}
	_, err := availableCollection.DeleteMany(context.TODO(), filter)

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
		_, err := availableCollection.InsertOne(c, structs.Available{
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

func ReadAvailableData(c *gin.Context, reviseData structs.ReviseAvailable) (bool, []int) {
	var result structs.Available
	filter := bson.D{{Key: "yymm", Value: reviseData.Yymm}, {Key: "date", Value: reviseData.Date}}
	err := availableCollection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}

	originArr := result.WorkTime
	if reviseData.WholeHour > 1 {
		gap := reviseData.WholeHour
		for i, v := range originArr {
			if i > reviseData.HourIndex && gap > 0 {
				if v != 1 {
					return false, originArr
				}
				gap--
			}
			if i == reviseData.HourIndex {
				if v != 1 {
					return false, originArr
				}
				gap--
			}
		}
	}

	if originArr[reviseData.HourIndex] != 1 {
		return false, originArr
	}

	return true, originArr
}

func UpdateAvailableData(c *gin.Context, wg *sync.WaitGroup, reviseData structs.ReviseAvailable, originArr []int) error {
	defer wg.Done()
	filter := bson.D{{Key: "yymm", Value: reviseData.Yymm}, {Key: "date", Value: reviseData.Date}}

	gap := reviseData.WholeHour
	if gap > 1 {
		for i := range originArr {
			if gap > 1 && i > reviseData.HourIndex {
				originArr[i] = 0
				gap--
			}
			if i == reviseData.HourIndex {
				originArr[i] = 0
			}
		}
	}

	for i := range originArr {
		if i == reviseData.HourIndex {
			originArr[i] = 0
		}
	}

	update := bson.D{{Key: "$set", Value: bson.D{{Key: "workTime", Value: originArr}}}}
	_, err := availableCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	return err
}
