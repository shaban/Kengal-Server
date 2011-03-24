package main

import (
	"os"
	"fmt"
	"json"
	"io/ioutil"
	"strings"
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
	dir, err := ioutil.ReadDir(fmt.Sprintf("%s/%s", DB_ROOT,kind))
	if err != nil {
		return nil, err
	}
	for _, fileInfo := range dir {
		// ignore unix backup files (since data format is human readable one can edit like this)
		if strings.HasSuffix(fileInfo.Name, "~") || strings.HasPrefix(fileInfo.Name, ".") {
			continue
		}
		item := scheme.New()

		data, err := ioutil.ReadFile(fmt.Sprintf("%s/%s/%s", DB_ROOT,kind, fileInfo.Name))
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(data, &item)
		if err != nil {
			return nil, err
		}
		scheme = scheme.Insert(item).(Serializer)
	}
	return scheme, nil
}

func saveItem(kind string, item interface{}, key int) os.Error {
	data, err := json.MarshalIndent(item, "", "\t")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(fmt.Sprintf("%s/%s/%v.json", DB_ROOT,kind, key), data, 0666)
	if err != nil {
		return err
	}
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
func MakeAudit(scheme interface{})([]byte, os.Error){
	data, err := json.Marshal(scheme)
	return data, err
}