package main

import (
	"fmt"
	"os"
	"http"
	"path"
	gobzip "github.com/shaban/kengal/gobzip"
)
var View = new(Page)
var PaginatorMax = 5

func (a *Article) getBlog() *Blog {
	for k, v := range View.Blogs {
		if v.ID == a.Blog {
			return View.Blogs[k]
		}
	}
	return nil
}
func (r *Rubric) getBlog() *Blog {
	for k, v := range View.Blogs {
		if v.ID == r.Blog {
			return View.Blogs[k]
		}
	}
	return nil
}

func (a Articles) Paginated() []*Article {
	if View.Index == 0 {
		return nil
	}
	l := len(a)
	if l < PaginatorMax {
		return a
	}
	if (View.Index-1)*PaginatorMax+PaginatorMax > l {
		return a[(View.Index-1)*PaginatorMax : l]
	}
	return a[(View.Index-1)*PaginatorMax : (View.Index-1)*PaginatorMax+PaginatorMax]
}
func (a Articles) Prev() string {
	if View.Index <= 1 {
		return ""
	}
	return fmt.Sprintf("/index/%v", View.Index-1)
}
func (a Articles) Next() string {
	if View.Index == 0 {
		return ""
	}
	if len(a) > View.Index*PaginatorMax {
		return fmt.Sprintf("/index/%v", View.Index+1)
	}
	return ""
}
func (a *Article) DateTime() string {
	return a.Date
}
func (a *Article) Path() string {
	return fmt.Sprintf("/artikel/%v/%s", a.ID, a.Url)
}
func (a *Article) RubricPath() string {
	for _, v := range View.Rubrics {
		if v.ID == a.Rubric {
			return v.Path()
		}
	}
	return ""
}

func (a *Article) RubricTitle() string {
	for _, v := range View.Rubrics {
		if v.ID == a.Rubric {
			return v.Title
		}
	}
	return ""
}
func (r *Rubric) Path() string {
	return fmt.Sprintf("/kategorie/%v/%s", r.ID, r.Url)
}

func RubricByID(id int) *Rubric {
	for k, v := range View.Rubrics {
		if v.ID == id {
			return View.Rubrics[k]
		}
	}
	return nil
}

func ArticleByID(id int) *Article {
	for k, v := range View.Articles {
		if v.ID == id {
			return View.Articles[k]
		}
	}
	return nil
}
func BlogByID(id int) *Blog {
	for k, v := range View.Blogs {
		if v.ID == id {
			return View.Blogs[k]
		}
	}
	return nil
}

func (s Servers) Current() *Server {
	for k, v := range s {
		if v.ID == View.Server {
			return s[k]
		}
	}
	return nil
}
func (s Globals) Current() *Global {
	for k, v := range s {
		if v.ID == View.Global {
			return s[k]
		}
	}
	return nil
}
func (s Resources) Current() *Resource {
	for k, v := range s {
		if v.ID == View.Resource {
			return s[k]
		}
	}
	return nil
}
func (s Resources) Index() []*Resource {
	slice := make([]*Resource, 0)
	for k, v := range s {
		if v.Template == View.Theme {
			slice = append(slice, s[k])
		}
	}
	return slice
}

func (s *Global) DataString() string {
	switch path.Ext(s.Name){
	case ".js",".css","html",".htm":
		return string(s.Data)
	default:
		return ""
	}
	return ""
}

func (srv *Server) Active() bool {
	if View.Server == 0 || srv.ID != View.Server {
		return false
	}
	return true
}

func (t *Theme) Active() bool {
	if View.Blogs.Current().Template == t.ID {
		return true
	}
	return false
}

func (b *Blog) Active() bool {
	if b.ID == View.Blogs.Current().ID {
		return true
	}
	return false
}

func (r *Rubric) Active() bool {
	if r.ID == View.Rubrics.Current().ID {
		return true
	}
	return false
}
func (b Blogs) Current() *Blog {
	for k, v := range b {
		if v.ID == View.Blog {
			return b[k]
		}
	}
	return nil
}
func (t Themes) Current() *Theme {
	for k, v := range t {
		if v.ID == View.Theme {
			return t[k]
		}
	}
	return nil
}

func (a Articles) Latest() []*Article {
	l := len(a)
	if l < 5 {
		return a
	}
	return a[0:5]
}
func (a Articles) Index() Articles {
	b := View.Blogs.Current()
	s := make([]*Article, 0)
	for k, v := range a {
		if b.ID == v.Blog {

			s = append(s, a[k])
		}
	}
	return s
}

func (a Articles) Current() *Article {
	if View.Article == 0 {
		return nil
	}
	for k, v := range a {
		if v.ID == View.Article {
			return a[k]
		}
	}
	return nil
}

func (a Articles) Rubric() []*Article {
	if View.Rubric == 0 {
		return nil
	}
	s := make([]*Article, 0)
	for k, v := range a {
		if v.Rubric == View.Rubric {
			s = append(s, a[k])
		}
	}
	if len(s) == 0 {
		return nil
	}
	return s
}
func (r Rubrics) Current() *Rubric {
	if View.Rubric == 0 {
		return nil
	}
	for k, v := range r {
		if v.ID == View.Rubric {
			return r[k]
		}
	}
	return nil
}
func (r Rubrics) Index() Rubrics {
	b := View.Blogs.Current()
	s := make([]*Rubric, 0)
	for k, v := range r {
		if v.Blog == b.ID {
			s = append(s, r[k])
		}
	}
	return s
}
func main() {
	View.Master = gobzip.DefaultMaster
	View.Master.Init(View, "db", "/admin/delete/", "/admin/replace/", "/admin/insert/", "/admin/audit/")
	View.Master.HandleForms()

	err := LoadAll()
	if err != nil {
		fmt.Println(err.String())
		os.Exit(1)
	}
	http.HandleFunc("/", AdminController)

	http.HandleFunc("/global/", GlobalController)
	http.HandleFunc("/images/", Images)
	http.HandleFunc("/style.css", Css)
	http.HandleFunc("/favicon.ico", GlobalController)

	http.HandleFunc("/js/", FileHelper)
	http.HandleFunc("/css/", FileHelper)
	http.HandleFunc("/html/", FileHelper)
	http.HandleFunc("/admin/clear/log/",DeleteLog)

	http.HandleFunc("/ckeditor/", FileHelper)

	http.ListenAndServe(":80", nil)
}
