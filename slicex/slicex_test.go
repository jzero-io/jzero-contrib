package slicex

import (
	"testing"
)

type Data struct {
	Id   int
	Name string
}

var datas = []Data{
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

func TestPaginate(t *testing.T) {
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

func TestToMap(t *testing.T) {
	toMap := ToMap(datas, func(row Data) int {
		return row.Id
	})
	if len(toMap) != 10 {
		t.Errorf("ToMap error")
	}
	if toMap[1].Id != 1 {
		t.Errorf("ToMap error")
	}
}
