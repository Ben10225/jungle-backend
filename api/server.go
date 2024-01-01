package api

import (
	"embed"
	"fmt"
	"mime"
	"strings"

	sqlc "jungle-proj/db/sqlc"
	"jungle-proj/repository/order"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

//go:embed web/dist
var staticFS embed.FS

type Server struct {
	store  *sqlc.Store
	router *gin.Engine
	rdb    *redis.Client
	Repo   *order.RedisRepo
}

func NewServer(store *sqlc.Store) (*Server, error) {
	server := &Server{
		store: store,
		rdb:   redis.NewClient(&redis.Options{}),
		Repo: &order.RedisRepo{
			Client: redis.NewClient(&redis.Options{}),
		},
	}

	err := server.rdb.Ping(&gin.Context{}).Err()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Connect Redis")

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.Use(corsMiddleware())

	router.POST("/login", server.loginUser)
	router.POST("/user", server.signUp)

	router.GET("/admin", server.GetAdminData)
	router.GET("/available", server.GetAvailableData)
	router.POST("/available", server.PostAvailableData)

	router.POST("/reserve", server.PostReserveData)
	router.PATCH("/bookingState", server.UpdateBookingState)

	router.GET("createTest", server.CreateTestData)
	router.DELETE("/availableAll", server.DeleteAvailable)

	router.GET("createReserveTest", server.CreateReserveTestData)
	router.DELETE("/reserveAll", server.DeleteReserve)

	router.NoRoute(EmbedReact)

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(fmt.Sprintf(":%v", address)) // fmt.Sprintf(":%v", address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func EmbedReact(c *gin.Context) {
	path := c.Request.URL.Path                                     // 获取请求路径
	s := strings.Split(path, ".")                                  // 分割路径，获取文件后缀
	prefix := "web/dist"                                           // 前缀路径
	if data, err := staticFS.ReadFile(prefix + path); err != nil { // 读取文件内容
		// 如果文件不存在，返回首页 index.html
		if data, err = staticFS.ReadFile(prefix + "/index.html"); err != nil {
			c.JSON(404, gin.H{
				"err": err,
			})
		} else {
			c.Data(200, mime.TypeByExtension(".html"), data)
		}
	} else {
		// 如果文件存在，根据请求的文件后缀，设置正确的mime type，并返回文件内容
		c.Data(200, mime.TypeByExtension(fmt.Sprintf(".%s", s[len(s)-1])), data)
	}
}
