package rep

import (
	"github.com/couchbase/gocb"
	"github.com/jeromedoucet/alienor-back/model"
)

func TeamExist(name string, consistency gocb.ConsistencyMode) (res bool, err error) {
	query := gocb.NewN1qlQuery("SELECT r.team.name FROM "+ bucketName +" a unnest a.roles r where a.type = $1 and r.team.name = $2")
	query.Consistency(consistency)
	var rows gocb.QueryResults
	rows, err = bucket.ExecuteN1qlQuery(query, []interface{}{model.USER, name})
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		return
	}
	var row interface{}
	res = rows.Next(&row)
	return
}
