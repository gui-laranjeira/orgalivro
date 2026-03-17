package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"orgalivro/backend/internal/config"
	"orgalivro/backend/internal/handler"
)

func New(
	cfg config.Config,
	ph *handler.ProfileHandler,
	bh *handler.BookHandler,
	eh *handler.EntryHandler,
	ih *handler.ISBNHandler,
) *gin.Engine {
	r := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = cfg.AllowedOrigins
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type"}
	r.Use(cors.New(corsConfig))

	v1 := r.Group("/api/v1")
	{
		v1.GET("/profiles", ph.List)
		v1.POST("/profiles", ph.Create)
		v1.DELETE("/profiles/:id", ph.Delete)

		v1.GET("/books", bh.List)
		v1.POST("/books", bh.Create)
		v1.GET("/books/:id", bh.Get)
		v1.PUT("/books/:id", bh.Update)
		v1.DELETE("/books/:id", bh.Delete)
		v1.PUT("/books/:id/owner", bh.TransferOwner)
		v1.GET("/authors", bh.Authors)
		v1.GET("/genres", bh.Genres)

		v1.GET("/profiles/:id/library", eh.List)
		v1.POST("/profiles/:id/library", eh.Add)
		v1.PUT("/profiles/:id/library/:book_id", eh.Update)
		v1.DELETE("/profiles/:id/library/:book_id", eh.Remove)

		v1.GET("/isbn/:isbn", ih.Lookup)
	}

	return r
}
