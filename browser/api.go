package browser

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/search/query"
	"github.com/bowd/quip-exporter/path"
	"github.com/bowd/quip-exporter/scraper"
	"github.com/bowd/quip-exporter/types"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"strings"
)

func rootHandler(c *gin.Context) {
	currentUser, err := repo.GetCurrentUser(scraper.NewCurrentUserNode(c, false))
	if err != nil {
		_ = c.Error(err)
		return
	} else {
		c.JSON(200, currentUser)
	}
}

func foldersHandler(c *gin.Context) {
	folderIDs := strings.Split(c.Query("ids"), ",")
	response := make(map[string]FolderResponse)
	for _, folderID := range folderIDs {
		folder := buildFolderResponse(c, folderID)
		if folder == nil {
			continue
		}
		response[folderID] = *folder
	}
	c.JSON(200, response)
}

func buildFolderResponse(ctx context.Context, folderID string) *FolderResponse {
	node := scraper.NewFolderNode(ctx, "/", folderID)
	folder, err := repo.GetFolder(node)
	if err != nil {
		return nil
	}
	folderPath, err := path.GetPathToFolder(ctx, folder, repo)
	if err != nil {
		return nil

	}
	return &FolderResponse{
		QuipFolder: folder,
		Path:       folderPath,
	}
}

func threadsHandler(c *gin.Context) {
	threadIDs := strings.Split(c.Query("ids"), ",")
	response := make(map[string]ThreadResponse)
	for _, threadID := range threadIDs {
		resp := buildThreadResponse(c, threadID)
		if resp == nil {
			continue
		}
		response[threadID] = *resp
	}
	c.JSON(200, response)
}

func buildThreadResponse(ctx context.Context, threadID string) *ThreadResponse {
	node := scraper.NewThreadNode(ctx, "/", threadID)
	thread, err := repo.GetThread(node)
	if err != nil {
		return nil
	}

	thread.InjectBlobHost(config.BlobHost)
	threadPath, err := path.GetPathToThread(ctx, thread, repo)
	if err != nil {
		return nil
	}
	return &ThreadResponse{
		QuipThread: thread,
		Path:       threadPath,
	}

}

func commentsHandler(c *gin.Context) {
	threadId := c.Query("threadId")
	thread := scraper.NewThreadNode(c, "/", threadId)
	node := scraper.NewThreadCommentsNode(thread.(*scraper.ThreadNode))
	comments, err := repo.GetThreadComments(node)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := CommentsResponse{
		Comments: comments,
		Users:    make(map[string]*types.QuipUser),
	}

	for _, comment := range comments {
		node := scraper.NewUserNode(c, comment.AuthorID)
		response.Users[comment.AuthorID], _ = repo.GetUser(node)
	}

	c.JSON(200, response)
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

func searchHandler(c *gin.Context) {
	requestBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		_ = c.Error(fmt.Errorf("error reading request body: %v", err))
		return
	}

	logger.Debugf("request body: %s", requestBody)

	// parse the request
	var searchRequest bleve.SearchRequest
	err = json.Unmarshal(requestBody, &searchRequest)
	if err != nil {
		_ = c.Error(fmt.Errorf("error parsing query: %v", err))
		return
	}

	logger.Debugf("parsed request %#v", searchRequest)

	// validate the query
	if srqv, ok := searchRequest.Query.(query.ValidatableQuery); ok {
		err = srqv.Validate()
		if err != nil {
			_ = c.Error(fmt.Errorf("error validating query: %v", err))
			return
		}
	}

	// execute the query
	results, err := search.Search(&searchRequest)
	if err != nil {
		_ = c.Error(fmt.Errorf("error executing query: %v", err))
		return
	}
	response := ThreadSearchResponse{
		Meta:    results,
		Threads: make(map[string]ThreadResponse),
		Folders: make(map[string]FolderResponse),
	}
	for _, hit := range results.Hits {
		thread := buildThreadResponse(c, hit.ID)
		if thread != nil {
			response.Threads[hit.ID] = *thread
		} else {
			folder := buildFolderResponse(c, hit.ID)
			if folder != nil {
				response.Folders[hit.ID] = *folder
			}
		}
	}
	c.JSON(200, response)
}
