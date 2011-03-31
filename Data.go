package main

import (
	"os"
	gobzip "github.com/shaban/kengal/gobzip"
	"strconv"
	"time"
	"fmt"
)

type Servers []*Server
type Articles []*Article
type Rubrics []*Rubric
type Blogs []*Blog
type Themes []*Theme
type Resources []*Resource
type Globals []*Global

func (s *Article) Key() int {
	return s.ID
}
func (s *Blog) Key() int {
	return s.ID
}
func (s *Global) Key() int {
	return s.ID
}
func (s *Resource) Key() int {
	return s.ID
}
func (s *Rubric) Key() int {
	return s.ID
}
func (s *Server) Key() int {
	return s.ID
}
func (s *Theme) Key() int {
	return s.ID
}
func (s *Article) Kind() string {
	return "articles"
}
func (s *Blog) Kind() string {
	return "blogs"
}
func (s *Global) Kind() string {
	return "globals"
}
func (s *Resource) Kind() string {
	return "resources"
}
func (s *Rubric) Kind() string {
	return "rubrics"
}
func (s *Server) Kind() string {
	return "servers"
}
func (s *Theme) Kind() string {
	return "themes"
}
func (ser Articles) Kind() string {
	return "articles"
}
func (ser Blogs) Kind() string {
	return "blogs"
}
func (ser Globals) Kind() string {
	return "globals"
}
func (ser Resources) Kind() string {
	return "resources"
}
func (ser Rubrics) Kind() string {
	return "rubrics"
}
func (e Servers) Kind() string {
	return "servers"
}
func (ser Themes) Kind() string {
	return "themes"
}
func (ser Articles) New() gobzip.Serial {
	return new(Article)
}
func (ser Blogs) New() gobzip.Serial {
	return new(Blog)
}
func (ser Globals) New() gobzip.Serial {
	return new(Global)
}
func (ser Resources) New() gobzip.Serial {
	return new(Resource)
}
func (ser Rubrics) New() gobzip.Serial {
	return new(Rubric)
}
func (ser Servers) New() gobzip.Serial {
	return new(Server)
}
func (ser Themes) New() gobzip.Serial {
	return new(Theme)
}
func (ser Articles) All(ins gobzip.Serializer) {
	View.Articles = ins.(Articles)
}
func (ser Blogs) All(ins gobzip.Serializer) {
	View.Blogs = ins.(Blogs)
}
func (ser Globals) All(ins gobzip.Serializer) {
	View.Globals = ins.(Globals)
}
func (ser Resources) All(ins gobzip.Serializer) {
	View.Resources = ins.(Resources)
}
func (ser Rubrics) All(ins gobzip.Serializer) {
	View.Rubrics = ins.(Rubrics)
}
func (ser Servers) All(ins gobzip.Serializer) {
	View.Servers = ins.(Servers)
}
func (ser Themes) All(ins gobzip.Serializer) {
	View.Themes = ins.(Themes)
}
func (ser Articles) NewKey() int {
	id := 1
	for _, v := range ser {
		if v.ID > id {
			id = v.ID
		}
	}
	return id + 1
}
func (ser Blogs) NewKey() int {
	id := 1
	for _, v := range ser {
		if v.ID > id {
			id = v.ID
		}
	}
	return id + 1
}
func (ser Globals) NewKey() int {
	id := 1
	for _, v := range ser {
		if v.ID > id {
			id = v.ID
		}
	}
	return id + 1
}
func (ser Resources) NewKey() int {
	id := 1
	for _, v := range ser {
		if v.ID > id {
			id = v.ID
		}
	}
	return id + 1
}
func (ser Rubrics) NewKey() int {
	id := 1
	for _, v := range ser {
		if v.ID > id {
			id = v.ID
		}
	}
	return id + 1
}
func (ser Servers) NewKey() int {
	id := 1
	for _, v := range ser {
		if v.ID > id {
			id = v.ID
		}
	}
	return id + 1
}
func (ser Themes) NewKey() int {
	id := 1
	for _, v := range ser {
		if v.ID > id {
			id = v.ID
		}
	}
	return id + 1
}
func (ser Articles) At(key int) gobzip.Serial {
	for k, v := range ser {
		if v.ID == key {
			return ser[k]
		}
	}
	return nil
}
func (ser Blogs) At(key int) gobzip.Serial {
	for k, v := range ser {
		if v.ID == key {
			return ser[k]
		}
	}
	return nil
}
func (ser Globals) At(key int) gobzip.Serial {
	for k, v := range ser {
		if v.ID == key {
			return ser[k]
		}
	}
	return nil
}
func (ser Resources) At(key int) gobzip.Serial {
	for k, v := range ser {
		if v.ID == key {
			return ser[k]
		}
	}
	return nil
}
func (ser Rubrics) At(key int) gobzip.Serial {
	for k, v := range ser {
		if v.ID == key {
			return ser[k]
		}
	}
	return nil
}
func (ser Servers) At(key int) gobzip.Serial {
	for k, v := range ser {
		if v.ID == key {
			return ser[k]
		}
	}
	return nil
}
func (ser Themes) At(key int) gobzip.Serial {
	for k, v := range ser {
		if v.ID == key {
			return ser[k]
		}
	}
	return nil
}
func (ser Articles) Insert(s gobzip.Serial) gobzip.Serializer {
	ser = append(ser, s.(*Article))
	return ser
}
func (ser Blogs) Insert(s gobzip.Serial) gobzip.Serializer {
	ser = append(ser, s.(*Blog))
	return ser
}
func (ser Globals) Insert(s gobzip.Serial) gobzip.Serializer {
	ser = append(ser, s.(*Global))
	return ser
}
func (ser Resources) Insert(s gobzip.Serial) gobzip.Serializer {
	ser = append(ser, s.(*Resource))
	return ser
}
func (ser Rubrics) Insert(s gobzip.Serial) gobzip.Serializer {
	ser = append(ser, s.(*Rubric))
	return ser
}
func (ser Servers) Insert(s gobzip.Serial) gobzip.Serializer {
	ser = append(ser, s.(*Server))
	return ser
}
func (ser Themes) Insert(s gobzip.Serial) gobzip.Serializer {
	ser = append(ser, s.(*Theme))
	return ser
}
func (ser Articles) Replace(s gobzip.Serial) os.Error {
	for k, v := range ser {
		if v.ID == s.Key() {
			ser[k] = s.(*Article)
			return nil
		}
	}
	return os.ENOENT
}
func (ser Blogs) Replace(s gobzip.Serial) os.Error {
	for k, v := range ser {
		if v.ID == s.Key() {
			ser[k] = s.(*Blog)
			return nil
		}
	}
	return os.ENOENT
}
func (ser Globals) Replace(s gobzip.Serial) os.Error {
	for k, v := range ser {
		if v.ID == s.Key() {
			ser[k] = s.(*Global)
			return nil
		}
	}
	return os.ENOENT
}
func (ser Resources) Replace(s gobzip.Serial) os.Error {
	for k, v := range ser {
		if v.ID == s.Key() {
			ser[k] = s.(*Resource)
			return nil
		}
	}
	return os.ENOENT
}
func (ser Rubrics) Replace(s gobzip.Serial) os.Error {
	for k, v := range ser {
		if v.ID == s.Key() {
			ser[k] = s.(*Rubric)
			return nil
		}
	}
	return os.ENOENT
}
func (ser Servers) Replace(s gobzip.Serial) os.Error {
	for k, v := range ser {
		if v.ID == s.Key() {
			ser[k] = s.(*Server)
			return nil
		}
	}
	return os.ENOENT
}
func (ser Themes) Replace(s gobzip.Serial) os.Error {
	for k, v := range ser {
		if v.ID == s.Key() {
			ser[k] = s.(*Theme)
			return nil
		}
	}
	return os.ENOENT
}
func (ser Articles) Init() gobzip.Serializer {
	s := make([]*Article, 0)
	var o Articles = s
	return o
}
func (ser Blogs) Init() gobzip.Serializer {
	s := make([]*Blog, 0)
	var o Blogs = s
	return o
}
func (ser Globals) Init() gobzip.Serializer {
	s := make([]*Global, 0)
	var o Globals = s
	return o
}
func (ser Resources) Init() gobzip.Serializer {
	s := make([]*Resource, 0)
	var o Resources = s
	return o
}
func (ser Rubrics) Init() gobzip.Serializer {
	s := make([]*Rubric, 0)
	var o Rubrics = s
	return o
}
func (ser Servers) Init() gobzip.Serializer {
	s := make([]*Server, 0)
	var o Servers = s
	return o
}
func (ser Themes) Init() gobzip.Serializer {
	s := make([]*Theme, 0)
	var o Themes = s
	return o
}
func (ser Articles) Keys() []int {
	keys := make([]int, 0)
	for _, v := range ser {
		keys = append(keys, v.ID)
	}
	return keys
}
func (ser Blogs) Keys() []int {
	keys := make([]int, 0)
	for _, v := range ser {
		keys = append(keys, v.ID)
	}
	return keys
}
func (ser Globals) Keys() []int {
	keys := make([]int, 0)
	for _, v := range ser {
		keys = append(keys, v.ID)
	}
	return keys
}
func (ser Resources) Keys() []int {
	keys := make([]int, 0)
	for _, v := range ser {
		keys = append(keys, v.ID)
	}
	return keys
}
func (ser Rubrics) Keys() []int {
	keys := make([]int, 0)
	for _, v := range ser {
		keys = append(keys, v.ID)
	}
	return keys
}
func (ser Servers) Keys() []int {
	keys := make([]int, 0)
	for _, v := range ser {
		keys = append(keys, v.ID)
	}
	return keys
}
func (ser Themes) Keys() []int {
	keys := make([]int, 0)
	for _, v := range ser {
		keys = append(keys, v.ID)
	}
	return keys
}

