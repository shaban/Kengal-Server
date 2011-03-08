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
func DataController(w http.ResponseWriter, r *http.Request) {
	j := r.FormValue("json")
	switch j {
	case "servers":
		data, _ := json.Marshal(View.Servers)
		w.Write(data)
		w.Flush()
	case "blogs":
		data, _ := json.Marshal(View.Blogs)
		w.Write(data)
		w.Flush()
	case "articles":
		data, _ := json.Marshal(View.Articles.Index())
		w.Write(data)
		w.Flush()
	case "templates":
		data, _ := json.Marshal(View.Themes)
		w.Write(data)
		w.Flush()
	case "template":
		data, _ := json.Marshal(View.Themes.Current())
		w.Write(data)
		w.Flush()
	case "resources":
		data, _ := json.Marshal(View.Resources)
		w.Write(data)
		w.Flush()
	case "globals":
		data, _ := json.Marshal(View.Globals)
		w.Write(data)
		w.Flush()
	case "rubrics":
		data, _ := json.Marshal(View.Rubrics.Index())
		w.Write(data)
		w.Flush()
	}
}

func AdminController(w http.ResponseWriter, r *http.Request) {
	//fmt.Println(r.FormValue("server"))
	i, err := strconv.Atoi(r.FormValue("server"))
	if err != nil {
		View.Server = app.Server
	} else {
		View.Server = i
	}
	//ParseParameters(r.URL.Path, r.Host)
	w.SetHeader("Content-Type", "text/html; charset=utf-8")
	w.SetHeader("Content-Encoding", "gzip")
	AdminDispatch(w)
}
