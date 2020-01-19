package client

import (
	"fmt"
	"github.com/parnurzeal/gorequest"
	"strings"
)

func (qc *QuipClient) testToken(token string) error {
	url := currentUserUrl()
	resp, rawBody, _ := qc.getWithToken(url, token).End()
	if resp.StatusCode > 400 {
		return fmt.Errorf("%s", rawBody)
	}
	return nil
}

const (
	BASE_URL     = "https://platform.quip.com/1"
	FOLDERS_MASK = "/folders/?ids=%s"
	THREADS_MASK = "/threads/?ids=%s"
	FOLDER_MASK  = "/folders/%s"
)

func foldersUrl(folders []string) string {
	mask := BASE_URL + FOLDERS_MASK
	idList := strings.Join(folders, ",")
	return fmt.Sprintf(mask, idList)
}

func threadsUrl(threads []string) string {
	mask := BASE_URL + THREADS_MASK
	idList := strings.Join(threads, ",")
	return fmt.Sprintf(mask, idList)
}

func currentUserUrl() string {
	return BASE_URL + "/users/current"
}

func (qc *QuipClient) getWithToken(url, token string) *gorequest.SuperAgent {
	qc.logger.Debugf("Requesting %s with %s", url, token)
	return gorequest.New().Get(url).Set("Authorization", "Bearer "+token)
}
