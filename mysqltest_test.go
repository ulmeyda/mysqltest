package mysqltest

import (
	"reflect"
	"testing"
)

func TestMySQLTest_createInsertQuery(t *testing.T) {
	type testData struct {
		ID   int    `db:"id"`
		Name string `db:"name"`
	}

	type args struct {
		tbl   string
		value interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   []interface{}
		wantErr bool
	}{
		{
			name: "slice value",
			args: args{
				tbl: "test",
				value: []testData{
					{ID: 1, Name: "aaa"},
					{ID: 2, Name: "bbb"},
				},
			},
			want:    "INSERT INTO `test` (`id`,`name`) VALUES (?,?),(?,?)",
			want1:   []interface{}{1, "aaa", 2, "bbb"},
			wantErr: false,
		},
		{
			name: "struct value",
			args: args{
				tbl:   "test",
				value: testData{ID: 1, Name: "aaa"},
			},
			want:    "INSERT INTO `test` (`id`,`name`) VALUES (?,?)",
			want1:   []interface{}{1, "aaa"},
			wantErr: false,
		},
		{
			name: "zero slice value",
			args: args{
				tbl:   "test",
				value: []testData{},
			},
			want:    "",
			want1:   nil,
			wantErr: true,
		},
		{
			name: "string",
			args: args{
				tbl:   "test",
				value: []testData{},
			},
			want:    "",
			want1:   nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mt := &MySQLTest{
				db:     nil,
				tag:    "db",
				tables: nil,
			}
			got, got1, err := mt.createInsertQuery(tt.args.tbl, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("MySQLTest.createInsertQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MySQLTest.createInsertQuery() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("MySQLTest.createInsertQuery() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
