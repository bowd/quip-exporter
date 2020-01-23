package browser

import "github.com/gin-gonic/gin"

type Config struct {
	Host string
	Port string
}

func Run(config Config) {
	r := gin.Default()
	r.LoadHTMLGlob("public/*.html")        // load the built dist path
	r.StaticFile("/", "public/index.html") // use the loaded source
	r.Static("/static/", "dist/")

	if err := r.Run(config.Host + ":" + config.Port); err != nil {
		panic(err)
	}
}
