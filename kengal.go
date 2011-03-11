package main

import (
	"fmt"
	"os"
	"flag"
	"mysql"
	"http"
)

type Servers []*Server
type Articles []*Article
type Rubrics []*Rubric
type Blogs []*Blog
type Themes []*Theme
type Resources []*Resource

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

func (a Articles) Next() string {
	if View.Index == 0 {
		return ""
	}
	if len(a) > View.Index*PaginatorMax {
		return fmt.Sprintf("/index/%v", View.Index+1)
	}
	return ""
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

type Server struct {
	ID     int
	IP     string
	Vendor string
	Type   string
	Item   int
}

type Blog struct {
	ID          int
	Title       string
	Url         string
	Template    int
	Keywords    string
	Description string
	Slogan      string
	Server      int
}

type Rubric struct {
	ID          int
	Title       string
	Url    string
	Keywords    string
	Description string
	Blog        int
}

type Article struct {
	ID          int
	Date        mysql.DateTime
	Title       string
	Keywords    string
	Description string
	Text        string
	Teaser      string
	Blog        int
	Rubric      int
	Url         string
}

type Resource struct {
	Name     string
	Template int
	Data     []byte
}

type BlogError struct {
	Code int
	Msg  string
}

type Page struct {
	HeadMeta  string
	Rubrics   Rubrics
	Articles  Articles
	Blogs     Blogs
	Blog      int
	Themes    Themes
	Resources Resources
	Globals   Resources
	Servers   Servers
	Index     int
	Rubric    int
	Article   int
	Server    int
	Imprint   bool
	Host      string
}
type Theme struct {
	ID      int
	Index   string
	Style   string
	Title   string
	FromUrl string
}
type Application struct {
	User     string
	Password string
	Database string
	DataHost string
	Server   int
}

var app = new(Application)
var View = new(Page)
var PaginatorMax = 5

func (a *Article) DateTime() string {
	return a.Date.String()
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
func (a *Article) ptrRubric() *Rubric {
	for k, v := range View.Rubrics {
		if v.ID == a.Rubric {
			return View.Rubrics[k]
		}
	}
	return nil
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

func main() {
	flag.StringVar(&app.User, "u", "root", "set Mysql User here, default is root")
	flag.StringVar(&app.Password, "p", "password", "set Mysql Password for selected User here")
	flag.StringVar(&app.Database, "db", "mysql", "set Database that MySql is supposed to connect to here")
	flag.StringVar(&app.DataHost, "h", "", "set Host MySql Adress like so -h myserver.com")
	flag.IntVar(&app.Server, "s", 0, "Set Server ID here")

	flag.Parse()

	if app.Password == "password" || app.Database == "mysql" || app.Server == 0 {
		flag.Usage()
		os.Exit(0)
	}

	err := View.loadBlogData()
	if err != nil {
		fmt.Println(err.String())
		os.Exit(1)
	}
	//popdb()
	http.HandleFunc("/", Controller)
	//http.HandleFunc("/admin/", AdminController)

	http.HandleFunc("/command/", CommandUnit)

	http.HandleFunc("/global/", GlobalController)
	http.HandleFunc("/images/", Images)
	http.HandleFunc("/style.css", Css)
	http.HandleFunc("/favicon.ico", GlobalController)
	//http.HandleFunc("/snippet/", DataController)
	//http.HandleFunc("/spry/", FileHelper)

	http.HandleFunc("/js/", FileHelper)
	http.HandleFunc("/css/", FileHelper)
	http.HandleFunc("/html/", FileHelper)
	
	http.HandleFunc("/ckeditor/", FileHelper)
	//http.HandleFunc("/templates/", FileHelper)

	http.ListenAndServe(":80", nil)
	os.Exit(0)
}
