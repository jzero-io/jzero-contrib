package nullx

import (
	"fmt"
	"testing"

	"github.com/guregu/null/v5"
)

func TestNullx(t *testing.T) {
	var a *int64
	var va int64 = 5
	a = &va
	ptr := null.IntFromPtr(a)
	fmt.Println(ptr.NullInt64)
}