func (ser Articles) NewFromForm(from map[string][]string) gobzip.Serial {
	a := new(Article)
	key := View.KeyFromForm(from)
	if key == 0 {
		a.ID = ser.NewKey()
	} else {
		a.ID = key
	}
	a.Blog, _ = strconv.Atoi(from["Blog"][0])
	a.Date = time.LocalTime().Format("02.01.2006 15:04:05")
	a.Description = from["Description"][0]
	a.Keywords = from["Keywords"][0]
	a.Rubric, _ = strconv.Atoi(from["Rubric"][0])
	a.Teaser = from["Teaser"][0]
	a.Text = from["Text"][0]
	a.Title = from["Title"][0]
	a.Url = from["Url"][0]
	return a
}
func (ser Blogs) NewFromForm(from map[string][]string) gobzip.Serial {
	a := new(Blog)
	key := View.KeyFromForm(from)
	if key == 0 {
		a.ID = ser.NewKey()
	} else {
		a.ID = key
	}
	a.Description = from["Description"][0]
	a.Keywords = from["Keywords"][0]
	a.Server, _ = strconv.Atoi(from["Server"][0])
	a.Slogan = from["Slogan"][0]
	a.Template, _ = strconv.Atoi(from["Template"][0])
	a.Title = from["Title"][0]
	a.Url = from["Url"][0]
	return a
}
func (ser Rubrics) NewFromForm(from map[string][]string) gobzip.Serial {
	a := new(Rubric)
	key := View.KeyFromForm(from)
	if key == 0 {
		a.ID = ser.NewKey()
	} else {
		a.ID = key
	}
	a.Blog, _ = strconv.Atoi(from["Blog"][0])
	a.Description = from["Description"][0]
	a.Keywords = from["Keywords"][0]
	a.Title = from["Title"][0]
	a.Url = from["Url"][0]
	return a
}
func (ser Globals) NewFromForm(from map[string][]string) gobzip.Serial {
	return nil
}
func (ser Resources) NewFromForm(from map[string][]string) gobzip.Serial {
	return nil
}
func (ser Themes) NewFromForm(from map[string][]string) gobzip.Serial {
	return nil
}
func (ser Servers) NewFromForm(from map[string][]string) gobzip.Serial {
	fmt.Println(from)
	a := new(Server)
	key := View.KeyFromForm(from)
	if key == 0 {
		a.ID = ser.NewKey()
	} else {
		a.ID = key
	}
	a.IP = from["IP"][0]
	a.Vendor = from["Vendor"][0]
	return a
}

