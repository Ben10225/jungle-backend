package api

import (
	"database/sql"
	sqlc "jungle-proj/db/sqlc"
	"jungle-proj/util"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (s *Server) signUp(ctx *gin.Context) {
	var user sqlc.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err := s.store.GetUserByEmail(ctx, user.Email)
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"OK":      false,
			"message": "此信箱已被註冊",
		})
		return
	}

	arg := sqlc.CreateUserParams{
		Uuid:       util.UuidGenerate(),
		Name:       user.Name,
		Email:      user.Email,
		Password:   user.Password,
		CreateTime: sql.NullTime{Time: time.Now().Add(time.Hour * 8), Valid: true},
	}

	_, err = s.store.CreateUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"OK": true,
	})
}

func (s *Server) loginUser(ctx *gin.Context) {
	// id, _ := strconv.Atoi(ctx.Param("id"))

	var user sqlc.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := sqlc.GetUserByEmailAndPwdParams{
		Email:    user.Email,
		Password: user.Password,
	}

	res, err := s.store.GetUserByEmailAndPwd(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// todo add token

	ctx.JSON(http.StatusOK, gin.H{
		"user": &res,
	})
}
