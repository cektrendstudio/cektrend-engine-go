package postgres

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/cektrendstudio/cektrend-engine-go/models"
	"github.com/cektrendstudio/cektrend-engine-go/pkg/serror"
	"github.com/jmoiron/sqlx"
)

func isZero(val reflect.Value) bool {
	if val.Kind() == reflect.Slice {
		return val.Len() == 0
	}
	return reflect.DeepEqual(val.Interface(), reflect.Zero(val.Type()).Interface())
}

func buildQueryDynamicFilter(selectQuery string, filterStruct interface{}, orderBy string, page string, limit string) (queryCount string, argsCount []interface{}, query string, args []interface{}, errx serror.SError) {
	baseErrLog := `[postgres][buildQueryDynamicFilter]`

	var filterQuery, countQuery string
	var execValues []interface{}

	x := reflect.ValueOf(filterStruct)
	for i := 0; i < x.NumField(); i++ {
		colName := x.Type().Field(i).Tag.Get("filter")
		distinct := x.Type().Field(i).Tag.Get("distinct")
		valCheck := x.Field(i)
		val := x.Field(i).Interface()
		if colName != "-" && !isZero(valCheck) {
			operator := x.Type().Field(i).Tag.Get("db_op")
			switch operator {
			case "ilike":
				execValues = append(execValues, "%"+x.Field(i).String()+"%")
				filterQuery += ` AND ` + colName + ` ILIKE ?`
			case "min":
				execValues = append(execValues, val)
				filterQuery += ` AND ` + colName + ` >= ?`
			case "max":
				execValues = append(execValues, val)
				filterQuery += ` AND ` + colName + ` <= ?`
			case "not":
				execValues = append(execValues, val)
				filterQuery += ` AND ` + colName + ` != ?`
			case "not_like":
				execValues = append(execValues, "%"+x.Field(i).String()+"%")
				filterQuery += ` AND ` + colName + ` NOT ILIKE ?`
			case "in":
				values := reflect.ValueOf(val)
				placeholders := make([]string, values.Len())
				for j := 0; j < values.Len(); j++ {
					execValues = append(execValues, values.Index(j).Interface())
					placeholders[j] = "?"
				}
				filterQuery += ` AND ` + colName + ` IN (` + strings.Join(placeholders, ", ") + `)`
			default:
				execValues = append(execValues, val)
				filterQuery += ` AND ` + colName + ` = ?`
			}
		}
		if distinct != "" && val != reflect.Zero(x.Type().Field(i).Type).Interface() {
			orderBy = distinct + ` DESC, ` + orderBy
			countQuery = fmt.Sprintf(`SELECT COUNT(distinct %s) FROM `, distinct)
		}
	}

	var paginationQuery string
	if page != "" && limit != "" {
		qc := strings.Split(selectQuery, "FROM")

		if countQuery == "" {
			countQuery = `SELECT count(*) FROM `
		}

		qCount, vCount, err := sqlx.In(countQuery+qc[1]+filterQuery, execValues...)
		if err != nil {
			return query, args, queryCount, argsCount, serror.NewFromErrorc(err, baseErrLog+" failed to build query count")
		}
		queryCount = qCount
		argsCount = vCount

		p, err := strconv.Atoi(page)
		if err != nil {
			return query, args, queryCount, argsCount, serror.NewFromErrorc(err, baseErrLog+" failed to convert page to int")
		}

		l, err := strconv.Atoi(limit)
		if err != nil {
			return query, args, queryCount, argsCount, serror.NewFromErrorc(err, baseErrLog+" failed to convert limit to int")
		}

		offset := 0
		if p != 1 {
			offset = p * l
		}
		paginationQuery = ` OFFSET ? LIMIT ?`
		execValues = append(execValues, offset, l)
	}

	queryOrderBy := ` ORDER BY ` + orderBy

	query, args, err := sqlx.In(selectQuery+filterQuery+queryOrderBy+paginationQuery, execValues...)
	if err != nil {
		return query, args, queryCount, argsCount, serror.NewFromErrorc(err, baseErrLog+" failed to build main query")
	}

	return queryCount, argsCount, query, args, nil
}

func buildQueryDynamicUpdate(tableName string, formStruct interface{}, keyCol string, keyVal interface{}, operatedBy string) (query string, args []interface{}, errx serror.SError) {

	var querySets []string
	var execValues []interface{}
	x := reflect.ValueOf(formStruct)
	for i := 0; i < x.NumField(); i++ {
		colName := x.Type().Field(i).Tag.Get("db")
		val := x.Field(i).Interface()
		if colName != "-" && val != reflect.Zero(x.Type().Field(i).Type).Interface() && colName != keyCol {
			execValues = append(execValues, val)
			querySets = append(querySets, colName+` = ?`)
		}
	}

	if operatedBy == "" {
		operatedBy = models.CreatedBySystem
	}
	modifedByQuery := ` updated_by = ` + fmt.Sprintf(`'%s'`, operatedBy)

	querySets = append(querySets, modifedByQuery, ` updated_at = now()`)
	execValues = append(execValues, keyVal)
	queryWhere := ` WHERE ` + keyCol + ` = ?`

	querySet := strings.Join(querySets, `, `)

	query, args, err := sqlx.In(`UPDATE `+tableName+` SET `+querySet+queryWhere, execValues...)
	if err != nil {
		errx = serror.NewFromErrorc(err, `[postgres][buildQueryDynamicUpdate][`+tableName+`] failed to build main query`)
		return
	}

	return
}
