package swaggerv2

import (
	"github.com/zeromicro/go-zero/rest"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

type Opts func(*swaggerConfig)

// SwaggerOpts configures the Doc middlewares.
type swaggerConfig struct {
	// SwaggerPath the path to find the spec for
	SwaggerPath string
	// SwaggerHost for the js that generates the swagger ui site, defaults to: http://petstore.swagger.io/
	SwaggerHost string
}

func RegisterRoutes(server *rest.Server, opts ...Opts) {
	config := &swaggerConfig{
		SwaggerPath: filepath.Join("app", "desc", "swagger"),
		SwaggerHost: "https://petstore.swagger.io"}
	for _, opt := range opts {
		opt(config)
	}

	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/swagger/:path",
		Handler: rawHandler(config),
	})

	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/swagger",
		Handler: uiHandler(config),
	})
}

func rawHandler(config *swaggerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		file, err := os.ReadFile(filepath.Join(config.SwaggerPath, path.Base(r.URL.Path)))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, _ = w.Write(file)
	}
}

func uiHandler(config *swaggerConfig) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "text/html; charset=utf-8")

		dir, err := os.ReadDir(config.SwaggerPath)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		var swaggerJsonsPath []string
		for _, fi := range dir {
			if fi.IsDir() {
				continue
			}
			if filepath.Ext(fi.Name()) == ".json" {
				swaggerJsonsPath = append(swaggerJsonsPath, filepath.Join("swagger", fi.Name()))
			}
		}

		uiHTML, _ := ParseTemplate(map[string]interface{}{
			"SwaggerHost":      config.SwaggerHost,
			"SwaggerJsonsPath": swaggerJsonsPath,
		}, []byte(swaggerTemplateV2))

		_, _ = rw.Write(uiHTML)
		rw.WriteHeader(http.StatusOK)
		return
	}
}

const swaggerTemplateV2 = `
	<!-- HTML for static distribution bundle build -->
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8">
    <title>API documentation</title>
    <link rel="stylesheet" type="text/css" href="{{ .SwaggerHost }}/swagger-ui.css" >
    <link rel="icon" type="image/png" href="{{ .SwaggerHost }}/favicon-32x32.png" sizes="32x32" />
    <link rel="icon" type="image/png" href="{{ .SwaggerHost }}/favicon-16x16.png" sizes="16x16" />
    <style>
      html
      {
        box-sizing: border-box;
        overflow: -moz-scrollbars-vertical;
        overflow-y: scroll;
      }

      *,
      *:before,
      *:after
      {
        box-sizing: inherit;
      }

      body
      {
        margin:0;
        background: #fafafa;
      }
    </style>
  </head>

  <body>
    <div id="swagger-ui"></div>

    <script src="{{ .SwaggerHost }}/swagger-ui-bundle.js"> </script>
    <script src="{{ .SwaggerHost }}/swagger-ui-standalone-preset.js"> </script>
    <script>
    window.onload = function() {
      // Begin Swagger UI call region
      const ui = SwaggerUIBundle({
        "dom_id": "#swagger-ui",
        deepLinking: true,
        presets: [
          SwaggerUIBundle.presets.apis,
          SwaggerUIStandalonePreset
        ],
        plugins: [
          SwaggerUIBundle.plugins.DownloadUrl
        ],
        layout: "StandaloneLayout",
		validatorUrl: null,
        urls: [
			{{range $k, $v := .SwaggerJsonsPath}}{url: "{{ $v }}", name: "{{ $v | base }}"},
			{{end}}
		]
      })

      // End Swagger UI call region
      window.ui = ui
    }
  </script>
  </body>
</html>`
