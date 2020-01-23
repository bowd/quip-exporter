package client

import (
	"fmt"
	"github.com/parnurzeal/gorequest"
	"strings"
)

func (qc *QuipClient) testToken() error {
	url := currentUserURL()
	resp, rawBody, _ := qc.getWithToken(url, Token{qc.token, 0}).End()
	if resp.StatusCode > 400 {
		return fmt.Errorf("%s", rawBody)
	}
	return nil
}

const (
	BASE_URL                    = "https://platform.quip.com/1"
	FOLDERS_MASK                = "/folders/?ids=%s"
	THREADS_MASK                = "/threads/?ids=%s"
	USERS_MASK                  = "/users/?ids=%s"
	FOLDER_MASK                 = "/folders/%s"
	EXPORT_MASK                 = "/threads/%s/export/%s"
	CURRENT_USER_PATH           = "/users/current"
	THREAD_COMMENTS_MASK        = "/messages/%s?count=100"
	THREAD_COMMENTS_CURSOR_MASK = "/messages/%s?count=100&max_created_usec=%d"
	BLOB_MASK                   = "/blob/%s/%s"
)

func batchURL(pathMask string, ids []string) string {
	mask := BASE_URL + pathMask
	idList := strings.Join(ids, ",")
	return fmt.Sprintf(mask, idList)
}

func currentUserURL() string {
	return BASE_URL + CURRENT_USER_PATH
}

func exportThreadURL(threadID string, exportType string) string {
	return BASE_URL + fmt.Sprintf(EXPORT_MASK, threadID, exportType)
}

func blobURL(threadID string, blobID string) string {
	return BASE_URL + fmt.Sprintf(BLOB_MASK, threadID, blobID)
}

func threadCommentsURL(threadID string, cursor *uint64) string {
	if cursor == nil {
		return BASE_URL + fmt.Sprintf(THREAD_COMMENTS_MASK, threadID)
	} else {
		return BASE_URL + fmt.Sprintf(THREAD_COMMENTS_CURSOR_MASK, threadID, *cursor)
	}
}

func (qc *QuipClient) getWithToken(url string, token Token) *gorequest.SuperAgent {
	qc.logger.Debugf("Requesting %s with %s", url, token)
	return gorequest.New().Get(url).Set("Authorization", "Bearer "+token.value)
}
