package castx

import (
	"fmt"
	"testing"
)

func TestToSlice(t *testing.T) {
	slice := ToSlice("create_time DESC")
	fmt.Println(slice)

	slice = ToSlice(5)
	fmt.Println(slice)

	slice = ToSlice([]int{1, 2, 3, 4, 5})
	fmt.Println(slice)
}