func (send *Article) Host() string {
	for _, v := range View.Servers {
		if send.getBlog().Server == v.ID {
			return v.IP
		}
	}
	return ""
}
func (send *Blog) Host() string {
	for _, v := range View.Servers {
		if send.Server == v.ID {
			return v.IP
		}
	}
	return ""
}
func (send *Rubric) Host() string {
	for _, v := range View.Servers {
		if send.getBlog().Server == v.ID {
			return v.IP
		}
	}
	return ""
}
func (send *Global) Host() string {
	return ""
}
func (send *Theme) Host() string {
	return ""
}
func (send *Resource) Host() string {
	return ""
}
/*func (multisend *Global) Hosts() []string {
	slice := make([]string, 0)
	for _, v := range View.Servers {
		slice = append(slice, v.IP)
	}
	return slice
}
func (multisend *Resource) Hosts() []string {
	slice := make([]string, 0)
	for _, v := range View.Servers {
		slice = append(slice, v.IP)
	}
	return slice
}
func (multisend *Theme) Hosts() []string {
	slice := make([]string, 0)
	for _, v := range View.Servers {
		slice = append(slice, v.IP)
	}
	return slice
}*/

func (p *Page) Delegate(kind string) gobzip.Serializer {
	switch kind {
	case "articles":
		return p.Articles
	case "blogs":
		return p.Blogs
	case "globals":
		return p.Globals
	case "resources":
		return p.Resources
	case "rubrics":
		return p.Rubrics
	case "servers":
		return p.Servers
	case "themes":
		return p.Themes
	}
	return nil
}
func (p *Page) KeyFromForm(from map[string][]string) int {
	if from["ID"] != nil {
		key, err := strconv.Atoi(from["ID"][0])
		if err != nil {
			return 0
		}
		return key
	}
	return 0
}

