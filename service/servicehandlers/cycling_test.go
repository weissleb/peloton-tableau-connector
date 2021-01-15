package servicehandlers

import (
	"reflect"
	"testing"
	"time"
)

type testType struct {
	s string
	i int32
	b bool
	f float32
	t time.Time
}

func Test_simpleSchemaForTable(t *testing.T) {
	type args struct {
		tableName string
		tableType interface{}
	}
	tests := []struct {
		name string
		args args
		want TableSchema
	}{
		{
			name: "test",
			args: args{
				tableName: "testType",
				tableType: testType{},
			},
			want: TableSchema{
				Name:        "testType",
				Description: "",
				Columns: []TableColumn{
					{Name: "s", GoType: "string"},
					{Name: "i", GoType: "int32"},
					{Name: "b", GoType: "bool"},
					{Name: "f", GoType: "float32"},
					{Name: "t", GoType: "time.Time"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := simpleSchemaForTable(tt.args.tableName, tt.args.tableType); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("simpleSchemaForTable() = %v, want %v", got, tt.want)
			}
		})
	}
}
