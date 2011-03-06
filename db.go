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
var statement = new(Statement)

type Statement struct {
	Blogs     *mysql.Statement
	Themes    *mysql.Statement
	Resources *mysql.Statement
	Rubrics   *mysql.Statement
	Articles  *mysql.Statement
	Servers   *mysql.Statement
	Globals   *mysql.Statement
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
	statement.Servers, err = db.Prepare("SELECT * FROM servers")
	if err != nil {
		return err
	}
	statement.Blogs, err = db.Prepare("SELECT * FROM blogs WHERE Server=?")
	if err != nil {
		return err
	}
	statement.Themes, err = db.Prepare("SELECT * FROM templates")
	if err != nil {
		return err
	}
	statement.Resources, err = db.Prepare("SELECT * FROM resources WHERE Template !=-1")
	if err != nil {
		return err
	}
	statement.Globals, err = db.Prepare("SELECT * FROM resources WHERE Template=-1")
	if err != nil {
		return err
	}
	statement.Rubrics, err = db.Prepare("SELECT * FROM rubrics WHERE Blog IN (SELECT ID FROM blogs WHERE Server=?)")
	if err != nil {
		return err
	}
	statement.Articles, err = db.Prepare("SELECT * FROM articles WHERE Blog IN (SELECT ID FROM blogs WHERE Server=?) ORDER BY Date DESC")
	if err != nil {
		return err
	}
	return nil
}

func (p *Page) loadBlogData() os.Error {
	err := InitMysql(app.DataHost, app.User, app.Password, app.Database)
	if err != nil {
		return err
	}
	err = prepareMysql()
	if err != nil {
		return err
	}
	db.Lock()
	// Load Blogs
	err = statement.Blogs.BindParams(app.Server)
	if err != nil {
		return err
	}
	err = statement.Blogs.Execute()
	if err != nil {
		return err
	}
	err = statement.Blogs.StoreResult()
	if err != nil {
		return err
	}
	p.Blogs = make([]*Blog, statement.Blogs.RowCount())
	for k, _ := range p.Blogs {
		p.Blogs[k] = new(Blog)
		err = statement.Blogs.BindResult(&p.Blogs[k].ID, &p.Blogs[k].Title, &p.Blogs[k].Url,
			&p.Blogs[k].Template, &p.Blogs[k].Keywords, &p.Blogs[k].Description, &p.Blogs[k].Slogan,
			&p.Blogs[k].Server)
		if err != nil {
			return err
		}
		eof, err := statement.Blogs.Fetch()
		if err != nil {
			return err
		}
		if eof {
			return os.EOF
		}
	}

	// load Templates
	err = statement.Themes.Execute()
	if err != nil {
		return err
	}
	err = statement.Themes.StoreResult()
	if err != nil {
		return err
	}
	p.Themes = make([]*Theme, statement.Themes.RowCount())
	for k, _ := range p.Themes {
		p.Themes[k] = new(Theme)
		err = statement.Themes.BindResult(&p.Themes[k].ID, &p.Themes[k].Index, &p.Themes[k].Style,
			&p.Themes[k].Title, &p.Themes[k].FromUrl)
		if err != nil {
			return err
		}
		eof, err := statement.Themes.Fetch()
		if err != nil {
			return err
		}
		if eof {
			return os.EOF
		}
	}
	// load nonglobal resources
	err = statement.Resources.Execute()
	if err != nil {
		return err
	}
	err = statement.Resources.StoreResult()
	if err != nil {
		return err
	}
	p.Resources = make([]*Resource, statement.Resources.RowCount())
	for k, _ := range p.Resources {
		p.Resources[k] = new(Resource)
		err = statement.Resources.BindResult(&p.Resources[k].Name, &p.Resources[k].Template,
			&p.Resources[k].Data)
		if err != nil {
			return err
		}
		eof, err := statement.Resources.Fetch()
		if err != nil {
			return err
		}
		if eof {
			return os.EOF
		}
	}

	// load global resources
	err = statement.Globals.Execute()
	if err != nil {
		return err
	}
	err = statement.Globals.StoreResult()
	if err != nil {
		return err
	}
	p.Globals = make([]*Resource, statement.Globals.RowCount())
	for k, _ := range p.Globals {
		p.Globals[k] = new(Resource)
		err = statement.Globals.BindResult(&p.Globals[k].Name, &p.Globals[k].Template,
			&p.Globals[k].Data)
		if err != nil {
			return err
		}
		eof, err := statement.Globals.Fetch()
		if err != nil {
			return err
		}
		if eof {
			return os.EOF
		}
	}

	// load Servers
	err = statement.Servers.Execute()
	if err != nil {
		return err
	}
	err = statement.Servers.StoreResult()
	if err != nil {
		return err
	}
	p.Servers = make([]*Server, statement.Servers.RowCount())
	for k, _ := range p.Servers {
		p.Servers[k] = new(Server)
		err = statement.Servers.BindResult(&p.Servers[k].ID, &p.Servers[k].IP, &p.Servers[k].Vendor,
			&p.Servers[k].Type, &p.Servers[k].Item)
		if err != nil {
			return err
		}
		eof, err := statement.Servers.Fetch()
		if err != nil {
			return err
		}
		if eof {
			return os.EOF
		}
	}

	// Load Rubrics 
	err = statement.Rubrics.BindParams(app.Server)
	if err != nil {
		return err
	}
	err = statement.Rubrics.Execute()
	if err != nil {
		return err
	}
	err = statement.Rubrics.StoreResult()
	if err != nil {
		return err
	}
	p.Rubrics = make([]*Rubric, statement.Rubrics.RowCount())
	for k, _ := range p.Rubrics {
		p.Rubrics[k] = new(Rubric)

		err = statement.Rubrics.BindResult(&p.Rubrics[k].ID, &p.Rubrics[k].Title, &p.Rubrics[k].ShortUrl,
			&p.Rubrics[k].Keywords, &p.Rubrics[k].Description, &p.Rubrics[k].Blog)
		if err != nil {
			return err
		}
		eof, err := statement.Rubrics.Fetch()
		if err != nil {
			return err
		}
		if eof {
			return os.EOF
		}
	}

	// Load Articles 
	err = statement.Articles.BindParams(app.Server)
	if err != nil {
		return err
	}
	err = statement.Articles.Execute()
	if err != nil {
		return err
	}
	err = statement.Articles.StoreResult()
	if err != nil {
		return err
	}
	p.Articles = make([]*Article, statement.Articles.RowCount())
	for k, _ := range p.Articles {
		p.Articles[k] = new(Article)

		err = statement.Articles.BindResult(&p.Articles[k].ID, &p.Articles[k].Title, &p.Articles[k].Rubric, &p.Articles[k].Text,
			&p.Articles[k].Teaser, &p.Articles[k].Blog, &p.Articles[k].Keywords, &p.Articles[k].Description,
			&p.Articles[k].Date, &p.Articles[k].Url)
		if err != nil {
			return err
		}
		eof, err := statement.Articles.Fetch()
		if err != nil {
			return err
		}
		if eof {
			return os.EOF
		}
	}
	db.Unlock()

	statement.Blogs.FreeResult()
	statement.Articles.FreeResult()
	statement.Resources.FreeResult()
	statement.Globals.FreeResult()
	statement.Rubrics.FreeResult()
	statement.Themes.FreeResult()
	statement.Servers.FreeResult()
	return nil
}
