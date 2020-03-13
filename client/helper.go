package client

import (
	"fmt"
	"github.com/parnurzeal/gorequest"
)

func (qc *QuipClient) testToken() error {
	url := qc.url.CurrentUser()
	resp, rawBody, _ := qc.getWithToken(url, Token{qc.token, 0}).End()
	if resp.StatusCode > 400 {
		return fmt.Errorf("%s", rawBody)
	}
	return nil
}

func (qc *QuipClient) getWithToken(url string, token Token) *gorequest.SuperAgent {
	qc.logger.WithField("url", url).WithField("token", token.index).Debugf("executing query")
	return gorequest.New().Get(url).Set("Authorization", "Bearer "+token.value)
}
