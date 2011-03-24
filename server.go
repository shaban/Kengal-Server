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
	"time"
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

type ServerError struct{
	Code int 
	Msg string
}
func (se *ServerError)Write(w http.ResponseWriter){
	w.WriteHeader(se.Code)
	errOut := fmt.Sprintf(errHtml,se.Code,se.Msg)
	w.Write([]byte(errOut))
	w.Flush()
}

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
func BlogSave(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	b := new(Blog)
	f := r.Form

	b.Title = f["Title"][0]
	b.Slogan = f["Slogan"][0]
	b.Keywords = f["Keywords"][0]
	b.Description = f["Description"][0]
	b.Template, _ = strconv.Atoi(f["Template"][0])
	b.ID, _ = strconv.Atoi(f["ID"][0])
	b.Server, _ = strconv.Atoi(f["Server"][0])
	b.Url = f["Url"][0]

	err := updateBlog(b)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Speichern von Blog %s fehlgeschlagen", b.Title)))
		return
	}
	w.Write([]byte(fmt.Sprintf("Blog %s erfolgreich gespeichert", b.Title)))
}

func RubricSave(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	rb := new(Rubric)
	f := r.Form
	fmt.Println(f)
	rb.Title = f["Title"][0]
	rb.Keywords = f["Keywords"][0]
	rb.Description = f["Description"][0]
	rb.Blog, _ = strconv.Atoi(f["Blog"][0])
	rb.ID, _ = strconv.Atoi(f["ID"][0])
	rb.Url = f["Url"][0]

	err := updateRubric(rb)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Speichern von Rubrik %s fehlgeschlagen", rb.Title)))
		return
	}
	w.Write([]byte(fmt.Sprintf("Rubrik %s erfolgreich gespeichert", rb.Title)))
}

func ArticleSave(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	a := new(Article)
	f := r.Form
	fmt.Println(f["Blog"][0])
	a.Title = f["Title"][0]
	a.Keywords = f["Keywords"][0]
	a.Description = f["Description"][0]
	a.Rubric, _ = strconv.Atoi(f["Rubric"][0])
	a.ID, _ = strconv.Atoi(f["ID"][0])
	a.Url = f["Url"][0]
	a.Blog, _ = strconv.Atoi(f["Blog"][0])
	a.Date = f["Date"][0]
	a.Teaser = f["Teaser"][0]
	a.Text = f["Text"][0]

	err := updateArticle(a)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Speichern von Artikel %s fehlgeschlagen", a.Title)))
		return
	}
	w.Write([]byte(fmt.Sprintf("Artikel %s erfolgreich gespeichert", a.Title)))
}

func BlogNew(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	b := new(Blog)
	f := r.Form

	b.Title = f["Title"][0]
	b.Slogan = f["Slogan"][0]
	b.Keywords = f["Keywords"][0]
	b.Description = f["Description"][0]
	b.Template, _ = strconv.Atoi(f["Template"][0])
	b.Server, _ = strconv.Atoi(f["Server"][0])
	b.Url = f["Url"][0]
	b.ID = View.Blogs.ID()

	err := insertBlog(b)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Anlegen von Blog %s fehlgeschlagen", b.Title)))
		return
	}
	w.Write([]byte(fmt.Sprintf("Blog %s erfolgreich angelegt", b.Title)))
}

func RubricNew(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	rb := new(Rubric)
	f := r.Form

	rb.Title = f["Title"][0]
	rb.Keywords = f["Keywords"][0]
	rb.Description = f["Description"][0]
	rb.Blog, _ = strconv.Atoi(f["Blog"][0])
	rb.Url = f["Url"][0]
	rb.ID = View.Rubrics.ID()

	err := insertRubric(rb)

	if err != nil {
		w.Write([]byte(fmt.Sprintf("Anlegen von Rubrik %s fehlgeschlagen", rb.Title)))
		return
	}
	w.Write([]byte(fmt.Sprintf("Rubrik %s erfolgreich angelegt", rb.Title)))
}

