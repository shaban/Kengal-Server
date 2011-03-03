package main

import (
	"fmt"
	"os"
	"flag"
	"mysql"
	"http"
)

type Server struct {
	ID      int
	IP      string
	Comment string
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
	ShortUrl    string
	Keywords    string
	Description string
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
type Theme struct {
	ID        int
	Article   string
	Index     string
	Imprint   string
	Rubric    string
	Resources []*Resource
	Global    []*Resource
	Error     string
	Style     string
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
type Paginator struct{}

type Page struct {
	HeadMeta   string
	HeadScript string
	HeadLink   string
	HeadStyle  string
	Index      []*Article
	Rubrics    []*Rubric
	Latest     []*Article
	Articles   []*Article
	MyArticle  *Article
	MyRubric   *Rubric
	Domain     *Blog
	Template   *Theme
	Imprint	bool
	Admin bool
}

var View = new(Page)
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
func (a *Article) RubricTitle() string {
	for _, v := range View.Rubrics {
		if v.ID == a.Rubric {
			return v.Title
		}
	}
	return ""
}
func (r *Rubric) Path() string {
	return fmt.Sprintf("/kategorie/%v/%s", r.ID, r.ShortUrl)
}

func main() {
	user := flag.String("u", "root", "set Mysql User here, default is root")
	pw := flag.String("p", "password", "set Mysql Password for selected User here")
	dbname := flag.String("db", "mysql", "set Database that MySql is supposed to connect to here")
	dbhost := flag.String("h", "", "set Host MySql Adress like so -h myserver.com")

	flag.Parse()

	if *pw == "password" || *dbname == "mysql" {
		flag.Usage()
		os.Exit(0)
	}

	err := InitMysql(*dbhost, *user, *pw, *dbname)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	prepareMysql()

	http.HandleFunc("/", Controller)
	http.HandleFunc("/kategorie/", RubricController)
	http.HandleFunc("/artikel/", ArticleController)
	http.HandleFunc("/imprint/", ImprintController)
	http.HandleFunc("/global/", GlobalController)

	http.HandleFunc("/images/", Images)
	http.HandleFunc("/style.css", Css)
	http.HandleFunc("/favicon.ico", GlobalController)
	http.HandleFunc("/js/", FileHelper)
	http.HandleFunc("/ckeditor/", FileHelper)
	http.HandleFunc("/templates/", FileHelper)

	http.ListenAndServe(":8080", nil)
	os.Exit(0)
}
