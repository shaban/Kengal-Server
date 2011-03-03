package main

import (
	"fmt"
	"mysql"
	"os"
)

func (e *BlogError) String() string {
	return fmt.Sprintf("<h1>%v</h1><p>%s</p>", e.Code, e.Msg)
}
// Those errors have to be rewritten once done debugging to make sense in 
// Server Error Page. Possibly need to make a function that wraps stuff in
// a few html tags like h1 p etc.

var DB_CONN_FAILED = &BlogError{506, "Mysql Connection failed. I called but noone at home."}
var DB_STMT_FAILED = &BlogError{507, "Mysql Statement failed. Better put Logging on and see what happened."}

var db *mysql.Client
var stmt = new(Stmt)

type Stmt struct {
	Blog      *mysql.Statement
	Template  *mysql.Statement
	Resources *mysql.Statement
	Global    *mysql.Statement
	Rubrics   *mysql.Statement
	Latest    *mysql.Statement
	Articles  *mysql.Statement
	MyArticle *mysql.Statement
}

func InitMysql(sock, user, pw, dbname string) os.Error {
	var err os.Error
	if sock == "" {
		db, err = mysql.DialUnix(mysql.DEFAULT_SOCKET, user, pw, dbname)
	} else {
		db, err = mysql.DialTCP(sock, user, pw, dbname)
	}
	if err != nil {
		return err
	}
	db.LogLevel = 0
	db.Query("SET NAMES 'utf8'")
	return nil
}
func prepareMysql() os.Error {
	var err os.Error
	stmt.Blog, err = db.Prepare("SELECT * FROM blogs WHERE Url=?")
	if err != nil {
		return err
	}
	stmt.Template, err = db.Prepare("SELECT * FROM templates WHERE ID=?")
	if err != nil {
		return err
	}
	stmt.Resources, err = db.Prepare("SELECT * FROM resources WHERE Template=?")
	if err != nil {
		return err
	}
	stmt.Global, err = db.Prepare("SELECT * FROM resources WHERE Template=?")
	if err != nil {
		return err
	}
	stmt.Rubrics, err = db.Prepare("SELECT * FROM rubrics WHERE ID IN (SELECT Distinct Rubric As ID FROM articles WHERE Blog=?)")
	if err != nil {
		return err
	}
	stmt.Latest, err = db.Prepare("SELECT * FROM articles WHERE Blog=? ORDER BY Date DESC LIMIT 0,5")
	if err != nil {
		return err
	}
	stmt.Articles, err = db.Prepare("SELECT * FROM articles WHERE Blog=? AND Rubric=? ORDER BY Date")
	if err != nil {
		return err
	}
	stmt.MyArticle, err = db.Prepare("SELECT * FROM articles where ID=?")
	if err != nil {
		return err
	}
	return nil
}

