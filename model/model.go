package model

import (
	"database/sql"
	"fmt"
	"strings"
	"strconv"
	"math"
	"logger/library"
)

type model struct {
}

// 分页返回数据 - 表字段名定义
type ResultRow struct {

}

// 分页返回数据 - 返回结果定义
type ResultData struct {
	List    []interface{}
	Count   int64
	PerPage int
	Page    int
	PageCount int
	Start   int
	Mark    int
}

func (c *model)Mysql() *sql.DB {
	return library.Db
}

// 构造分页返回结果
func (mdl *model) Multi(count int64, perpage int, page int) ResultData {
	var data ResultData
	data.Count = int64(math.Max(float64(count), 0))
	data.Page = page
	data.PerPage = int(math.Max(float64(perpage), 1))
	data.PageCount = int(math.Ceil(float64(data.Count) / float64(data.PerPage)))
	data.Start = int(math.Max(float64(data.PerPage * data.Page - data.PerPage), 0))
	data.Mark = data.Start + 1
	return data
}

func (mdl *model) buildParams(data []interface{}, split string) string {
	var cond []string
	for _, d := range data {
		dd := d.(map[string]interface{})
		v := library.ToString(dd["val"])
		cond = append(cond, fmt.Sprintf("`%s` %s '%s'", dd["key"], dd["oper"], library.Addslashes(v)))
	}
	return strings.Join(cond, split)
}

func (mdl *model) BuildCond(data map[string]interface{}) string {
	var cond []interface{}
	var condStr []string
	for k, d := range data {
		kv := make(map[string]interface{})
		if k == "_in_" {
			condStr = append(condStr, library.ToString(d))
		}else{
			switch d.(type) {
			case map[string]interface{}:
				dd := d.(map[string]interface{})
				if dd["oper"] == nil {
					kv["oper"] = "="
				} else {
					kv["oper"] = dd["oper"]
				}
				kv["key"] = k
				kv["val"] = dd["val"]
			default:
				kv["key"] = k
				kv["oper"] = "="
				kv["val"] = d
			}
			cond = append(cond, kv)
		}
	}
	if len(data) == 0 {
		return ""
	}
	result := " WHERE "
	addAnd := false
	if len(condStr) != 0 {
		result += strings.Join(condStr, "and")
		addAnd = true
	}
	if len(cond) != 0 {
		if addAnd {
			result += " and "
		}
		result += mdl.buildParams(cond, " and ")
	}
	return result
}

func (mdl *model) BuildVal(data map[string]interface{}) string {
	var cond []interface{}
	for k, d := range data {
		kv := map[string]interface{} {
			"key": k,
			"oper": "=",
			"val": d,
		}
		cond = append(cond, kv)
	}
	return mdl.buildParams(cond, ",")
}

// 查询统计个数
func (mdl *model) Count(table string, cond map[string]interface{}) int64 {
	parsedCond := mdl.BuildCond(cond)

	if table == "" {
		panic("query error: no table")
	}
	queryString := fmt.Sprintf("SELECT count(*) as nums FROM %s ", table)
	if parsedCond != "" {
		queryString += parsedCond
	}
	data := mdl.GetOne(queryString)
	count,err := strconv.ParseInt(data["nums"], 10, 64)
	if err != nil {
		count = 0
	}
	return count
}

// 拼接查询语句字符串
func (mdl *model) MakeQueryString(table string, fields string, cond string, group string, order string, limit int, start int) string {

	queryString := fmt.Sprintf("SELECT %s FROM %s ", fields, table)

	if cond != "" {
		queryString += " " + cond
	}

	if group != "" {
		queryString += " GROUP BY " + group
	}

	if order != "" {
		queryString += " ORDER BY " + order
	}

	if limit > 0 && start > 0 {
		queryString += " LIMIT " + strconv.Itoa(start) + "," + strconv.Itoa(limit)
	}else if limit > 0 {
		queryString += " LIMIT " + strconv.Itoa(limit)
	}

	return queryString
}

