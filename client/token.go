package client

import "time"

type Token struct {
	value string
	index int
}

func (qc *QuipClient) checkoutToken() Token {
	qc.tokenIndexMutex.Lock()
	index := (qc.lastTokenIndex + 1) % qc.tokenConcurrency
	qc.lastTokenIndex = index
	qc.tokenIndexMutex.Unlock()
	qc.tokenMutex[index].Lock()
	return Token{qc.token, index}
}

func (qc *QuipClient) checkinToken(token Token) {
	go func() {
		<-time.After(time.Duration(int(1000/qc.rps) * int(time.Millisecond)))
		qc.tokenMutex[token.index].Unlock()
	}()
}
