package gobzip

import (
	"os"
	"fmt"
	"bytes"
	"strings"
	"gob"
	"io"
	"compress/gzip"
	"http"
	"path"
)
// Datainterfaces
type Delegator interface{
	Delegate(kind string)Serializer
	KeyFromForm(from map[string][]string)int
}
type Sender interface{
	Host()string
}
type MultiSender interface{
	Hosts()[]string
}
type Serial interface{
	Key()int
	Kind()string
}
type SerialSender interface{
	Serial
	Sender
}
type SerialMultiSender interface{
	Serial
	MultiSender
}
type Serializer interface {
	New() Serial
	Init()Serializer
	Insert(insert Serial) Serializer
	Replace (elem Serial)os.Error
	Kind ()string
	All(ser Serializer)
	NewFromForm(from map[string][]string)Serial
	At(int)Serial
	Keys()[]int
}
// Databaseinterfaces
type FileSystem struct{
	root string
	Delegator Delegator
	deletePattern string
	replacePattern string
	insertPattern string
	auditPattern string
}
type MasterFileSystem struct{
	FileSystem
}
type ClientFileSystem struct{
	FileSystem
	//MasterAudit string
}

type Database interface{
	Delete(s Serial)os.Error
	Save(s Serial)os.Error
	Init(Root, DeletePattern, ReplacePattern, InsertPattern,AuditPattern string)
}
type MasterDatabase interface{
	Database
	LoadKind(ser Serializer)os.Error
	//HandleForms(DeletePattern, ReplacePattern, InsertPattern,AuditPattern string)
}
type ClientDatabase interface{
	Database
	SaveKind(ser Serializer)os.Error
}

func (db *FileSystem)Delete(s Serial)os.Error{
	err := os.Remove(fmt.Sprintf("%s/%s/%v.bin.gz", db.root,s.Kind(), s.Key()))
	return err
}
var DefaultMaster *MasterFileSystem = new(MasterFileSystem)
var DefaultClient *ClientFileSystem = new(ClientFileSystem)

func (db *FileSystem)Save(s Serial)os.Error{
	f, err := os.Open(fmt.Sprintf("%s/%s/%v.bin.gz", db.root,s.Kind(), s.Key()), os.O_CREATE|os.O_WRONLY, 0666)
	defer f.Close()
	if err != nil {
		return err
	}
	buf := bytes.NewBufferString("")
	
	genc := gob.NewEncoder(buf)
	err = genc.Encode(s)
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
func (db *FileSystem)Init(Root, DeletePattern, ReplacePattern, InsertPattern,AuditPattern string){
	db.root = Root
	db.deletePattern = DeletePattern
	db.insertPattern = InsertPattern
	db.replacePattern = ReplacePattern
	db.auditPattern = AuditPattern
}
func (db *MasterFileSystem)HandleForms(){
	http.HandleFunc(db.deletePattern, handleDeleteForm)
	http.HandleFunc(db.replacePattern, handleReplaceForm)
	http.HandleFunc(db.insertPattern, handleInsertForm)
	http.HandleFunc(db.auditPattern, handleAudit)
}
func handleDeleteForm(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	_, kind := path.Split(r.URL.Path)
	key := DefaultMaster.Delegator.KeyFromForm(r.Form)
	ser := DefaultMaster.Delegator.Delegate(kind)
	keys := ser.Keys()
	n := ser.Init()
	for _,v := range keys{
	fmt.Println(v)
		if v != key{
			n = n.Insert(ser.At(v))
		}
	}
	ser.All(n)
	s := ser.At(key)
	DefaultMaster.Delete(s)
}
func handleReplaceForm(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	_, kind := path.Split(r.URL.Path)
	ser := DefaultMaster.Delegator.Delegate(kind)
	s := ser.NewFromForm(r.Form)
	ser.Replace(s)
	DefaultMaster.Save(s)
}
func handleInsertForm(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	_, kind := path.Split(r.URL.Path)
	ser := DefaultMaster.Delegator.Delegate(kind)
	s := ser.NewFromForm(r.Form)
	ser.All(ser.Insert(s))
	DefaultMaster.Save(s)
}
func handleAudit(w http.ResponseWriter, r *http.Request){
	
}
func (db *MasterFileSystem)LoadKind(ser Serializer)os.Error{
	fdir, err := os.Open(fmt.Sprintf("%s/%s", db.root,ser.Kind()),os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer fdir.Close()
	dir, err := fdir.Readdirnames(-1)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	for _, fileName := range dir {
		if strings.HasSuffix(fileName, "~") || strings.HasPrefix(fileName, ".") {
			continue
		}
		
		f, err := os.Open(fmt.Sprintf("%s/%s/%s", db.root,ser.Kind(), fileName),os.O_RDONLY, 0)
		defer f.Close()
		if err != nil {
			return  err
		}
		gz, err := gzip.NewReader(f)
		defer gz.Close()
		
		item := ser.New()
		gdec := gob.NewDecoder(gz)
		err = gdec.Decode(item)
		if err != nil {
			return err
		}
		ser.All(ser.Insert(item))
		fmt.Println(ser.Kind())
	}
	return nil
}
func (db *ClientFileSystem)SaveKind(ser Serializer)os.Error{
	r, _, err := http.Get(db.auditPattern)
	if err!= nil{
		return err
	}
	gz, err := gzip.NewReader(r.Body)
	if err != nil {
		return err
	}
	defer gz.Close()
	defer r.Body.Close()
	decoder := gob.NewDecoder(gz)
	err = decoder.Decode(ser)
	if err != nil {
		return err
	}
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