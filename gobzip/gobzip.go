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
	"log"
	"strconv"
)
// Datainterfaces
type logMsg struct { 
        data []byte 
        next *logMsg 
}
type reverseBuffer struct { 
        logMsgs *logMsg 
}
type reverseReader struct { 
        data []byte 
        logMsgs *logMsg 
}
func (w *reverseBuffer) Write(data []byte) (int, os.Error) { 
        if len(data) == 0 { 
                return 0, nil 
        } 
        w.logMsgs = &logMsg{append([]byte(nil), data...), w.logMsgs} 
        return len(data), nil 
}
func (w *reverseBuffer) Reader() io.Reader { 
        return &reverseReader{nil, w.logMsgs} 
}
func (r *reverseReader) Read(data []byte) (int, os.Error) { 
        if len(r.data) == 0 { 
                if r.logMsgs == nil { 
                        return 0, os.EOF 
                } 
                r.data = r.logMsgs.data 
                r.logMsgs = r.logMsgs.next 
        } 
        n := copy(data, r.data) 
        r.data = r.data[n:] 
        return n, nil 
} 
type Delegator interface {
	Delegate(kind string) Serializer
	KeyFromForm(from map[string][]string) int
}
type Sender interface {
	Host() string
}
type Serial interface {
	Key() int
	Kind() string
	Log()string
}
type SerialSender interface {
	Serial
	Sender
}
type Serializer interface {
	New() Serial
	Init() Serializer
	Insert(insert Serial) Serializer
	Replace(elem Serial) os.Error
	Kind() string
	All(ser Serializer)
	NewFromForm(from map[string][]string) Serial
	At(int) Serial
	Keys() []int
}
// Databaseinterfaces
type FileSystem struct {
	root           string
	delegator      Delegator
	deletePattern  string
	replacePattern string
	insertPattern  string
	auditPattern   string
	logData *reverseBuffer
	Logger *log.Logger
}
type MasterFileSystem struct {
	FileSystem
}
type ClientFileSystem struct {
	FileSystem
}

type Database interface {
	Delete(s Serial) os.Error
	Save(s Serial) os.Error
	Init(deleg Delegator, Root, DeletePattern, ReplacePattern, InsertPattern, AuditPattern string)
	Logged()string
	ClearLog()
}
type MasterDatabase interface {
	Database
	LoadKind(ser Serializer) os.Error
	HandleForm(pattern string, w http.ResponseWriter, r *http.Request)
}
type ClientDatabase interface {
	Database
	SaveKind(ser Serializer) os.Error
}

func (db *FileSystem) Delete(s Serial) os.Error {
	err := os.Remove(fmt.Sprintf("%s/%s/%v.bin.gz", db.root, s.Kind(), s.Key()))
	return err
}

var DefaultMaster *MasterFileSystem = new(MasterFileSystem)
var DefaultClient *ClientFileSystem = new(ClientFileSystem)

