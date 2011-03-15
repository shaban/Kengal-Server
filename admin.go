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
	"mysql"
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
func popdb() {
	pbl, err := db.Prepare("Insert Into blogs (Description,Keywords,Server, Slogan, Template, Title,Url)Values(?,?,?,?,?,?,?)")
	if err != nil {
		fmt.Println(err)
	}
	prb, err := db.Prepare("Insert Into rubrics (Blog, Description, Keywords, Title, Url)Values(?,?,?,?,?)")
	if err != nil {
		fmt.Println(err)
	}
	par, err := db.Prepare("Insert INTO articles (Blog, Date, Description, Keywords, Rubric,Teaser,Text, Title,Url)Values(?,?,?,?,?,?,?,?,?)")
	if err != nil {
		fmt.Println(err)
	}

	for i := 0; i < 100; i++ {
		b := new(Blog)
		b.Description = fmt.Sprintf("Blog Description %v", i)
		b.Keywords = fmt.Sprintf("blog, key, words %v", i)
		b.Server = 2
		b.Slogan = fmt.Sprintf("Blog Slogan %v", i)
		b.Template = 2
		b.Title = fmt.Sprintf("Blog Title %v", i)
		b.Url = fmt.Sprintf("127.0.0.%v", i+2)

		err = pbl.BindParams(b.Description, b.Keywords, b.Server, b.Slogan, b.Template, b.Title, b.Url)
		if err != nil {
			fmt.Println("pbl.BindParams")
			fmt.Println(err)
		}
		err = pbl.Execute()
		if err != nil {
			fmt.Println("pbl.Execute")
			fmt.Println(err)
		}
		linsb := pbl.LastInsertId
		for n := 0; n < 10; n++ {
			r := new(Rubric)
			r.Blog = int(linsb)
			r.Description = fmt.Sprintf("Rubric Description %v", n)
			r.Keywords = fmt.Sprintf("rubric, key, words %v", n)
			r.Title = fmt.Sprintf("Rubric Title %v", n)
			r.Url = fmt.Sprintf("url%v", n)

			err = prb.BindParams(r.Blog, r.Description, r.Keywords, r.Title, r.Url)
			if err != nil {
				fmt.Println("prb.BindParams")
				fmt.Println(err)
			}
			err = prb.Execute()
			if err != nil {
				fmt.Println("prb.Execute")
				fmt.Println(err)
			}

			linsr := prb.LastInsertId

			for x := 0; x < 10; x++ {
				a := new(Article)

				dt := new(mysql.DateTime)
				dt.Day = 9
				dt.Year = 2011
				dt.Month = 3
				dt.Hour = 19
				dt.Minute = 38
				dt.Second = uint8(x)

				a.Blog = int(linsb)
				adate := fmt.Sprintf("2011-03-09 20:06:0%v", x)
				//a.Date = *dt
				a.Description = fmt.Sprintf("Article Description %v", x)
				a.Keywords = fmt.Sprintf("article, key, words %v", x)
				a.Rubric = int(linsr)
				a.Teaser = fmt.Sprintf("<p>aber hallo %v</p>", x)
				a.Text = "<p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p><p>TestText lang mit allem pipapo und kurzen übersichtlichen Passagen und Absätzen</p>"
				a.Title = fmt.Sprintf("Article Title %v", x)
				a.Url = fmt.Sprintf("url%v", x)

				fmt.Println(linsb)
				fmt.Println(linsr)
				err = par.BindParams(a.Blog, adate, a.Description, a.Keywords, a.Rubric, a.Teaser, a.Text, a.Title, a.Url)
				if err != nil {
					fmt.Println("par.BindParams")
					fmt.Println(err)
				}
				err = par.Execute()
				if err != nil {
					fmt.Println("par.Execute")
					fmt.Println(err)
				}
			}
		}
	}

}

func SnippetController(w http.ResponseWriter, r *http.Request) {
	j := r.FormValue("which")
	
	s, _ := strconv.Atoi(r.FormValue("server"))
	h := r.FormValue("host")
	a, _ := strconv.Atoi(r.FormValue("article"))
	rb, _ := strconv.Atoi(r.FormValue("rubric"))

	os := View.Server
	oh := View.Host
	oa := View.Article
	orb := View.Rubric

	View.Server = s
	View.Host = h
	View.Rubric = rb
	View.Article = a

	tplFile := ""
	w.SetHeader("Content-Type", "text/html; charset=utf-8")

	switch j {
	case "blogs":
		tplFile = "html/blogSnippet.html"
	case "articles":
		tplFile = "html/articleSnippet.html"
	case "rubrics":
		tplFile = "html/rubricSnippet.html"
	case "editor":
		tplFile = "html/editorSnippet.html"
	case "templates":
		data, _ := json.Marshal(View.Themes)
		w.Write(data)
		w.Flush()
	case "newblog":
		tplFile = "html/blogNewSnippet.html"
	case "newrubric":
		tplFile = "html/rubricNewSnippet.html"
	}
	var Templ = template.New(nil)
	snippetTempl, err := ioutil.ReadFile(tplFile)
	if err != nil {
		fmt.Println(err)
	}
	buf := bytes.NewBufferString("")
	err = Templ.Parse(string(snippetTempl))
	if err != nil {
		fmt.Println(err)
	}
	err = Templ.Execute(buf, View)
	if err != nil {
		fmt.Println(err)
	}
	w.Write(buf.Bytes())
	View.Server = os
	View.Host = oh
	View.Article = oa
	View.Rubric = orb
}

func AdminController(w http.ResponseWriter, r *http.Request) {
	//fmt.Println(r.FormValue("server"))
	i, err := strconv.Atoi(r.FormValue("server"))
	if err != nil {
		View.Server = app.Server
	} else {
		app.Server = View.Server
		View.Server = i
		if View.Server != app.Server {
			fmt.Println(app.Server)
			fmt.Println(View.Server)
			View.reloadBlogData()
			fmt.Println("success")
		}
	}
	//ParseParameters(r.URL.Path, r.Host)
	w.SetHeader("Content-Type", "text/html; charset=utf-8")
	w.SetHeader("Content-Encoding", "gzip")
	AdminDispatch(w)
}

