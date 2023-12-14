package api

import (
	"fmt"
	"jungle-proj/db/mongo"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) GetGameResult(c *gin.Context) {

	var req struct {
		ThisMonth string
		NextMonth string
	}

	err := c.BindJSON(&req)
	if err != nil {
		log.Fatal(err)
	}
	thisMonth := req.ThisMonth
	nextMonth := req.NextMonth

	fmt.Println(thisMonth, nextMonth)

	result, err := mongo.GetDateData(c, thisMonth, nextMonth)

	if err != nil {
		fmt.Println(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"result": gin.H{
			"thisMonth": result.ThisMonth,
			"nextMonth": result.NextMonth,
		},
	})
}
