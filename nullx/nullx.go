package nullx

import (
	"database/sql"

	"github.com/guregu/null/v5"
)

func NewInt(value sql.NullInt64) null.Int {
	return null.NewInt(value.Int64, value.Valid)
}

func NewString(value sql.NullString) null.String {
	return null.NewString(value.String, value.Valid)
}
