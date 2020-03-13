package browser

import (
	"github.com/blevesearch/bleve"
	bleveHttp "github.com/blevesearch/bleve/http"
	"github.com/bowd/quip-exporter/interfaces"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"time"
)

type Config struct {
	Host     string
	Port     string
	BlobHost string
}

var config Config
var repo interfaces.IRepository
var search bleve.Index
var logger = logrus.WithField("module", "api")

func Run(_config Config, _repo interfaces.IRepository, _search bleve.Index) {
	config = _config
	repo = _repo
	search = _search
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

	r.Static("/blob", "./output/data/blobs")
	r.Static("/profile", "./output/data/profile_pictures")

	bleveHttp.RegisterIndexName("threads", search)
	// searchHandler := bleveHttp.NewSearchHandler("threads")

	api := r.Group("/api")
	api.POST("/search", searchHandler)
	api.GET("/root", rootHandler)
	api.GET("/folders", foldersHandler)
	api.GET("/threads", threadsHandler)
	api.GET("/threads/comments", commentsHandler)
	api.GET("/users", usersHandler)

	if err := r.Run(config.Host + ":" + config.Port); err != nil {
		panic(err)
	}
}
