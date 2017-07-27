package main

import (
	"net/http"
	"io"
	"log"
	"os"
	"io/ioutil"
	"html/template"
	"path"
	"strings"
	"runtime/debug"
	"fmt"
)

const (
	UPLOAD_DIR   = "src/photoweb/uploads"
	TEMPLATE_DIR = "src/photoweb/views"
)

var templates = make(map[string]*template.Template)

func init() {
	fmt.Println(os.Getwd())
	fileInfoArr, err := ioutil.ReadDir(TEMPLATE_DIR)
	if err != nil {
		panic(err)
		return
	}
	var templateName, templatePath string
	for _, fileInfo := range fileInfoArr {
		templateName = fileInfo.Name()
		if ext := path.Ext(templateName); ext != ".html" {
			continue
		}
		templatePath = TEMPLATE_DIR + "/" + templateName
		log.Println("Loading template:", templatePath)
		t := template.Must(template.ParseFiles(templatePath))
		tmpl := strings.Split(templateName, ".")[0]
		templates[tmpl] = t
	}
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		err := renderHtml(w, "upload", nil)
		check(err)
	}
	if r.Method == "POST" {
		f, h, err := r.FormFile("image")
		check(err)
		filename := h.Filename
		defer f.Close()
		t, err := os.Create(UPLOAD_DIR + "/" + filename)
		check(err)
		defer t.Close()
		_, err = io.Copy(t, f)
		check(err)
		http.Redirect(w, r, "/view?id="+filename, http.StatusFound)
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	imageId := r.FormValue("id")
	imagePath := UPLOAD_DIR + "/" + imageId
	if exists := isExists(imagePath); !exists {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "image")
	http.ServeFile(w, r, imagePath)
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	fileInfoArr, err := ioutil.ReadDir(UPLOAD_DIR)
	check(err)

	locals := make(map[string]interface{})
	images := []string{}
	for _, fileInfo := range fileInfoArr {
		images = append(images, fileInfo.Name())
	}
	locals["images"] = images
	err = renderHtml(w, "list", locals)
	check(err)
}

func renderHtml(w http.ResponseWriter, tmpl string, locals map[string]interface{}) (err error) {
	err = templates[tmpl].Execute(w, locals)
	return
}

func isExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return os.IsExist(err)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func safeHandler(fn http.HandlerFunc) http.HandlerFunc  {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if e, ok := recover().(error); ok {
				http.Error(w, e.Error(), http.StatusInternalServerError)
				// 或者输出自定义的 50x 错误页面
				// w.WriteHeader(http.StatusInternalServerError)
				// renderHtml(w, "error", e)
				// logging
				log.Println("WARN: panic in %v - %v", fn, e)
				log.Println(string(debug.Stack()))
			}
		}()
		fn(w, r)
	}

}

func main() {
	http.HandleFunc("/", safeHandler(listHandler))
	http.HandleFunc("/view", safeHandler(viewHandler))
	http.HandleFunc("/upload", safeHandler(uploadHandler))
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err.Error())
	}
}
