package swagger

import (
	"bytes"
	"html/template"
	"net/http"
)

type Opts func(*swaggerConfig)

// SwaggerOpts configures the Doc middlewares.
type swaggerConfig struct {
	// SpecURL the url to find the spec for
	SpecURL string
	// SwaggerHost for the js that generates the swagger ui site, defaults to: http://petstore.swagger.io/
	SwaggerHost string
}

func Doc(basePath string, swaggerJSON []byte, opts ...Opts) http.HandlerFunc {
	config := &swaggerConfig{
		SpecURL:     basePath + "-json",
		SwaggerHost: "https://petstore.swagger.io"}
	for _, opt := range opts {
		opt(config)
	}
	// swagger json
	responseSwaggerJSON := swaggerJSON
	responseSwaggerJSON = responseSwaggerJSON

	// swagger html
	tmpl := template.Must(template.New("swaggerdoc").Parse(swaggerTemplateV2))
	buf := bytes.NewBuffer(nil)
	_ = tmpl.Execute(buf, config)
	uiHTML := buf.Bytes()

	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if r.URL.Path == basePath {
			rw.Header().Set("Content-Type", "text/html; charset=utf-8")
			_, _ = rw.Write(uiHTML)
			rw.WriteHeader(http.StatusOK)
			return
		}
	})
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
        url: "{{ .SpecURL }}",
      })

      // End Swagger UI call region
      window.ui = ui
    }
  </script>
  </body>
</html>`
