package utils

import "net/http"


type HttpWriterMock struct {
	Head http.Header
	Data []byte
	Code int
}

func (m *HttpWriterMock) Header() http.Header {
	return m.Head
}

func (m *HttpWriterMock) Write(data []byte) (nb int, err error) {
	m.Data = data
	return len(data), nil
}

func (m *HttpWriterMock) WriteHeader(code int) {
	m.Code = code
}


