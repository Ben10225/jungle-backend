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

var (
	wg sync.WaitGroup
	mu sync.Mutex
)

func (s *Server) GetAvailableData(c *gin.Context) {
	thisMonth := c.Query("thisMonth")
	nextMonth := c.Query("nextMonth")

	var result structs.TimeTable
	var err error

	wg.Add(1)
	go func() {
		result, err = mongo.GetWorkData(c, &wg, thisMonth, nextMonth, "")
	}()
	wg.Wait()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("here")
	data, err := s.Repo.FindByID(c, 11)
	if err != nil {
		fmt.Println("failed to find:", err)
		return
	}
	fmt.Println(data)

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

	err := c.BindJSON(&req)
	if err != nil {
		log.Fatal(err)
	}

	err = s.Repo.Insert(c, req.AddReserve)
	if err != nil {
		fmt.Println("failed to insert:", err)
		return
	}

	// lock
	mu.Lock()
	defer mu.Unlock()

	canReserve, originTime := mongo.ReadAvailableData(c, req.ReviseAvailable)

	if !canReserve {
		c.JSON(http.StatusOK, gin.H{
			"result": "time repeat",
		})
		return
	}

	var errR, errA error
	wg.Add(2)
	go func() {
		errR = mongo.UpdateAvailableData(c, &wg, req.ReviseAvailable, originTime)
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
