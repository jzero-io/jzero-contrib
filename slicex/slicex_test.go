package slicex

import (
	"testing"
)

func TestPaginate(t *testing.T) {
	type Data struct {
		Id   int
		Name string
	}

	datas := []Data{
		{Id: 1, Name: "a"},
		{Id: 2, Name: "b"},
		{Id: 3, Name: "c"},
		{Id: 4, Name: "d"},
		{Id: 5, Name: "e"},
		{Id: 6, Name: "f"},
		{Id: 7, Name: "g"},
		{Id: 8, Name: "h"},
		{Id: 9, Name: "i"},
		{Id: 10, Name: "j"},
	}

	paginate := Paginate[Data](datas, 1, 3)
	if len(paginate) != 3 {
		t.Errorf("Paginate error")
	}
	if paginate[0].Id != 1 {
		t.Errorf("Paginate error")
	}
	if paginate[1].Id != 2 {
		t.Errorf("Paginate error")
	}
	if paginate[2].Id != 3 {
		t.Errorf("Paginate error")
	}
}
