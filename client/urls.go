package client

import (
	"fmt"
	"strings"
)

type URLProvider interface {
	CurrentUser() string
	Folders(ids []string) string
	Threads(ids []string) string
	Users(ids []string) string
	ExportThread(threadID, exportType string) string
	Blob(threadID, blobID string) string
	ThreadComments(threadID string, cursor *uint64) string
}

type DefaultURLProvider struct {
	BaseURL string
}

func NewDefaultURLProvider(baseURL string) URLProvider {
	return DefaultURLProvider{BaseURL: baseURL}
}

func (provider DefaultURLProvider) CurrentUser() string {
	return provider.BaseURL + "/users/current"
}

func (provider DefaultURLProvider) Folders(ids []string) string {
	idList := strings.Join(ids, ",")
	return fmt.Sprintf(provider.BaseURL+"/folders?ids=%s", idList)
}

func (provider DefaultURLProvider) Threads(ids []string) string {
	idList := strings.Join(ids, ",")
	return fmt.Sprintf(provider.BaseURL+"/threads?ids=%s", idList)
}

func (provider DefaultURLProvider) Users(ids []string) string {
	idList := strings.Join(ids, ",")
	return fmt.Sprintf(provider.BaseURL+"/users?ids=%s", idList)
}

func (provider DefaultURLProvider) ExportThread(threadID, exportType string) string {
	return fmt.Sprintf(provider.BaseURL+"/threads/%s/export/%s", threadID, exportType)
}

func (provider DefaultURLProvider) Blob(threadID, blobID string) string {
	return fmt.Sprintf(provider.BaseURL+"/blobs/%s/%s", threadID, blobID)
}

func (provider DefaultURLProvider) ThreadComments(threadID string, cursor *uint64) string {
	if cursor == nil {
		return fmt.Sprintf(provider.BaseURL+"/messages/%s?count=100", threadID)
	} else {
		return fmt.Sprintf(provider.BaseURL+"/messages/%s?count=100&max_created_usec=%d", threadID, *cursor)
	}
}

type AdminURLProvider struct {
	*DefaultURLProvider
	CompanyID string
}

func NewAdminURLProvider(baseURL string, companyID string) URLProvider {
	return AdminURLProvider{
		DefaultURLProvider: &DefaultURLProvider{BaseURL: baseURL},
		CompanyID:          companyID,
	}
}

func (provider AdminURLProvider) Folders(ids []string) string {
	idList := strings.Join(ids, ",")
	return fmt.Sprintf(provider.BaseURL+"/admin/folders?ids=%s&company_id=%s", idList, provider.CompanyID)
}

func (provider AdminURLProvider) Threads(ids []string) string {
	idList := strings.Join(ids, ",")
	return fmt.Sprintf(provider.BaseURL+"/admin/threads?ids=%s&company_id=%s", idList, provider.CompanyID)
}

func (provider AdminURLProvider) Users(ids []string) string {
	idList := strings.Join(ids, ",")
	return fmt.Sprintf(provider.BaseURL+"/admin/users?ids=%s", idList)
}

func (provider AdminURLProvider) ExportThread(threadID, exportType string) string {
	return fmt.Sprintf(provider.BaseURL+"/admin/threads/%s/export/%s?company_id=%s", threadID, exportType, provider.CompanyID)
}

func (provider AdminURLProvider) Blob(threadID, blobID string) string {
	return fmt.Sprintf(provider.BaseURL+"/admin/blobs/%s/%s?company_id=%s", threadID, blobID, provider.CompanyID)
}

func (provider AdminURLProvider) ThreadComments(threadID string, cursor *uint64) string {
	if cursor == nil {
		return fmt.Sprintf(provider.BaseURL+"/admin/messages/%s?count=100&company_id=%s", threadID, provider.CompanyID)
	} else {
		return fmt.Sprintf(provider.BaseURL+"/admin/messages/%s?count=100&max_created_usec=%d&company_id=%s", threadID, *cursor, provider.CompanyID)
	}
}
