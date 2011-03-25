package main

import (
	"os"
	"fmt"
	"bytes"
	"strings"
	"gob"
	"io"
	"compress/gzip"
)

type Serializer interface {
	//Init()
	New() interface{}
	Insert(insert interface{}) interface{}
	//At(key string) interface{}
	//All() []interface{}
}
const DB_ROOT = "db"
func loadAll() os.Error {
	res, err := loadKind("articles", View.Articles)
	if err != nil {
		return err
	}
	View.Articles = res.(Articles)
	res, err = loadKind("blogs", View.Blogs)
	if err != nil {
		return err
	}
	View.Blogs = res.(Blogs)
	res, err = loadKind("globals", View.Globals)
	if err != nil {
		return err
	}
	View.Globals = res.(Globals)
	res, err = loadKind("resources", View.Resources)
	if err != nil {
		return err
	}
	View.Resources = res.(Resources)
	res, err = loadKind("rubrics", View.Rubrics)
	if err != nil {
		return err
	}
	View.Rubrics = res.(Rubrics)
	res, err = loadKind("servers", View.Servers)
	if err != nil {
		return err
	}
	View.Servers = res.(Servers)
	res, err = loadKind("themes", View.Themes)
	if err != nil {
		return err
	}
	View.Themes = res.(Themes)
	return nil
}
func loadKind(kind string, scheme Serializer) (interface{}, os.Error) {
	fdir, err := os.Open(fmt.Sprintf("%s/%s", DB_ROOT,kind),os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	defer fdir.Close()
	dir, err := fdir.Readdirnames(-1)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	for _, fileName := range dir {
		// ignore unix backup files and hidden files (since data format is human readable one can edit like this)
		if strings.HasSuffix(fileName, "~") || strings.HasPrefix(fileName, ".") {
			continue
		}
		
		f, err := os.Open(fmt.Sprintf("%s/%s/%s", DB_ROOT,kind, fileName),os.O_RDONLY, 0)
		defer f.Close()
		if err != nil {
			return nil, err
		}
		gz, err := gzip.NewReader(f)
		defer gz.Close()
		
		item := scheme.New()
		gdec := gob.NewDecoder(gz)
		err = gdec.Decode(item)
		if err != nil {
			return nil, err
		}
		scheme = scheme.Insert(item).(Serializer)
	}
	return scheme, nil
}

func saveItem(kind string, item interface{}, key int) os.Error {
	f, err := os.Open(fmt.Sprintf("%s/%s/%v.bin.gz", DB_ROOT,kind, key), os.O_CREATE|os.O_WRONLY, 0666)
	defer f.Close()
	if err != nil {
		return err
	}
	buf := bytes.NewBufferString("")
	
	genc := gob.NewEncoder(buf)
	err = genc.Encode(item)
	if err != nil {
		return err
	}
	gz, err := gzip.NewWriter(f)
	if err != nil {
		return err
	}
	_, err = gz.Write(buf.Bytes())
	if err != nil {
		return err
	}
	gz.Close()
	return nil
}
func deleteItem(kind string, id int) os.Error {
	err := os.Remove(fmt.Sprintf("%s/%s/%v.json", DB_ROOT,kind, id))
	return err
}

func updateBlog(b *Blog) os.Error {
	err := saveItem("blogs", b, b.ID)
	if err != nil {
		return err
	}
	View.Blogs.Replace(b)
	return nil
}

func updateRubric(rb *Rubric) os.Error {
	err := saveItem("rubrics", rb,rb.ID)
	if err != nil {
		return err
	}
	View.Rubrics.Replace(rb)
	return nil
}

func updateArticle(a *Article) os.Error {
	err := saveItem("articles", a, a.ID)
	if err != nil {
		return err
	}
	View.Articles.Replace(a)
	return nil
}

func insertBlog(b *Blog) os.Error {
	err := saveItem("blogs", b, b.ID)
	if err != nil {
		return err
	}
	View.Blogs = append(View.Blogs, b)
	return nil
}

func insertRubric(rb *Rubric) os.Error {
	err := saveItem("rubrics", rb, rb.ID)
	if err != nil {
		return err
	}
	View.Rubrics = append(View.Rubrics, rb)
	return nil
}

func insertArticle(a *Article) os.Error {
	err := saveItem("articles", a, a.ID)
	if err != nil {
		return err
	}
	View.Articles = append(View.Articles, a)
	return nil
}
func deleteRubric(id int) os.Error {
	s := make([]*Rubric, 0)
	for k, v := range View.Rubrics {
		if v.ID == id {
			err := deleteItem("rubrics", v.ID)
			if err != nil {
				return err
			}

		} else {
			s = append(s, View.Rubrics[k])
		}
	}
	View.Rubrics = s
	return nil
}

func deleteArticle(id int) os.Error {
	s := make([]*Article, 0)
	for k, v := range View.Articles {
		if v.ID == id {
			err := deleteItem("articles", v.ID)
			if err != nil {
				return err
			}
		} else {
			s = append(s, View.Articles[k])
		}
	}
	View.Articles = s
	return nil
}
func MakeAudit(w io.Writer, scheme interface{})(os.Error){
	buf := bytes.NewBufferString("")
	genc := gob.NewEncoder(buf)
	err := genc.Encode(scheme)
	if err != nil {
		return err
	}
	gz, err := gzip.NewWriter(w)
	if err != nil {
		return err
	}
	_, err = gz.Write(buf.Bytes())
	if err != nil {
		return err
	}
	gz.Close()
	return nil
}