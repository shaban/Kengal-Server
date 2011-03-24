package main

import (
	"http"
	"compress/gzip"
	"template"
	"fmt"
	"bytes"
	"io/ioutil"
	"strconv"
	"json"
)

func AdminDispatch(w http.ResponseWriter) {
	w.SetHeader("Content-Type", "text/html; charset=utf-8")
	w.SetHeader("Content-Encoding", "gzip")
	var Templ = template.New(nil)
	admintemplate, _ := ioutil.ReadFile("admtpl.html")

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

func SnippetController(w http.ResponseWriter, r *http.Request) {
	j := r.FormValue("which")
	
	s, _ := strconv.Atoi(r.FormValue("server"))
	h := r.FormValue("host")
	a, _ := strconv.Atoi(r.FormValue("article"))
	rb, _ := strconv.Atoi(r.FormValue("rubric"))

	os := View.Server
	oh := View.Host
	oa := View.Article
	orb := View.Rubric

	View.Server = s
	View.Host = h
	View.Rubric = rb
	View.Article = a

	tplFile := ""
	w.SetHeader("Content-Type", "text/html; charset=utf-8")

	switch j {
	case "blogs":
		tplFile = "html/blogSnippet.html"
	case "articles":
		tplFile = "html/articleSnippet.html"
	case "rubrics":
		tplFile = "html/rubricSnippet.html"
	case "editor":
		tplFile = "html/editorSnippet.html"
	case "templates":
		data, _ := json.Marshal(View.Themes)
		w.Write(data)
		w.Flush()
	case "newblog":
		tplFile = "html/blogNewSnippet.html"
	case "newrubric":
		tplFile = "html/rubricNewSnippet.html"
	case "newarticle":
		tplFile = "html/articleNewSnippet.html"
	}
	var Templ = template.New(nil)
	snippetTempl, err := ioutil.ReadFile(tplFile)
	if err != nil {
		fmt.Println(err)
	}
	buf := bytes.NewBufferString("")
	err = Templ.Parse(string(snippetTempl))
	if err != nil {
		fmt.Println(err)
	}
	err = Templ.Execute(buf, View)
	if err != nil {
		fmt.Println(err)
	}
	w.Write(buf.Bytes())
	View.Server = os
	View.Host = oh
	View.Article = oa
	View.Rubric = orb
}

func AdminController(w http.ResponseWriter, r *http.Request) {
	i, err := strconv.Atoi(r.FormValue("server"))
	if err != nil {
		View.Server = app.Server
	} else {
		app.Server = View.Server
		View.Server = i
	}
	//ParseParameters(r.URL.Path, r.Host)
	w.SetHeader("Content-Type", "text/html; charset=utf-8")
	w.SetHeader("Content-Encoding", "gzip")
	AdminDispatch(w)
}