// 查询并返回多条记录，且包含分页信息
func (mdl *model) FetchWithPage(table string, fields string, cond map[string]interface{}, order string, limit int, page int, priKey string) *ResultData {
	count := mdl.Count(table, cond)
	data := mdl.Multi(count, limit, page)

	start := int(math.Max(float64(page * limit - limit), 0))

	data.List = mdl.Fetch(table, fields, cond, order, limit, start, priKey)
	return &data
}

// 解析结果集
func (mdl *model) ParseData(rows *sql.Rows) []interface{} {
	data := []interface{}{}

	columns,err := rows.Columns()
	library.CheckError(err)
	fCount := len(columns)
	fieldPtr := make([]interface{}, fCount)
	fieldArr := make([]sql.RawBytes, fCount)
	fieldToID := make(map[string]int64, fCount)
	for k,v := range columns {
		fieldPtr[k] = &fieldArr[k]
		fieldToID[v] = int64(k)
	}

	for rows.Next() {
		err = rows.Scan(fieldPtr...)
		library.CheckError(err)

		m := make(map[string]string, fCount)

		for k, v := range fieldToID {
			if fieldArr[v] == nil {
				m[k] = ""
			} else {
				m[k] = string(fieldArr[v])
			}
		}
		data = append(data, m)
	}

	err = rows.Err()
	library.CheckError(err)
	return data
}

// 查询并返回多条记录
func (mdl *model) Fetch(table string, fields string, cond map[string]interface{}, order string, limit int, start int, priKey string) []interface{} {

	parsedCond := mdl.BuildCond(cond)

	queryString := mdl.MakeQueryString(table, fields, parsedCond, "", order, limit, start)

	rows, err := mdl.Mysql().Query(queryString)
	library.CheckQueryError(err, queryString)
	defer rows.Close()

	lists := mdl.ParseData(rows)

	return lists
}

func (mdl *model) GetAll(query string) []interface{} {
	rows, err := mdl.Mysql().Query(query)
	library.CheckQueryError(err, query)
	defer rows.Close()
	outArr := mdl.ParseData(rows)

	return outArr
}

func (mdl *model) GetOne(query string) map[string]string {
	rows, err := mdl.Mysql().Query(query)
	library.CheckQueryError(err, query)
	defer rows.Close()

	columns,err := rows.Columns()
	library.CheckError(err)
	fCount := len(columns)
	fieldPtr := make([]interface{}, fCount)
	fieldArr := make([]sql.RawBytes, fCount)
	fieldToID := make(map[string]int64, fCount)
	for k,v := range columns {
		fieldPtr[k] = &fieldArr[k]
		fieldToID[v] = int64(k)
	}

	m := make(map[string]string, fCount)
	if rows.Next() {
		err = rows.Scan(fieldPtr...)
		library.CheckError(err)

		for k, v := range fieldToID {
			if fieldArr[v] == nil {
				m[k] = ""
			} else {
				m[k] = string(fieldArr[v])
			}
		}
	}
	err = rows.Err()
	library.CheckError(err)
	return m
}


func (mdl *model) Insert(table string, set map[string]interface{}) int64 {
	stmt, err := mdl.Mysql().Prepare("Insert into " +table+ " set " + mdl.BuildVal(set))
	library.CheckError(err)
	defer stmt.Close()

	res, err := stmt.Exec()
	library.CheckError(err)

	id, err := res.LastInsertId()
	library.CheckError(err)
	return id
}

func (mdl *model) Delete(table string, cond map[string]interface{}) int64 {
	stmt, err := mdl.Mysql().Prepare("delete from " +table+ " " + mdl.BuildCond(cond))
	library.CheckError(err)
	defer stmt.Close()

	res, err := stmt.Exec()
	library.CheckError(err)

	id, err := res.RowsAffected()
	library.CheckError(err)
	return id
}

func (mdl *model) Update(table string, set map[string]interface{}, cond map[string]interface{}) int64 {
	where := mdl.BuildCond(cond)
	if len(where) <= 0 {
		return 0
	}

	stmt, err := mdl.Mysql().Prepare("update " +table+ " set " + mdl.BuildVal(set) + where)
	library.CheckError(err)
	defer stmt.Close()

	res, err := stmt.Exec()
	library.CheckError(err)

	id, err := res.RowsAffected()
	library.CheckError(err)
	return id
}