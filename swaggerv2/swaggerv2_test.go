package swaggerv2

import (
	"bytes"
	"fmt"
	"html/template"
	"testing"
)

func TestParseSwaggerV2Template(t *testing.T) {
	tmpl := template.Must(template.New("swagger-ui").Parse(swaggerTemplateV2))
	buf := bytes.NewBuffer(nil)

	_ = tmpl.Execute(buf, map[string]any{
		"SwaggerHost":      "https://petstore.swagger.io",
		"SwaggerJsonsPath": []string{"swagger.json", "swagger2.json"},
	})
	fmt.Println(string(buf.Bytes()))
}
