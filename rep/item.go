package rep

import (
	"errors"
	"github.com/couchbase/gocb"
	"github.com/jeromedoucet/alienor-back/model"
)

type ItemRepository struct {
}

func (ItemRepository) Get(parentId, identifier string, doc model.Document) (cas gocb.Cas, err error) {
	item, isItem := doc.(*model.Item)
	if !isItem {
		return 0, errors.New("Cannot Get a non item entity")
	}
	cas, err = bucket.Get(string(model.ITEM)+":"+parentId+":"+identifier, item)
	item.Version = uint64(cas)
	return
}

func (ItemRepository) Insert(parentId string, doc model.Document) (err error) {
	item, isItem := doc.(*model.Item)
	if !isItem {
		err = errors.New("can only insert item")
		return
	} else if item.Type == "" {
		err = errors.New("missing type in item")
		return
	} else if item.Id == "" {
		err = errors.New("missing id in item")
		return
	} else if item.State != model.Newly {
		err = errors.New("bad item status")
		return
	}
	var cas gocb.Cas
	cas, err = bucket.Insert(string(model.ITEM)+":"+parentId+":"+item.Id, item, 0)
	item.Version = uint64(cas)
	return
}

func (ItemRepository) Update(parentId string, doc model.Document, cas gocb.Cas) error {
	return nil
}

func (ItemRepository) Delete(parentId, id string, doc model.Document) error {
	item, isItem := doc.(*model.Item)
	if !isItem {
		return errors.New("Cannot delete a non item entity")
	}
	_, err := bucket.Remove(string(model.ITEM)+":"+parentId+":"+item.Id, gocb.Cas(item.Version))
	return err
}
