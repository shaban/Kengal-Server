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
func BlogSave(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	b := new(Blog)
	f := r.Form
	
	b.Title=f["Title"][0]
	b.Slogan=f["Slogan"][0]
	b.Keywords=f["Keywords"][0]
	b.Description=f["Description"][0]
	b.Template,_=strconv.Atoi(f["Template"][0])
	b.ID,_=strconv.Atoi(f["ID"][0])
	b.Server,_=strconv.Atoi(f["Server"][0])
	b.Url=f["Url"][0]
	
	err := statement.UpdateBlog.BindParams(b.Description, b.Keywords, b.Server, b.Slogan, b.Template, b.Title, b.Url, b.ID)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Speichern von Blog %s fehlgeschlagen",b.Title)))
		return
	}
	err = statement.UpdateBlog.Execute()
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Speichern von Blog %s fehlgeschlagen",b.Title)))
		return
	}
	View.Blogs.Replace(b)
	w.Write([]byte(fmt.Sprintf("Blog %s erfolgreich gespeichert", b.Title)))
}

func RubricSave(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	rb := new(Rubric)
	f := r.Form
	fmt.Println(f)
	rb.Title=f["Title"][0]
	rb.Keywords=f["Keywords"][0]
	rb.Description=f["Description"][0]
	rb.Blog,_=strconv.Atoi(f["Blog"][0])
	rb.ID,_=strconv.Atoi(f["ID"][0])
	rb.Url=f["Url"][0]
	
	err := statement.UpdateRubric.BindParams(rb.Description, rb.Keywords, rb.Blog, rb.Title, rb.Url, rb.ID)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Speichern von Rubrik %s fehlgeschlagen | BindParams",rb.Title)))
		return
	}
	err = statement.UpdateRubric.Execute()
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Speichern von Rubrik %s fehlgeschlagen | Execute",rb.Title)))
		return
	}
	View.Rubrics.Replace(rb)
	w.Write([]byte(fmt.Sprintf("Rubrik %s erfolgreich gespeichert", rb.Title)))
}

func ArticleSave(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	a := new(Article)
	f := r.Form
	fmt.Println(f["Blog"][0])
	a.Title=f["Title"][0]
	a.Keywords=f["Keywords"][0]
	a.Description=f["Description"][0]
	a.Rubric,_=strconv.Atoi(f["Rubric"][0])
	a.ID,_=strconv.Atoi(f["ID"][0])
	a.Url=f["Url"][0]
	a.Blog,_=strconv.Atoi(f["Blog"][0])
	a.Date = f["Date"][0]
	a.Teaser = f["Teaser"][0]
	a.Text = f["Text"][0]
	
	err := statement.UpdateArticle.BindParams(a.Description, a.Keywords, a.Blog, a.Rubric, a.Date, a.Title, a.Url, a.Text, a.Teaser, a.ID)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Speichern von Artikel %s fehlgeschlagen | BindParams",a.Title)))
		return
	}
	err = statement.UpdateArticle.Execute()
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Speichern von Artikel %s fehlgeschlagen | Execute",a.Title)))
		return
	}
	View.Articles.Replace(a)
	w.Write([]byte(fmt.Sprintf("Artikel %s erfolgreich gespeichert", a.Title)))
}

func BlogNew(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	b := new(Blog)
	f := r.Form
	
	b.Title=f["Title"][0]
	b.Slogan=f["Slogan"][0]
	b.Keywords=f["Keywords"][0]
	b.Description=f["Description"][0]
	b.Template,_=strconv.Atoi(f["Template"][0])
	b.Server,_=strconv.Atoi(f["Server"][0])
	b.Url=f["Url"][0]
	
	err := statement.InsertBlog.BindParams(b.Description, b.Keywords, b.Server, b.Slogan, b.Template, b.Title, b.Url)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Anlegen von Blog %s fehlgeschlagen",b.Title)))
		return
	}
	err = statement.InsertBlog.Execute()
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Anlegen von Blog %s fehlgeschlagen",b.Title)))
		return
	}
	b.ID = int(statement.InsertBlog.LastInsertId)
	View.Blogs = append(View.Blogs,b)
	w.Write([]byte(fmt.Sprintf("Blog %s erfolgreich angelegt", b.Title)))
}

func RubricNew(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	rb := new(Rubric)
	f := r.Form
	fmt.Println(f)
	rb.Title=f["Title"][0]
	rb.Keywords=f["Keywords"][0]
	rb.Description=f["Description"][0]
	rb.Blog,_=strconv.Atoi(f["Blog"][0])
	rb.Url=f["Url"][0]
	
	err := statement.InsertRubric.BindParams(rb.Description, rb.Keywords, rb.Blog, rb.Title, rb.Url)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Anlegen von Rubrik %s fehlgeschlagen",rb.Title)))
		return
	}
	err = statement.InsertRubric.Execute()
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Anlegen von Rubrik %s fehlgeschlagen",rb.Title)))
		return
	}
	rb.ID = int(statement.InsertRubric.LastInsertId)
	View.Rubrics = append(View.Rubrics,rb)
	w.Write([]byte(fmt.Sprintf("Rubrik %s erfolgreich angelegt", rb.Title)))
}

func ArticleNew(w http.ResponseWriter, r *http.Request) {
}


func CommandUnit(w http.ResponseWriter, r *http.Request) {	
	/*err := View.loadBlogData()
	if err != nil {
		fmt.Println(err.String())
		os.Exit(1)
	}*/
}
func Controller(w http.ResponseWriter, r *http.Request) {
	ParseParameters(r.URL.Path, r.Host)
	if r.FormValue("mode") == "admin"{
		AdminController(w,r)
		return
	}
	
	if r.FormValue("mode") == "snippet"{
		SnippetController(w,r)
		return
	}
	
	w.SetHeader("Content-Type", "text/html; charset=utf-8")
	
	if View.Blogs.Current() == nil{
		kw := new(KengalWebError)
		kw.Code = 403
		kw.Msg = "Zugriff nicht erlaubt"
		kw.Write(w)
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
	kw := new(KengalWebError)
		kw.Code = 404
		kw.Msg = "Möglicherweise haben Sie eine falsche Url eingegeben"
		kw.Write(w)
		return
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
