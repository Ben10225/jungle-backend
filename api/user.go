package api

import (
	"jungle-proj/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (s *Server) createUser(ctx *gin.Context) {
	var user model.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := s.store.AddUser(&user); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// users, _ := s.store.GetAllUser()
	// for _, v := range users {
	// 	fmt.Println(v.Name)
	// }

	ctx.JSON(http.StatusOK, gin.H{
		"OK": true,
	})
}

func (s *Server) loginUser(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))

	user, err := s.store.GetUserByID(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "no user",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

func (s *Server) getInput(ctx *gin.Context) {
	var user model.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"name":     user.Name,
		"email":    user.Email,
		"password": "secret",
	})
}
