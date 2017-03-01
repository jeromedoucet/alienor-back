package rep

import (
	"github.com/jeromedoucet/alienor-back/model"
	"errors"
	"github.com/couchbase/gocb"
)

type ItemRepository struct {
}

func (ItemRepository) Get(parentId, identifier string, document model.Document) (gocb.Cas, error) {
	item, isItem := document.(*model.Item)
	if !isItem {
		return 0, errors.New("Cannot Get a non item entity")
	}

	return bucket.Get(string(model.ITEM)+":"+parentId+":"+identifier, item)
}

func (ItemRepository) Insert(parentId string, document model.Document) (err error) {
	item, isItem := document.(*model.Item)
	if !isItem {
		err = errors.New("can only insert item")
		return
	} else if item.Type == "" {
		err = errors.New("missing type in item")
		return
	} else if item.Id == "" {
		err = errors.New("missing id in item")
		return
	} else if item.State != model.Newly { //todo is it really a desired behavior ?
		err = errors.New("bad item status")
		return
	}
	_, err = bucket.Insert(string(model.ITEM)+":"+parentId+":"+item.Id, item, 0)
	return
}

func (ItemRepository) Update(parentId string, document model.Document, cas gocb.Cas) error {
	return nil
}
