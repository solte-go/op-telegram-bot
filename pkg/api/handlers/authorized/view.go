package authorized

import (
	"bytes"
	"html/template"
	"io"
	"net/http"
	"path"
	"path/filepath"
)

type View struct {
	Template *template.Template
	Layout   string
}

var (
	fileDir       = "templates/layouts"
	fileExtension = ".gohtml"
	templateDir   = "templates"
	staticContent = "static"
)

func setDirPath(dir string) {
	fileDir = filepath.Join(dir, fileDir)
	templateDir = path.Join(dir, templateDir)
	staticContent = path.Join(dir, staticContent)
}

func NewView(dir, layout string, files ...string) *View {
	setDirPath(dir)
	prependDir(files)

	files = append(files, filePath()...)
	t, err := template.ParseFiles(files...)
	if err != nil {
		panic(err)
	}
	return &View{
		Template: t,
		Layout:   layout,
	}

}

func prependDir(files []string) {
	for i, file := range files {
		files[i] = templateDir + "/" + file
	}
}

func filePath() (files []string) {
	files, err := filepath.Glob(fileDir + "/*" + fileExtension)
	if err != nil {
		panic(err)
	}
	return files
}

func (v *View) renderView(w http.ResponseWriter, r *http.Request, data interface{}) error {
	w.Header().Set("Content-Type", "text/html")
	var buffer bytes.Buffer
	if err := v.Template.ExecuteTemplate(&buffer, v.Layout, data); err != nil {
		return err
	}
	_, err := io.Copy(w, &buffer)
	if err != nil {
		return err
	}
	return nil
}
