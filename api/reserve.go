package api

import (
	"fmt"
	"jungle-proj/db/mongo"
	"jungle-proj/structs"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

var wg sync.WaitGroup

func (s *Server) GetAvailableData(c *gin.Context) {
	// admin with cookie res three months data

	thisMonth := c.Query("thisMonth")
	nextMonth := c.Query("nextMonth")
	theMonthAfterNext := c.Query("theMonthAfterNext")
	r := c.Query("r")

	result, err := mongo.GetWorkData(c, thisMonth, nextMonth, theMonthAfterNext)
	if err != nil {
		fmt.Println(err)
	}

	if r == "admin" {
		c.JSON(http.StatusOK, gin.H{
			"result": gin.H{
				"thisMonth":         result.ThisMonth,
				"nextMonth":         result.NextMonth,
				"theMonthAfterNext": result.TheMonthAfterNext,
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": gin.H{
			"thisMonth": result.ThisMonth,
			"nextMonth": result.NextMonth,
		},
	})

}

func (s *Server) PostAvailableData(c *gin.Context) {
	var req struct {
		Create []structs.Available
		Update []structs.Available
	}

	err := c.BindJSON(&req)
	if err != nil {
		log.Fatal(err)
	}

	create := req.Create
	update := req.Update

	err = mongo.PostWorkData(c, create, update)

	if err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"ok": true,
	})
}

func (s *Server) DeleteAvailable(c *gin.Context) {
	err := mongo.DeleteAllAvailableData(c)
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"ok": true,
	})
}

func (s *Server) CreateTestData(c *gin.Context) {
	err := mongo.CreateTestAvailableData(c)
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"ok": true,
	})
}

func (s *Server) PostReserveData(c *gin.Context) {
	var req struct {
		ReviseAvailable structs.ReviseAvailable
		AddReserve      structs.ReserveData
	}

	// need check time still have first

	err := c.BindJSON(&req)
	if err != nil {
		log.Fatal(err)
	}

	var errR, errA error

	wg.Add(2)
	go func() {
		errR = mongo.UpdateAvailableData(c, &wg, req.ReviseAvailable)
		errA = mongo.CreateReserveData(c, &wg, req.AddReserve)
	}()

	wg.Wait()
	if errR != nil && errA != nil {
		log.Fatal(errA, errR)
	}

	c.JSON(http.StatusOK, gin.H{
		"result": "ok",
	})
}
