package api

import (
	"database/sql"

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
	r.GET("/getNews", s.GetNews)
}