func ArticleNew(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	a := new(Article)
	f := r.Form
	fmt.Println(f)

	a.Title = f["Title"][0]
	a.Keywords = f["Keywords"][0]
	a.Description = f["Description"][0]
	a.Blog, _ = strconv.Atoi(f["Blog"][0])
	a.Url = f["Url"][0]
	a.Rubric, _ = strconv.Atoi(f["Rubric"][0])
	a.Teaser = "<p>Geben Sie hier Ihren Teasertext ein</p>"
	a.Text = "<p>Geben Sie hier den Text des Artikels ein</p>"
	a.Date = time.LocalTime().Format("02.01.2006 15:04:05")
	fmt.Println(a.Date)
	a.ID = View.Articles.ID()

	err := insertArticle(a)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Anlegen von Artikel %s fehlgeschlagen", a.Title)))
		return
	}
	w.Write([]byte(fmt.Sprintf("Artikel %s erfolgreich angelegt", a.Title)))
}
func ArticleDelete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.FormValue("ID"))
	err := deleteArticle(id)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Artikel %v konnte nicht gelöscht werden", id)))
	}
	w.Write([]byte(fmt.Sprintf("Artikel %v erfolgreich gelöscht", id)))
}
func RubricDelete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.FormValue("ID"))
	err := deleteRubric(id)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Rubrik %v konnte nicht gelöscht werden", id)))
	}
	w.Write([]byte(fmt.Sprintf("Rubrik %v erfolgreich gelöscht", id)))
}
func Audit(w http.ResponseWriter, r *http.Request) {
	_, file := path.Split(r.URL.Path)
	var data []byte
	var err os.Error
	switch file{
		case "articles":
			data, err = MakeAudit(View.Articles)
		case "blogs":
			data, err = MakeAudit(View.Blogs)
		case "globals":
			data, err = MakeAudit(View.Globals)
		case "resources":
			data, err = MakeAudit(View.Resources)
		case "rubrics":
			data, err = MakeAudit(View.Rubrics)
		case "servers":
			data, err = MakeAudit(View.Servers)
		case "themes":
			data, err = MakeAudit(View.Themes)
	}
	if err != nil{
		w.WriteHeader(500)
		return
	}
	gz, err := gzip.NewWriter(w)
	if err != nil{
		w.WriteHeader(500)
		return
	}
	w.SetHeader("Content-Type", "application/json; charset=utf-8")
	w.SetHeader("Content-Encoding", "gzip")
	gz.Write(data)
	gz.Close()
}

func Controller(w http.ResponseWriter, r *http.Request) {
	ParseParameters(r.URL.Path, r.Host)
	if r.FormValue("mode") == "admin" {
		AdminController(w, r)
		return
	}

	if r.FormValue("mode") == "snippet" {
		SnippetController(w, r)
		return
	}

	w.SetHeader("Content-Type", "text/html; charset=utf-8")

	if View.Blogs.Current() == nil {
		se := &ServerError{403, "Forbidden"}
		se.Write(w)
		return
	}

	if View.Index != 0 {
		w.SetHeader("Content-Encoding", "gzip")
		View.HeadMeta = fmt.Sprintf(`<meta name="description" content="%s" />`, View.Blogs.Current().Description)
		View.HeadMeta += fmt.Sprintf(`<meta name="keywords" content="%s" />`, View.Blogs.Current().Keywords)
		Dispatch(w)
		return
	}
	if View.Article != 0 {
		w.SetHeader("Content-Encoding", "gzip")
		View.HeadMeta = fmt.Sprintf(`<meta name="description\" content=\"%s\" />`, View.Articles.Current().Description)
		View.HeadMeta += fmt.Sprintf(`<meta name="keywords" content="%s" />`, View.Articles.Current().Keywords)
		Dispatch(w)
		return
	}
	if View.Rubric != 0 {
		w.SetHeader("Content-Encoding", "gzip")
		View.HeadMeta = fmt.Sprintf(`<meta name="description" content="%s" />`, View.Rubrics.Current().Description)
		View.HeadMeta += fmt.Sprintf(`<meta name="keywords" content="%s" />`, View.Rubrics.Current().Keywords)
		Dispatch(w)
		return
	}
	if View.Imprint {
		w.SetHeader("Content-Encoding", "gzip")
		View.HeadMeta = fmt.Sprintf(`<meta name="description" content="%s" />`, "Impressum")
		View.HeadMeta += fmt.Sprintf(`<meta name="keywords" content="%s" />`, "Impressum")
		Dispatch(w)
		return
	}
	se := &ServerError{404, "Not Found"}
	se.Write(w)
	return
}
func Images(w http.ResponseWriter, r *http.Request) {
	imagePath := path.Base(r.URL.Path)
	mimeType := mime.TypeByExtension(path.Ext(imagePath))

	w.SetHeader("Content-Type", mimeType)
	w.SetHeader("Cache-Control", "max-age=31104000, public")
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
