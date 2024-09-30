package swaggerv2

import (
	"bytes"
	"html/template"

	sprig "github.com/Masterminds/sprig/v3"
)

// ParseTemplate template
func ParseTemplate(data any, tplT []byte) ([]byte, error) {
	t := template.Must(template.New("production").Funcs(sprig.HtmlFuncMap()).Parse(string(tplT)))

	buf := new(bytes.Buffer)
	err := t.Execute(buf, data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), err
}
