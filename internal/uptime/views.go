package uptime

import (
	"bytes"
	"embed"
	"html/template"
)

//go:embed templates/dashboard.html static/style.css
var assets embed.FS

func dashboardHTML() ([]byte, error) {
	page, err := template.ParseFS(assets, "templates/dashboard.html")
	if err != nil {
		return nil, err
	}
	var body bytes.Buffer
	if err := page.Execute(&body, nil); err != nil {
		return nil, err
	}
	return body.Bytes(), nil
}
