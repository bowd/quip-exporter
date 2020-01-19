package client

import (
	"encoding/json"
)

func (qc *QuipClient) getFolders(ids []string) (map[string][]byte, error) {
	qc.logger.Debug("Waiting for token")
	token := qc.checkoutToken()
	defer qc.checkinToken(token)

	qc.logger.Infof("Querying folder batch [%d folders]", len(ids))
	return qc.getMap(foldersUrl(ids), token)
}

func (qc *QuipClient) getThreads(ids []string) (map[string][]byte, error) {
	qc.logger.Debug("Waiting for token")
	token := qc.checkoutToken()
	defer qc.checkinToken(token)

	qc.logger.Infof("Querying threads batch [%d threads]", len(ids))
	return qc.getMap(threadsUrl(ids), token)
}

func (qc *QuipClient) getMap(url, token string) (map[string][]byte, error) {
	resp, rawBody, errs := qc.getWithToken(url, token).End()
	if errs != nil {
		return nil, errs[0]
	}
	var body map[string]interface{}
	err := json.Unmarshal([]byte(rawBody), &body)
	if err != nil {
		qc.logger.Debug(resp)
		qc.logger.Error(rawBody)
		return nil, err
	}
	result := make(map[string][]byte)
	for key, value := range body {
		result[key], _ = json.Marshal(value)
	}
	return result, nil
}
