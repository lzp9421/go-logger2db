package model

var table string = "log"

type Log struct {
	model
}

func (this *Log)Insert(data map[string]interface{}) int64 {
	return this.model.Insert(table, data)
}