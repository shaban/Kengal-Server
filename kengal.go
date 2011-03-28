package main

import (
	"fmt"
	"os"
	//"flag"
	"http"
	gobzip "github.com/shaban/kengal/gobzip"
)

//var app = new(Application)
var master = gobzip.DefaultMaster
var View = new(Page)
var PaginatorMax = 5

func (a *Article)getBlog() *Blog{
	for k, v := range View.Blogs{
		if v.ID == a.Blog{
			return View.Blogs[k]
		}
	}
	return nil
}
func (r *Rubric)getBlog() *Blog{
	for k, v := range View.Blogs{
		if v.ID == r.Blog{
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

func RubricByID(id int)* Rubric{
	for k, v := range View.Rubrics{
		if v.ID == id{
			return View.Rubrics[k]
		}
	}
	return nil
}

func ArticleByID(id int)* Article{
	for k, v := range View.Articles{
		if v.ID == id{
			return View.Articles[k]
		}
	}
	return nil
}
func BlogByID(id int)* Blog{
	for k, v := range View.Blogs{
		if v.ID == id{
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

func (srv *Server) Active() bool {
	if View.Server == 0 || srv.ID != View.Server{
		return false
	}
	return true
}

func (t *Theme) Active() bool {
	if View.Blogs.Current().Template == t.ID{
		return true
	}
	return false
}

func (b *Blog) Active() bool {
	if b.ID == View.Blogs.Current().ID{
		return true
	}
	return false
}

func (r *Rubric) Active() bool {
	if r.ID == View.Rubrics.Current().ID{
		return true
	}
	return false
}



func (b Blogs) Current() *Blog {
	for k, v := range b {
		if v.Url == View.Host {
			return b[k]
		}
	}
	return nil
}

func (t Themes) Current() *Theme {
	current := View.Blogs.Current()
	for k, v := range t {
		if v.ID == current.Template {
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
/*type Application struct {
	User     string
	Password string
	Database string
	LogLevel int
	Server   int
}*/
func main() {
	//flag.IntVar(&app.Server, "s", 0, "Geben Sie hier die ID des Servers an")
	//flag.StringVar(&master.Root, "f", "db", "Geben Sie hier die Wurzel der Datenbank an")

	/*flag.Parse()

	if app.Server == 0 {
		flag.Usage()
		os.Exit(0)
	}*/
	
	View.Server=0
	master.Delegator = View
	master.Init("db","/admin/delete/","/admin/replace/","/admin/insert/","/admin/audit/")
	master.HandleForms()

	err := LoadAll()
	if err != nil {
		fmt.Println(err.String())
		os.Exit(1)
	}
	//View.Articles = make([]*Article,0)
	//var err os.Error
	
	
	for _, v := range View.Articles{
		fmt.Printf("Art\tk: %v v:%s\n",v.ID,v.Title)
	}
	for _, v := range View.Blogs{
		fmt.Printf("Blg\tk: %v v:%s\n",v.ID,v.Title)
	}
	for _, v := range View.Globals{
		fmt.Printf("Glb\tk: %v v:%s\n",v.ID,v.Name)
	}
	for _, v := range View.Resources{
		fmt.Printf("Rsr\tk: %v v:%s\n",v.ID,v.Name)
	}
	for _, v := range View.Rubrics{
		fmt.Printf("Rbr\tk: %v v:%s\n",v.ID,v.Title)
	}
	for _, v := range View.Servers{
		fmt.Printf("Srv\tk: %v v:%s\n",v.ID,v.Vendor)
	}
	for _, v := range View.Themes{
		fmt.Printf("Thm\tk: %v v:%s\n",v.ID,v.Title)
	}
	//os.Exit(0)

	/*http.HandleFunc("/admin/blog/save", BlogSave)
	http.HandleFunc("/admin/rubric/save", RubricSave)
	http.HandleFunc("/admin/article/save", ArticleSave)*/
	
	http.HandleFunc("/snippet/", SnippetController)
	http.HandleFunc("/", AdminController)
	
	//http.HandleFunc("/admin/audit/",Audit)
	
	/*http.HandleFunc("/admin/blog/new", BlogNew)
	http.HandleFunc("/admin/rubric/new", RubricNew)
	http.HandleFunc("/admin/article/new", ArticleNew)
	
	http.HandleFunc("/admin/blog/delete/", BlogDelete)
	http.HandleFunc("/admin/rubric/delete/", RubricDelete)
	http.HandleFunc("/admin/article/delete/", ArticleDelete)*/

	http.HandleFunc("/global/", GlobalController)
	http.HandleFunc("/images/", Images)
	http.HandleFunc("/style.css", Css)
	http.HandleFunc("/favicon.ico", GlobalController)

	http.HandleFunc("/js/", FileHelper)
	http.HandleFunc("/css/", FileHelper)
	http.HandleFunc("/html/", FileHelper)
	http.HandleFunc("/tpl/", FileHelper)
	
	http.HandleFunc("/ckeditor/", FileHelper)

	http.ListenAndServe(":6060", nil)
	os.Exit(0)
}