func (p *Page) Hosts() []string {
	slice := make([]string, 0)
	for _, v := range p.Servers {
		slice = append(slice, v.IP)
	}
	return slice
}
func (p *Page) Senders(kind string) gobzip.Serializer {
	switch kind {
	case "articles":
		return p.Articles
	case "blogs":
		return p.Blogs
	case "globals":
		return p.Globals
	case "resources":
		return p.Resources
	case "rubrics":
		return p.Rubrics
	case "servers":
		return p.Servers
	case "themes":
		return p.Themes
	}
	return nil
}
func (p *Page) Console()string {
	return master.Logged()
}
type Server struct {
	ID     int
	IP     string
	Vendor string
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
	Url         string
	Keywords    string
	Description string
	Blog        int
}

type Article struct {
	ID          int
	Date        string
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
	ID       int
	Name     string
	Template int
	Data     []byte
}

type Global struct {
	ID   int
	Name string
	Data []byte
}

type Page struct {
	HeadMeta  string
	Rubrics   Rubrics
	Articles  Articles
	Blogs     Blogs
	Themes    Themes
	Resources Resources
	Globals   Globals
	Servers   Servers
	Index     int
	Blog      int
	Rubric    int
	Article   int
	Server    int
	Theme     int
	Global    int
	Resource	int
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

func LoadAll() os.Error {
	err := master.LoadKind(&View.Articles)
	if err != nil {
		return err
	}
	err = master.LoadKind(&View.Blogs)
	if err != nil {
		return err
	}
	err = master.LoadKind(&View.Globals)
	if err != nil {
		return err
	}
	err = master.LoadKind(&View.Resources)
	if err != nil {
		return err
	}
	err = master.LoadKind(&View.Rubrics)
	if err != nil {
		return err
	}
	err = master.LoadKind(&View.Servers)
	if err != nil {
		return err
	}
	err = master.LoadKind(&View.Themes)
	if err != nil {
		return err
	}
	return nil
}
