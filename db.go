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

	UpdateBlog    *mysql.Statement
	UpdateRubric  *mysql.Statement
	UpdateArticle *mysql.Statement

	InsertBlog    *mysql.Statement
	InsertRubric  *mysql.Statement
	InsertArticle *mysql.Statement
	
	DeleteArticle	*mysql.Statement
	DeleteRubric	*mysql.Statement
}

func InitMysql() os.Error {
	var err os.Error
	db, err = mysql.DialUnix(mysql.DEFAULT_SOCKET, app.User, app.Password, app.Database)
	if err != nil {
		return err
	}
	db.LogLevel = uint8(app.LogLevel)
	db.Query("SET NAMES 'utf8'")
	return nil
}
func prepareMysql() os.Error {
	var err os.Error
	statement.UpdateBlog, err = db.Prepare("UPDATE blogs SET Description=?, Keywords=?, Server=?, Slogan=?, Template=?, Title=? ,Url=? WHERE ID=?")
	if err != nil {
		return err
	}
	
	statement.UpdateRubric, err = db.Prepare("UPDATE rubrics SET Description=?, Keywords=?, Blog=?, Title=? ,Url=? WHERE ID=?")
	if err != nil {
		return err
	}
	
	statement.UpdateArticle, err = db.Prepare("UPDATE articles SET Description=?, Keywords=?, Blog=?, Rubric=?, Date=?, Title=?, Url=?, Text=?, Teaser=? WHERE ID=?")
	if err != nil {
		return err
	}
	
	statement.InsertBlog, err = db.Prepare("INSERT INTO blogs (Description, Keywords, Server, Slogan, Template, Title, Url) VALUES (?,?,?,?,?,?,?)")
	if err != nil {
		return err
	}
	
	statement.InsertRubric, err = db.Prepare("INSERT INTO rubrics (Description, Keywords, Blog, Title, Url) VALUES (?,?,?,?,?)")
	if err != nil {
		return err
	}
	
	statement.InsertArticle, err = db.Prepare("INSERT INTO articles (Description, Keywords, Blog, Rubric, Title, Url, Text, Teaser,Date) VALUES (?,?,?,?,?,?,?,?,?)")
	if err != nil {
		return err
	}
	
	statement.DeleteRubric, err = db.Prepare("DELETE FROM rubrics where ID=?")
	if err != nil {
		return err
	}
	statement.DeleteArticle, err = db.Prepare("DELETE FROM articles where ID=?")
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
	statement.Themes, err = db.Prepare("SELECT * FROM templates")
	if err != nil {
		return err
	}
	statement.Servers, err = db.Prepare("SELECT * FROM servers")
	if err != nil {
		return err
	}
	statement.Blogs, err = db.Prepare("SELECT * FROM blogs ORDER by Title")
	if err != nil {
		return err
	}
	statement.Rubrics, err = db.Prepare("SELECT * FROM rubrics")
	if err != nil {
		return err
	}
	statement.Articles, err = db.Prepare("SELECT * FROM articles ORDER BY Date DESC")
	if err != nil {
		return err
	}
	return nil
}

func (p *Page) reloadBlogData() os.Error {
	p.Blogs = nil
	p.Articles = nil
	p.Rubrics = nil

	db.Lock()
	// Load Blogs
	/*err := statement.Blogs.BindParams(View.Server)
	if err != nil {
		return err
	}*/
	err := statement.Blogs.Execute()
	if err != nil {
		return err
	}
	err = statement.Blogs.StoreResult()
	if err != nil {
		return err
	}
	rc := statement.Blogs.RowCount()
	if rc == 0 {
		return os.ENOENT
	}
	p.Blogs = make([]*Blog, rc)
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

	// Load Rubrics 
	/*err = statement.Rubrics.BindParams(View.Server)
	if err != nil {
		return err
	}*/
	err = statement.Rubrics.Execute()
	if err != nil {
		return err
	}
	err = statement.Rubrics.StoreResult()
	if err != nil {
		return err
	}
	rc = statement.Rubrics.RowCount()
	if rc == 0 {
		return os.ENOENT
	}
	p.Rubrics = make([]*Rubric, rc)
	for k, _ := range p.Rubrics {
		p.Rubrics[k] = new(Rubric)

		err = statement.Rubrics.BindResult(&p.Rubrics[k].ID, &p.Rubrics[k].Title, &p.Rubrics[k].Url,
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
	/*err = statement.Articles.BindParams(View.Server)
	if err != nil {
		return err
	}*/
	err = statement.Articles.Execute()
	if err != nil {
		return err
	}
	err = statement.Articles.StoreResult()
	if err != nil {
		return err
	}
	rc = statement.Articles.RowCount()
	if rc == 0 {
		return os.ENOENT
	}
	p.Articles = make([]*Article, rc)
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
	statement.Rubrics.FreeResult()
	return nil
}

func (p *Page) loadBlogData() os.Error {
	db.Lock()
	// load Templates
	err := statement.Themes.Execute()
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

	// Load Blogs
	/*err = statement.Blogs.BindParams(View.Server)
	if err != nil {
		return err
	}*/
	err = statement.Blogs.Execute()
	if err != nil {
		return err
	}
	err = statement.Blogs.StoreResult()
	if err != nil {
		return err
	}
	rc := statement.Blogs.RowCount()
	if rc == 0 {
		return os.ENOENT
	}
	p.Blogs = make([]*Blog, rc)
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

	// Load Rubrics 
	/*err = statement.Rubrics.BindParams(View.Server)
	if err != nil {
		return err
	}*/
	err = statement.Rubrics.Execute()
	if err != nil {
		return err
	}
	err = statement.Rubrics.StoreResult()
	if err != nil {
		return err
	}
	rc = statement.Rubrics.RowCount()
	if rc == 0 {
		return os.ENOENT
	}
	p.Rubrics = make([]*Rubric, rc)
	for k, _ := range p.Rubrics {
		p.Rubrics[k] = new(Rubric)

		err = statement.Rubrics.BindResult(&p.Rubrics[k].ID, &p.Rubrics[k].Title, &p.Rubrics[k].Url,
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
	/*err = statement.Articles.BindParams(View.Server)
	if err != nil {
		return err
	}*/
	err = statement.Articles.Execute()
	if err != nil {
		return err
	}
	err = statement.Articles.StoreResult()
	if err != nil {
		return err
	}
	rc = statement.Articles.RowCount()
	if rc == 0 {
		return os.ENOENT
	}
	p.Articles = make([]*Article, rc)
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

func (b *Blog) updateBlog() {
	err := statement.UpdateBlog.BindParams(b.Description, b.Keywords, b.Server, b.Slogan, b.Template, b.Title, b.Url, b.ID)
	if err != nil {
		fmt.Println("pbl.BindParams")
		fmt.Println(err)
	}
	err = statement.UpdateBlog.Execute()
	if err != nil {
		fmt.Println("pbl.Execute")
		fmt.Println(err)
	}
	View.Blogs.Replace(b)
}

func updateRubric(b *Rubric) {
}

func updateArticle(b *Rubric) {
}

func insertBlog(b *Blog) {
}

func insertRubric(b *Rubric) {
}

func insertArticle(b *Rubric) {
}
