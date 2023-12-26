package api

import (
	"jungle-proj/db/mongo"
	"jungle-proj/structs"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) GetAdminData(c *gin.Context) {
	// admin with cookie res three months data

	thisMonth := c.Query("thisMonth")
	nextMonth := c.Query("nextMonth")
	theMonthAfterNext := c.Query("theMonthAfterNext")

	var errW, errB error

	var workTime structs.TimeTable
	var bookings structs.Booking

	wg.Add(2)
	go func() {
		workTime, errW = mongo.GetWorkData(c, &wg, thisMonth, nextMonth, theMonthAfterNext)
		bookings, errB = mongo.RealAllBooking(c, &wg, thisMonth, nextMonth, theMonthAfterNext)
	}()

	wg.Wait()

	if errW != nil && errB != nil {
		log.Fatal(errW, errB)
	}

	c.JSON(http.StatusOK, gin.H{
		"workTime": gin.H{
			"thisMonth":         workTime.ThisMonth,
			"nextMonth":         workTime.NextMonth,
			"theMonthAfterNext": workTime.TheMonthAfterNext,
		},
		"bookings": gin.H{
			"thisMonth":         bookings.ThisMonth,
			"nextMonth":         bookings.NextMonth,
			"theMonthAfterNext": bookings.TheMonthAfterNext,
		},
	})
}

func (s *Server) UpdateBookingState(c *gin.Context) {
	var req struct {
		Yymm      string
		Date      string
		HourIndex int
		NewState  int
	}

	err := c.BindJSON(&req)
	if err != nil {
		log.Fatal(err)
	}

	err = mongo.UpdateBookingStateData(c, req.Yymm, req.Date, req.HourIndex, req.NewState)
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"result": "ok",
	})
}
