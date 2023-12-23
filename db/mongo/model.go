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
var reserveCollection *mongo.Collection = configs.GetCollection(configs.DB, "reserve")

func GetWorkData(c *gin.Context, thisMonth, nextMonth, theMonthAfterNext string) (structs.TimeTable, error) {
	data := structs.Available{
		Yymm: thisMonth,
	}

	query := bson.M{
		"yymm": bson.M{"$in": []string{thisMonth, nextMonth, theMonthAfterNext}},
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
		// newCreateDatas := []structs.Available{}
		for _, v := range createData {
			// newCreateDatas = append(newCreateDatas, v)
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

func UpdateAvailableData(c *gin.Context, wg *sync.WaitGroup, reviseData structs.ReviseAvailable) error {
	defer wg.Done()
	filter := bson.D{{Key: "yymm", Value: reviseData.Yymm}, {Key: "date", Value: reviseData.Date}}

	var availeData structs.Available
	err := availableCollection.FindOne(context.TODO(), filter).Decode(&availeData)
	if err != nil {
		log.Fatal(err)
	}

	gap := reviseData.WholeHour
	if gap > 1 {
		for i := range availeData.WorkTime {
			if gap > 1 && i > reviseData.HourIndex {
				availeData.WorkTime[i] = 0
				gap--
			}
			if i == reviseData.HourIndex {
				availeData.WorkTime[i] = 0
			}
		}
	}

	for i := range availeData.WorkTime {
		if i == reviseData.HourIndex {
			availeData.WorkTime[i] = 0
		}
	}

	update := bson.D{{Key: "$set", Value: bson.D{{Key: "workTime", Value: availeData.WorkTime}}}}
	_, err = availableCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	return err
}

func CreateReserveData(c *gin.Context, wg *sync.WaitGroup, addData structs.ReserveData) error {
	defer wg.Done()
	_, err := reserveCollection.InsertOne(c, structs.ReserveData{
		Titles:    addData.Titles,
		Time:      addData.Time,
		Cost:      addData.Cost,
		HourIndex: addData.HourIndex,
		WholeHour: addData.WholeHour,
	})
	if err != nil {
		return err
	}

	return nil
}
