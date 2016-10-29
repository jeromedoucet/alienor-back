package utils

import (
	"net/http"
	"github.com/couchbase/gocb"
	"github.com/jeromedoucet/alienor-back/model"
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
	DoGet func (identifier string, document model.Document) (gocb.Cas, error)
	DoInsert func (document model.Document) error
	DoUpdate func (document model.Document, cas gocb.Cas) error
}

func (rep * RepositoryHeader) Get(identifier string, document model.Document) (gocb.Cas, error) {
	return rep.DoGet(identifier, document)
}

func (rep * RepositoryHeader) Insert(document model.Document) error {
	return rep.DoInsert(document)
}

func (rep * RepositoryHeader) Update(document model.Document, cas gocb.Cas) error {
	return rep.DoUpdate(document, cas)
}

type MockDocument struct {
	Id string
}

func (d * MockDocument) Identifier() string {
	return d.Id
}




