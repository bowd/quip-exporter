package browser

import "github.com/gin-gonic/gin"

type Config struct {
	Host string
	Port string
}

func Run(config Config) {
	r := gin.Default()
	r.LoadHTMLGlob("frontend/dist/*.html")        // load the built dist path
	r.StaticFile("/", "frontend/dist/index.html") // use the loaded source
	r.Static("/static", "frontend/dist/")

	if err := r.Run(config.Host + ":" + config.Port); err != nil {
		panic(err)
	}
}