func (db *FileSystem) Save(s Serial) os.Error {
	f, err := os.Open(fmt.Sprintf("%s/%s/%v.bin.gz", db.root, s.Kind(), s.Key()), os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	err = db.save(f,s)
	return err
}
func (db *FileSystem) save(w io.Writer,s interface{}) os.Error {
	buf := bytes.NewBufferString("")

	genc := gob.NewEncoder(buf)
	err := genc.Encode(s)
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
func (db *FileSystem) SaveKind(w io.Writer,s Serializer) os.Error {
	return db.save(w,s)
}
func (db *FileSystem)Logged()string{
	out := bytes.NewBufferString("")
	io.Copy(out,db.logData.Reader())
	return out.String()
}
func (db *FileSystem)ClearLog(){
	db.logData = new(reverseBuffer)
	db.Logger = log.New(db.logData,"",log.Ldate | log.Ltime)
}

func (db *FileSystem) Init(deleg Delegator, Root, DeletePattern, ReplacePattern, InsertPattern, AuditPattern string) {
	db.delegator = deleg
	db.root = Root
	db.deletePattern = DeletePattern
	db.insertPattern = InsertPattern
	db.replacePattern = ReplacePattern
	db.auditPattern = AuditPattern
	db.logData = new(reverseBuffer)
	db.Logger = log.New(db.logData,"",log.Ldate | log.Ltime)
}
func (db *MasterFileSystem) HandleForms() {
	http.HandleFunc(db.deletePattern, handleDeleteForm)
	http.HandleFunc(db.replacePattern, handleReplaceForm)
	http.HandleFunc(db.insertPattern, handleInsertForm)
	http.HandleFunc(db.auditPattern, handleAudit)
}
func (db *MasterFileSystem) HandleForm(pattern string, w http.ResponseWriter, r *http.Request){
	dir,_ := path.Split(pattern)
	switch dir{
		case db.deletePattern:
			handleDeleteForm(w,r)
		case db.replacePattern:
			handleReplaceForm(w,r)
		case db.insertPattern:
			handleInsertForm(w,r)
	}
}
func handleDeleteForm(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	_, kind := path.Split(r.URL.Path)
	key := DefaultMaster.delegator.KeyFromForm(r.Form)
	ser := DefaultMaster.delegator.Delegate(kind)
	keys := ser.Keys()
	n := ser.Init()
	for _, v := range keys {
		if v != key {
			n = n.Insert(ser.At(v))
		}
	}
	ser.All(n)
	s := ser.At(key)
	DefaultMaster.Delete(s)
	DefaultMaster.Logger.Printf("%v erfolgreich gelöscht",s.Log())
	redir:= "http://"+r.Host+r.FormValue("Redir")
	w.SetHeader("Location",redir)
	w.WriteHeader(302)
}
func handleReplaceForm(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	_, kind := path.Split(r.URL.Path)
	ser := DefaultMaster.delegator.Delegate(kind)
	s := ser.NewFromForm(r.Form)
	ser.Replace(s)
	DefaultMaster.Save(s)
	DefaultMaster.Logger.Printf("%v erfolgreich modifiziert",s.Log())
	redir:= "http://"+r.Host+r.FormValue("Redir")
	w.SetHeader("Location",redir)
	w.WriteHeader(302)
}
func handleInsertForm(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	_, kind := path.Split(r.URL.Path)
	ser := DefaultMaster.delegator.Delegate(kind)
	s := ser.NewFromForm(r.Form)
	ser.All(ser.Insert(s))
	DefaultMaster.Save(s)
	DefaultMaster.Logger.Printf("%v erfolgreich angelegt",s.Log())
	redir:= "http://"+r.Host+r.FormValue("Redir")+strconv.Itoa(s.Key())
	w.SetHeader("Location",redir)
	w.WriteHeader(302)
}
func handleAudit(w http.ResponseWriter, r *http.Request) {
	_, kind := path.Split(r.URL.Path)
	ip := r.FormValue("IP")
	ser := DefaultMaster.delegator.Delegate(kind)
	keys := ser.Keys()
	n := ser.Init()
	for _, v := range keys {
		ss := ser.At(v).(SerialSender)
		host := ss.Host()
		if host == ip || host == ""{
			n=n.Insert(ss)
		}
	}
	if len(n.Keys()) != 0{
		w.SetHeader("Content-Encoding", "gzip")
		w.SetHeader("Content-Type", "application/octet-stream")
		DefaultMaster.SaveKind(w,n)
		return
	}
	w.WriteHeader(404)
}
func (db *MasterFileSystem) LoadKind(ser Serializer) os.Error {
	fdir, err := os.Open(fmt.Sprintf("%s/%s", db.root, ser.Kind()), os.O_RDONLY, 0)
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

		f, err := os.Open(fmt.Sprintf("%s/%s/%s", db.root, ser.Kind(), fileName), os.O_RDONLY, 0)
		defer f.Close()
		if err != nil {
			return err
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
	}
	return nil
}
func (db *ClientFileSystem) SaveKind(ser Serializer) os.Error {
	r, _, err := http.Get(db.auditPattern)
	if err != nil {
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
func (db *ClientFileSystem)Audit(MasterIP, ClientIP string, ser Serializer){
	r,_,err := http.Get(fmt.Sprintf("http://%s%s%s?IP=%s",MasterIP,db.auditPattern,ser.Kind(),ClientIP))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if r.StatusCode == 404{
		return
	}
	gz, err := gzip.NewReader(r.Body)
	defer gz.Close()
	defer r.Body.Close()
	gdec := gob.NewDecoder(gz)
	err = gdec.Decode(ser)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
func MakeAudit(w io.Writer, scheme interface{}) os.Error {
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