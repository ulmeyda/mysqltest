package mysqltest

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"reflect"

	"github.com/Masterminds/squirrel"
	orderedmap "github.com/wk8/go-ordered-map"
)

// Data テストデータ
type Data struct {
	datamap *orderedmap.OrderedMap
}

// NewData データコンストラクタ
func NewData() *Data {
	return &Data{
		datamap: orderedmap.New(),
	}
}

// Set セット
func (m *Data) Set(key, value interface{}) *Data {
	m.datamap.Set(key, value)
	return m
}

// MySQLTest テストデータ
type MySQLTest struct {
	db     *sql.DB
	tag    string
	tables []string
}

// New コンストラクタ
func New(db *sql.DB, tag string) *MySQLTest {
	return &MySQLTest{
		db:  db,
		tag: tag,
	}
}

// CleaningTables defer時に削除したいテーブルを追加する
func (t *MySQLTest) CleaningTables(tbl ...string) *MySQLTest {
	t.tables = append(t.tables, tbl...)
	return t
}

// Exec データセットし、deferで使うclean関数を呼び出す
func (t *MySQLTest) Exec(d *Data) func() {
	for pair := d.datamap.Oldest(); pair != nil; pair = pair.Next() {
		tbl := pair.Key.(string)
		t.insert(tbl, pair.Value)
		t.tables = append(t.tables, tbl)
	}
	return func() {
		t.clean(t.tables...)
		t.tables = []string{}
	}
}

// clean Tableデータのclear
func (t *MySQLTest) clean(tbl ...string) {
	_, err := t.db.Exec("SET FOREIGN_KEY_CHECKS = 0")
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range tbl {
		_, err := t.db.Exec(fmt.Sprintf("TRUNCATE TABLE %s", v))
		if err != nil {
			log.Fatal(err)
		}
	}
	_, err = t.db.Exec("SET FOREIGN_KEY_CHECKS = 1")
	if err != nil {
		log.Fatal(err)
	}
}

// insert データ挿入
func (t *MySQLTest) insert(tbl string, value interface{}) {
	query, args, err := t.createInsertQuery(tbl, value)
	if err != nil {
		log.Fatal(err)
	}
	_, err = t.db.Exec(query, args...)
	if err != nil {
		log.Fatalf("err=%+v, query=%s, ", err, query)
	}
}

// createInsertQuery insertクエリの作成
func (t *MySQLTest) createInsertQuery(tbl string, value interface{}) (string, []interface{}, error) {
	var ifs []interface{}

	v := reflect.Indirect(reflect.ValueOf(value))

	switch v.Type().Kind() {
	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			ifs = append(ifs, v.Index(i).Interface())
		}
	case reflect.Struct:
		ifs = []interface{}{value}
	default:
		return "", nil, errors.New("value should be struct or slice.")
	}

	if len(ifs) == 0 {
		return "", nil, errors.New("value is empty")
	}

	quote := func(s string) string {
		return fmt.Sprintf("`%s`", s)
	}

	rv := reflect.ValueOf(ifs[0])
	if rv.Kind() == reflect.Ptr {
		rv = reflect.ValueOf(ifs[0]).Elem()
	}
	rt := rv.Type()
	columns := make([]string, 0, rt.NumField())
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		columns = append(columns, quote(f.Tag.Get(t.tag)))
	}
	builder := squirrel.Insert(quote(tbl)).Columns(columns...)

	values := make([]string, 0, len(columns))
	for _, column := range columns {
		values = append(values, fmt.Sprintf("%s = VALUES(%s)", column, column))
	}

	for _, vv := range ifs {
		rv := reflect.ValueOf(vv)
		if rv.Kind() == reflect.Ptr {
			rv = reflect.ValueOf(vv).Elem()
		}
		rt := rv.Type()
		var values []interface{}
		for i := 0; i < rt.NumField(); i++ {
			f := rt.Field(i)
			column := rv.FieldByName(f.Name).Interface()
			values = append(values, column)
		}
		builder = builder.Values(values...)
	}
	return builder.ToSql()
}
