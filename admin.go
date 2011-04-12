package main

import (
	"http"
	"compress/gzip"
	"template"
	"fmt"
	"bytes"
	"io/ioutil"
	"strconv"
	"path"
	"strings"
	"mime"
	"sort"
)

const errHtml = `
	<html>
		<head>
		<style>
		body{
			font-family: Arial;
		}
		div{
			padding: 15px 30px;
			border: 1px solid black;
			width: 640px;
			margin-top: 40px;
			margin-left: auto;
			margin-right: auto;
			background: white;
			color: #369;
		}
		</style>
		</head><body>
			<div>
				<h1>%v - %s</h1>
				<p>Kengal 0.9.1</p>
			</div>
		</body>
	</html>`

type ServerError struct {
	Code int
	Msg  string
}

func (se *ServerError) Write(w http.ResponseWriter) {
	w.WriteHeader(se.Code)
	errOut := fmt.Sprintf(errHtml, se.Code, se.Msg)
	w.Write([]byte(errOut))
	w.Flush()
}
func Images(w http.ResponseWriter, r *http.Request) {
	imagePath := path.Base(r.URL.Path)
	mimeType := mime.TypeByExtension(path.Ext(imagePath))

	w.SetHeader("Content-Type", mimeType)
	w.SetHeader("Cache-Control", "max-age=31104000, public")
	current := View.Themes.Current()
	if current == nil{
		se := &ServerError{404, "Not Found"}
		se.Write(w)
		return
	}
	for _, v := range View.Resources {
		if v.Template == current.ID {
			if v.Name == imagePath {
				w.Write(v.Data)
				w.Flush()
				return
			}
		}
	}
	se := &ServerError{404, "Not Found"}
	se.Write(w)
	return
}
func GlobalController(w http.ResponseWriter, r *http.Request) {
	imagePath := path.Base(r.URL.Path)
	mimeType := mime.TypeByExtension(path.Ext(imagePath))
	w.SetHeader("Content-Type", mimeType)
	for _, v := range View.Globals {
		if v.Name == imagePath {
			w.Write(v.Data)
			w.Flush()
		}
	}
}
func Css(w http.ResponseWriter, r *http.Request) {
	w.SetHeader("Content-Encoding", "gzip")
	w.SetHeader("Content-Type", "text/css")

	gz, _ := gzip.NewWriter(w)
	gz.Write([]byte(View.Themes.Current().Style))
	gz.Close()
}
func DeleteLog(w http.ResponseWriter, r *http.Request) {
	View.Master.ClearLog()
	http.Redirect(w, r, r.Referer, http.StatusTemporaryRedirect)
}

func FileHelper(w http.ResponseWriter, r *http.Request) {
	mimeType := mime.TypeByExtension(path.Ext(r.URL.Path))
	w.SetHeader("Content-Encoding", "gzip")
	w.SetHeader("Content-Type", mimeType)
	//w.SetHeader("Expires", "Fri, 30 Oct 2013 14:19:41 GMT")
	b, _ := ioutil.ReadFile("." + r.URL.Path)
	gz, _ := gzip.NewWriter(w)
	gz.Write(b)
	gz.Close()
}
func AdminDispatch(w http.ResponseWriter, kind string) {
	w.SetHeader("Content-Type", "text/html; charset=utf-8")
	w.SetHeader("Content-Encoding", "gzip")
	var Templ = template.New(nil)
	admintemplate, _ := ioutil.ReadFile(fmt.Sprintf("html/%s.html", kind))

	err := Templ.Parse(string(admintemplate))
	if err != nil {
		fmt.Println(err)
	}

	bufNozip := bytes.NewBufferString("")
	err = Templ.Execute(bufNozip, View)
	gz, _ := gzip.NewWriter(w)
	gz.Write(bufNozip.Bytes())
	gz.Close()
	if err != nil {
		fmt.Println(err)
	}
}
func AdminController(w http.ResponseWriter, r *http.Request) {
	View.Blog = 0
	View.Rubric = 0
	View.Article = 0
	View.Server = 0
	View.Theme = 0
	View.Global = 0
	View.Resource = 0

	rd, err := r.MultipartReader()
	if err == nil {
		r.Form = make(map[string][]string)
		for {
			pt, _ := rd.NextPart()
			if pt == nil {
				break
			}
			var fh [1]string
			fd, _ := ioutil.ReadAll(pt)
			fh[0] = strings.TrimSpace(string(fd))
			r.Form[strings.TrimSpace(pt.FormName())] = fh[0:1]
		}
	}
	// Originalpfad der Url zwischenspeichern und nach Redirect widerherstellen
	route := r.FormValue("Route")
	if route != "" {
		orig := r.URL.Path
		r.URL.Path = route
		View.Master.HandleForm(route, w, r)
		r.URL.Path = orig
	}

	dir, file := path.Split(r.URL.Path)
	ids := strings.Split(file, ",", -1)
	
	if ! sort.IsSorted(View.Blogs){
		sort.Sort(View.Blogs)
	}
	if ! sort.IsSorted(View.Themes){
		sort.Sort(View.Themes)
	}

	kind := strings.Replace(dir, "/", "", -1)
	switch kind {
	case "blogs", "newrubrics":
		View.Blog, _ = strconv.Atoi(ids[0])
	case "rubrics", "newarticles":
		View.Blog, _ = strconv.Atoi(ids[0])
		View.Rubric, _ = strconv.Atoi(ids[1])
	case "articles":
		View.Blog, _ = strconv.Atoi(ids[0])
		View.Rubric, _ = strconv.Atoi(ids[1])
		View.Article, _ = strconv.Atoi(ids[2])
	case "servers":
		View.Server, _ = strconv.Atoi(ids[0])
	case "globals":
		View.Global, _ = strconv.Atoi(ids[0])
	case "themes", "newresources":
		View.Theme, _ = strconv.Atoi(ids[0])
	case "resources":
		View.Theme, _ = strconv.Atoi(ids[0])
		View.Resource, _ = strconv.Atoi(ids[1])
	case "":
		kind = "admin"
	}

	w.SetHeader("Content-Type", "text/html; charset=utf-8")
	w.SetHeader("Content-Encoding", "gzip")
	AdminDispatch(w, kind)
}
