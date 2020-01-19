package client

import "time"

func (qc *QuipClient) checkoutToken() string {
	qc.tokenIndexMutex.Lock()
	qc.lastTokenIndex = (qc.lastTokenIndex + 1) % len(qc.tokens)
	token := qc.tokens[qc.lastTokenIndex]
	qc.tokenIndexMutex.Unlock()
	qc.tokenMutex[token].Lock()
	return token
}

func (qc *QuipClient) checkinToken(token string) {
	<-time.After(time.Duration(int(1000/qc.rps) * int(time.Millisecond)))
	qc.tokenMutex[token].Unlock()
}
