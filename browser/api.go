package browser

import (
	"github.com/bowd/quip-exporter/scraper"
	"github.com/bowd/quip-exporter/types"
	"github.com/gin-gonic/gin"
	"strings"
)

func rootHandler(c *gin.Context) {
	currentUser, err := repo.GetCurrentUser(scraper.NewCurrentUserNode(c))
	if err != nil {
		_ = c.Error(err)
		return
	} else {
		c.JSON(200, currentUser)
	}
}

func foldersHandler(c *gin.Context) {
	folderIDs := strings.Split(c.Query("ids"), ",")
	response := make(map[string]*types.QuipFolder)
	for _, folderID := range folderIDs {
		node := scraper.NewFolderNode(c, "/", folderID)
		folder, _ := repo.GetFolder(node)
		response[folderID] = folder
	}
	c.JSON(200, response)
}

func threadsHandler(c *gin.Context) {
	threadIDs := strings.Split(c.Query("ids"), ",")
	response := make(map[string]*types.QuipThread)
	for _, threadID := range threadIDs {
		node := scraper.NewThreadNode(c, "/", threadID)
		thread, _ := repo.GetThread(node)
		response[threadID] = thread
	}
	c.JSON(200, response)
}

func commentsHandler(c *gin.Context) {
	threadId := c.Query("ids")
	thread := scraper.NewThreadNode(c, "/", threadId)
	node := scraper.NewThreadCommentsNode(thread.(*scraper.ThreadNode))
	comments, err := repo.GetThreadComments(node)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, comments)
}

func usersHandler(c *gin.Context) {
	userIDs := strings.Split(c.Query("ids"), ",")
	response := make(map[string]*types.QuipUser)
	for _, userID := range userIDs {
		node := scraper.NewUserNode(c, userID)
		user, _ := repo.GetUser(node)
		response[userID] = user
	}
	c.JSON(200, response)
}
