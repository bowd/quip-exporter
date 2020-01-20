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
	BASE_URL     = "https://platform.quip.com/1"
	FOLDERS_MASK = "/folders/?ids=%s"
	THREADS_MASK = "/threads/?ids=%s"
	USERS_MASK   = "/users/?ids=%s"
	FOLDER_MASK  = "/folders/%s"
)

func batchURL(pathMask string, ids []string) string {
	mask := BASE_URL + pathMask
	idList := strings.Join(ids, ",")
	return fmt.Sprintf(mask, idList)
}

func currentUserURL() string {
	return BASE_URL + "/users/current"
}

func (qc *QuipClient) getWithToken(url string, token Token) *gorequest.SuperAgent {
	qc.logger.Debugf("Requesting %s with %s", url, token)
	return gorequest.New().Get(url).Set("Authorization", "Bearer "+token.value)
}
