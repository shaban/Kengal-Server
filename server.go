package main

import (
	"fmt"
	"template"
	"bytes"
	"strings"
	"http"
	"io/ioutil"
	"compress/gzip"
	//"os"
	"mime"
	"path"
)

func Dispatch(layout string, w http.ResponseWriter) {
	w.SetHeader("Content-Type", "text/html; charset=utf-8")
	w.SetHeader("Content-Encoding", "gzip")
	var Templ = template.New(nil)

	err := Templ.Parse(layout)
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
func Controller(w http.ResponseWriter, r *http.Request) {
	w.SetHeader("Content-Type", "text/html; charset=utf-8")
	w.SetHeader("Content-Encoding", "gzip")

	err := View.loadBlogData(r.Host)
	if err != nil {
		fmt.Println(err.String())
		//View.Error = err.String()
	}

	View.HeadMeta = fmt.Sprintf("<meta name=\"description\" content=\"%s\" />\n", View.Domain.Description)
	View.HeadMeta += fmt.Sprintf("<meta name=\"keywords\" content=\"%s\" />", View.Domain.Keywords)

	Dispatch(View.Template.Index, w)
}
func RubricController(w http.ResponseWriter, r *http.Request) {
	err := View.loadBlogData(r.Host)
	params := strings.Split(r.URL.Path, "/", -1)
	err = View.loadArticlesInRubric(params[2])

	if err != nil {
		fmt.Println(err.String())
		//View.Error = "<h1>404</h1><p>Datei nicht gefunden<br/><a href='/'>Zur Startseite</a></p>"
		//Dispatch("rubric", w)
		return
	}
	for _, rub := range View.Rubrics {
		if rub.ShortUrl == params[3] {
			View.MyRubric = rub
			break
		}
	}
	View.HeadMeta = fmt.Sprintf("<meta name=\"description\" content=\"%s\" />\n", View.MyRubric.Description)
	View.HeadMeta += fmt.Sprintf("<meta name=\"keywords\" content=\"%s\" />", View.MyRubric.Keywords)
	Dispatch(View.Template.Rubric, w)
}
func ArticleController(w http.ResponseWriter, r *http.Request) {
	err := View.loadBlogData(r.Host)
	params := strings.Split(r.URL.Path, "/", -1)
	err = View.loadMyArticle(params[2])
	if err != nil {
		//View.Error = "<h1>404</h1><p>Datei nicht gefunden<br/><a href='/'>Zur Startseite</a></p>"
		//fmt.Println(View.Error)
		Dispatch("article", w)
		return
	}
	View.HeadMeta = fmt.Sprintf("<meta name=\"description\" content=\"%s\" />\n", View.MyArticle.Description)
	View.HeadMeta += fmt.Sprintf("<meta name=\"keywords\" content=\"%s\" />", View.MyArticle.Keywords)

	Dispatch(View.Template.Article, w)
}
func ImprintController(w http.ResponseWriter, r *http.Request) {
	View.loadBlogData(r.Host)
	Dispatch(View.Template.Imprint, w)
}
func Images(w http.ResponseWriter, r *http.Request) {
	imagePath := path.Base(r.URL.Path)
	mimeType := mime.TypeByExtension(path.Ext(imagePath))

	w.SetHeader("Content-Type", mimeType)
	//w.SetHeader("Cache-Control", "public")
	for _, v := range View.Template.Resources {
		if v.Name == imagePath {
			w.Write(v.Data)
			w.Flush()
		}
	}
}
func GlobalController(w http.ResponseWriter, r *http.Request) {
	imagePath := path.Base(r.URL.Path)
	mimeType := mime.TypeByExtension(path.Ext(imagePath))
	w.SetHeader("Content-Type", mimeType)
	for _, v := range View.Template.Global {
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
	gz.Write([]byte(View.Template.Style))
	gz.Close()
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
