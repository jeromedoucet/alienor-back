package rep

import (
	"github.com/jeromedoucet/alienor-back/model"
	"errors"
	"github.com/couchbase/gocb"
)

type ItemRepository struct {

}

func (ItemRepository) Get(identifier string, document model.Document) (gocb.Cas, error) {
	return 0, nil
}

func (ItemRepository) Insert(document model.Document) (err error) {
	item, isItem := document.(*model.Item)
	if !isItem {
		err = errors.New("Cannot Insert a non user entity !")
		return
	}
	// todo check mandatory field
	_, err = bucket.Insert(string(model.ITEM) + ":" + item.TeamId + ":" + item.Id, item, 0)
	return
}

func (ItemRepository) Update(document model.Document, cas gocb.Cas) error {
	return nil
}
