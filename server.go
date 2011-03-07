package main

import (
	"fmt"
	"template"
	"bytes"
	"http"
	"io/ioutil"
	"compress/gzip"
	"mime"
	"path"
	"os"
	"strconv"
)

func Dispatch(w http.ResponseWriter) {
	w.SetHeader("Content-Type", "text/html; charset=utf-8")
	w.SetHeader("Content-Encoding", "gzip")
	var Templ = template.New(nil)

	err := Templ.Parse(View.Themes.Current().Index)
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

func ParseParameters(url, host string) os.Error {
	View.Host = host
	View.Imprint = false
	View.Index = 0
	View.Rubric = 0
	View.Article = 0

	dir, file := path.Split(url)
	if file == "" {
		// is Index Startpage
		View.Index = 1
		return nil
	}
	i, err := strconv.Atoi(file)
	if err == nil {
		//is Index non Startpage
		View.Index = i
		if ((i - 1) * PaginatorMax) > len(View.Articles.Index()) {
			View.Index = 0
			return nil
		}
		return nil
	}
	if file == "impressum" {
		// is Impressum
		View.Imprint = true
		return nil
	}
	nextdir := path.Clean(dir)
	dir, file = path.Split(nextdir)
	i, err = strconv.Atoi(file)
	if err != nil {
		return err
	}
	if dir == "/kategorie/" {
		//is Rubricpage
		View.Rubric = i
		if View.Rubrics.Index().Current() == nil {
			View.Rubric = 0
		}
		return nil
	}
	if dir == "/artikel/" {
		//is Rubricpage
		View.Article = i
		if View.Articles.Index().Current() == nil {
			View.Article = 0
		}
		return nil
	}
	return os.ENOTDIR
}
func CommandUnit(w http.ResponseWriter, r *http.Request) {
	err := View.loadBlogData()
	if err != nil {
		fmt.Println(err.String())
		os.Exit(1)
	}
}
func Controller(w http.ResponseWriter, r *http.Request) {
	ParseParameters(r.URL.Path, r.Host)
	w.SetHeader("Content-Type", "text/html; charset=utf-8")
	w.SetHeader("Content-Encoding", "gzip")

	if View.Index != 0 {
		View.HeadMeta = fmt.Sprintf(`<meta name="description" content="%s" />`, View.Blogs.Current().Description)
		View.HeadMeta += fmt.Sprintf(`<meta name="keywords" content="%s" />`, View.Blogs.Current().Keywords)
		Dispatch(w)
		return
	}
	if View.Article != 0 {
		View.HeadMeta = fmt.Sprintf(`<meta name="description\" content=\"%s\" />`, View.Articles.Current().Description)
		View.HeadMeta += fmt.Sprintf(`<meta name="keywords" content="%s" />`, View.Articles.Current().Keywords)
		Dispatch(w)
		return
	}
	if View.Rubric != 0 {
		View.HeadMeta = fmt.Sprintf(`<meta name="description" content="%s" />`, View.Rubrics.Current().Description)
		View.HeadMeta += fmt.Sprintf(`<meta name="keywords" content="%s" />`, View.Rubrics.Current().Keywords)
		Dispatch(w)
		return
	}
	if View.Imprint {
		View.HeadMeta = fmt.Sprintf(`<meta name="description" content="%s" />`, "Impressum")
		View.HeadMeta += fmt.Sprintf(`<meta name="keywords" content="%s" />`, "Impressum")
		Dispatch(w)
		return
	}

	bufNozip := bytes.NewBufferString(Error)
	gz, _ := gzip.NewWriter(w)
	gz.Write(bufNozip.Bytes())
	gz.Close()
}
func Images(w http.ResponseWriter, r *http.Request) {
	imagePath := path.Base(r.URL.Path)
	mimeType := mime.TypeByExtension(path.Ext(imagePath))

	w.SetHeader("Content-Type", mimeType)
	w.SetHeader("Cache-Control","max-age=31104000, public")
	current := View.Themes.Current()
	for _, v := range View.Resources {
		if v.Template == current.ID {
			if v.Name == imagePath {
				w.Write(v.Data)
				w.Flush()
			}
		}
	}
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