func (p *Page) loadBlogData(host string) os.Error {
	// set some variables that arent database driven manually
	db.Lock()

	// Load Data for the Blog which Url matches the Host as designated by Client Browser
	p.Domain = new(Blog)
	err := stmt.Blog.BindExec(host)
	if err != nil {
		return err
	}
	err = stmt.Blog.BindFetch(&p.Domain.ID, &p.Domain.Title, &p.Domain.Url, &p.Domain.Template, &p.Domain.Keywords,
	  &p.Domain.Description, &p.Domain.Slogan)
	if err != nil {
		return err
	}

	// load Html Templates for this Blog
	p.Template = new(Theme)
	err = stmt.Template.BindExec(p.Domain.Template)
	if err != nil {
		return err
	}
	err = stmt.Template.BindFetch(&p.Template.ID, &p.Template.Index, &p.Template.Style)
	if err != nil {
		return err
	}
	
	// load JS,CSS,Images for this Template
	err = stmt.Resources.BindExec(p.Domain.Template)
	if err != nil {
		return err
	}
	p.Template.Resources = make([]*Resource, stmt.Resources.RowCount())
	for k, _ := range p.Template.Resources {
		p.Template.Resources[k] = new(Resource)
		err = stmt.Resources.BindFetch(&p.Template.Resources[k].Name, &p.Template.Resources[k].Template, 
		&p.Template.Resources[k].Data)
		if err != nil {
			return err
		}
	}

	// loadGlobal JS,Images.Css etc. -1 stands for Resources 
	// which are not connected to a specific Template
	err = stmt.Global.BindExec(-1)
	if err != nil {
		return err
	}
	p.Template.Global = make([]*Resource, stmt.Global.RowCount())
	for k, _ := range p.Template.Global {
		p.Template.Global[k] = new(Resource)
		err = stmt.Global.BindFetch(&p.Template.Global[k].Name, &p.Template.Global[k].Template, 
		&p.Template.Global[k].Data)
		if err != nil {
			return err
		}
	}

	// Load all distinct Rubrics that occur in this Blog 
	err = stmt.Rubrics.BindExec(p.Domain.ID)
	if err != nil {
		return err
	}
	p.Rubrics = make([]*Rubric, stmt.Rubrics.RowCount())
	for k, _ := range p.Rubrics {
		p.Rubrics[k] = new(Rubric)

		err = stmt.Rubrics.BindFetch(&p.Rubrics[k].ID, &p.Rubrics[k].Title, &p.Rubrics[k].ShortUrl, 
		&p.Rubrics[k].Keywords, &p.Rubrics[k].Description)
		if err != nil {
			return err
		}
	}

	// Load Latest 5 Articles that occur in this Blog 
	err = stmt.Latest.BindExec(p.Domain.ID)
	if err != nil {
		return err
	}
	p.Latest = make([]*Article, stmt.Latest.RowCount())
	for k, _ := range p.Latest {
		p.Latest[k] = new(Article)

		err = stmt.Latest.BindFetch(&p.Latest[k].ID, &p.Latest[k].Title, &p.Latest[k].Rubric, &p.Latest[k].Text, 
		&p.Latest[k].Teaser, &p.Latest[k].Blog, &p.Latest[k].Keywords, &p.Latest[k].Description, 
		&p.Latest[k].Date, &p.Latest[k].Url)
		if err != nil {
			return err
		}
	}
	db.Unlock()
	return nil
}
func (p *Page) loadArticlesInRubric(id string) os.Error {
	// Load all Articles that are in the rubric
	db.Lock()
	err := stmt.Articles.BindExec(p.Domain.ID, id)
	if err != nil {
		return err
	}
	p.Articles = make([]*Article, stmt.Articles.RowCount())
	for k, _ := range p.Articles {
		p.Articles[k] = new(Article)
		err = stmt.Articles.BindFetch(&p.Articles[k].ID, &p.Articles[k].Title, &p.Articles[k].Rubric, 
		&p.Articles[k].Text, &p.Articles[k].Teaser, &p.Articles[k].Blog, &p.Articles[k].Keywords, 
		&p.Articles[k].Description, &p.Articles[k].Date, &p.Articles[k].Url)
		if err != nil {
			return err
		}
	}
	db.Unlock()
	return nil
}
func (p *Page) loadMyArticle(id string) os.Error {
	db.Lock()
	// Load the full selected Article
	err := stmt.MyArticle.BindExec(id)
	if err != nil {
		return err
	}
	p.MyArticle = new(Article)
	err = stmt.MyArticle.BindFetch(&p.MyArticle.ID, &p.MyArticle.Title, &p.MyArticle.Rubric, &p.MyArticle.Text, 
	&p.MyArticle.Teaser, &p.MyArticle.Blog, &p.MyArticle.Keywords, &p.MyArticle.Description, &p.MyArticle.Date, 
	&p.MyArticle.Url)
	if err != nil {
		return err
	}
	db.Unlock()
	return nil
}
