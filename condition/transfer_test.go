package condition

import (
	"testing"

	"github.com/huandu/go-sqlbuilder"
)

func TestSelectWhere(t *testing.T) {
	type args struct {
		sb              *sqlbuilder.SelectBuilder
		originalField   string
		paramJoinString string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test1",
			args: args{
				sb:              sqlbuilder.NewSelectBuilder(),
				originalField:   "`sys_user_id` = ? and `sys_authority_authority_id` = ?",
				paramJoinString: "sysUserId,sysAuthorityAuthorityId",
			},
		},

		{
			name: "test2",
			args: args{
				sb:              sqlbuilder.NewSelectBuilder(),
				originalField:   "`sys_user_id` = ?",
				paramJoinString: "sysUserId",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SelectByWhereRawSql(tt.args.sb, tt.args.originalField, tt.args.paramJoinString)

			sql, arguments := tt.args.sb.Build()
			t.Log(sql, arguments)
		})
	}
}
