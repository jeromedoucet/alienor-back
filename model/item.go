package model

const Newly ItemState = "New"
const Incomplete = "Incomplete"
const Complete ItemState = "Complete"
const Error ItemState = "Error"

type ItemState string

type Item struct {
	Id    string `json:"id"`
	Type     DocType `json:"type"`
	State ItemState `json:"state"`
	TeamId string `json:"team_id"`
	Values map [string]string `json:"values"`
}

func NewItem() *Item {
	i := Item{State:Newly, Type:ITEM}
	return &i
}

func (i *Item) Identifier() string {
	return i.Id
}

