package utils

import (
	"net/http"
	"github.com/couchbase/gocb"
)

// simple mock for htp writer
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

type RepositoryHeader struct {
	DoGet func (identifier string, entity interface{}) (gocb.Cas, error)
	DoInsert func (entity interface{}) error
	DoUpdate func (entity interface{}, cas gocb.Cas) error
}

func (rep * RepositoryHeader) Get(identifier string, entity interface{}) (gocb.Cas, error) {
	return rep.DoGet(identifier, entity)
}

func (rep * RepositoryHeader) Insert(entity interface{}) error {
	return rep.DoInsert(entity)
}

func (rep * RepositoryHeader) Update(entity interface{}, cas gocb.Cas) error {
	return rep.DoUpdate(entity, cas)
}




