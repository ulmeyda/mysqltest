# mysqltest

test時ににmysqlにデータを挿入するヘルパー

## example

```go

type TestData struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

t.Run("test", func(t *testing.T) {

	testData1 := []TestData{
		{ID: 1, Name: "aaa"},
		{ID: 2, Name: "bbb"},
	}

	insertData := mysqltest.NewData().
		Set("table_name1", testData1).
		Set("table_name2", testData2)

	db, _ := sql.Open("mysql", "")
	defer mysqltest.New(db, "db").CleaningTables("table_name3").Exec(insertData)()
})

```
