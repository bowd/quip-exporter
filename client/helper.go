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
	BaseUrl                  = "https://platform.quip.com/1"
	FoldersMask              = "/admin/folders/?ids=%s&company_id=%s"
	ThreadsMask              = "/admin/threads/?ids=%s&company_id=%s"
	UsersMask                = "/admin/users/?ids=%s&company_id=%s"
	FolderMask               = "/admin/folders/%s?company_id=%s"
	ExportMask               = "/admin/threads/%s/export/%s?company_id=%s"
	CurrentUserPath          = "/users/current"
	ThreadCommentsMask       = "/admin/messages/%s?count=100&company_id=%s"
	ThreadCommentsCursorMask = "/admin/messages/%s?count=100&max_created_usec=%d&company_id=%s"
	BlobMask                 = "/admin/blob/%s/%s?company_id=%s"
)

func batchURL(pathMask, companyID string, ids []string) string {
	mask := BaseUrl + pathMask
	idList := strings.Join(ids, ",")
	return fmt.Sprintf(mask, idList, companyID)
}

func currentUserURL() string {
	return BaseUrl + CurrentUserPath
}

func exportThreadURL(threadID, exportType, companyID string) string {
	return BaseUrl + fmt.Sprintf(ExportMask, threadID, exportType, companyID)
}

func blobURL(threadID, blobID, companyID string) string {
	return BaseUrl + fmt.Sprintf(BlobMask, threadID, blobID, companyID)
}

func threadCommentsURL(threadID string, cursor *uint64, companyID string) string {
	if cursor == nil {
		return BaseUrl + fmt.Sprintf(ThreadCommentsMask, threadID, companyID)
	} else {
		return BaseUrl + fmt.Sprintf(ThreadCommentsCursorMask, threadID, *cursor, companyID)
	}
}

func (qc *QuipClient) getWithToken(url string, token Token) *gorequest.SuperAgent {
	qc.logger.WithField("url", url).WithField("token", token.index).Debugf("executing query")
	return gorequest.New().Get(url).Set("Authorization", "Bearer "+token.value)
}
