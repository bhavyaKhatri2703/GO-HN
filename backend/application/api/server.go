package api

import (
	"database/sql"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

type Server struct {
	Router   *gin.Engine
	DB       *sql.DB
	GrpcConn *grpc.ClientConn
}

func NewServer(db *sql.DB, grpcConn *grpc.ClientConn) *Server {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	s := &Server{
		Router:   router,
		DB:       db,
		GrpcConn: grpcConn,
	}

	return s
}

func (s *Server) SetRoutes() {
	r := s.Router
	r.POST("/saveInterests", s.InterestsHandler)
	r.POST("/getNews", s.GetNews)
}
