package browser

import (
	"github.com/bowd/quip-exporter/interfaces"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

type Config struct {
	Host string
	Port string
}

var repo interfaces.IRepository

func Run(config Config, _repo interfaces.IRepository) {
	repo = _repo
	r := gin.Default()
	// r.LoadHTMLGlob("frontend/dist/*.html")        // load the built dist path
	// r.StaticFile("/", "frontend/dist/index.html") // use the loaded source
	// r.Static("/static", "frontend/dist/")
	r.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"http://localhost:3000"},
		AllowMethods:  []string{"PUT", "PATCH", "GET", "POST"},
		AllowHeaders:  []string{"Origin", "Content-Type"},
		ExposeHeaders: []string{"Content-Length"},
		MaxAge:        12 * time.Hour,
	}))

	api := r.Group("/api")
	api.GET("/root", rootHandler)
	api.GET("/folders", foldersHandler)
	api.GET("/threads", threadsHandler)
	api.GET("/threads/comments", commentsHandler)
	api.GET("/users", usersHandler)

	if err := r.Run(config.Host + ":" + config.Port); err != nil {
		panic(err)
	}
}
