package main

import (
"http"
"fmt"
)
const ErrorNotFound = `
	<html>
	<head>
	<style>
	body{
		font-family: Arial;
	}
	div{
		padding: 15px 30px;
		border: 1px solid black;
		width: 640px;
		margin-top: 40px;
		margin-left: auto;
		margin-right: auto;
		background: white;
		color: #369;
	}
	a {
		font-weight: bold;
		color: #321;
	}
	a:hover{
		color: #963;
	}
	</style>
	</head><body>
		<div>
			<h1>404 - Not Found</h1>
			<p>Möglicherweise haben Sie die falsche Url eingegeben.</p>
			<br/>Zurück zur
			<a href="/">Startseite</a>
			<p>Kengal 0.9.1</p>
		</div>
	</body>
	</html>`

const ErrorForbidden = `
	<html>
	<head>
	<style>
	body{
		font-family: Arial;
	}
	div{
		padding: 15px 30px;
		border: 1px solid black;
		width: 640px;
		margin-top: 40px;
		margin-left: auto;
		margin-right: auto;
		background: white;
		color: #369;
	}
	a {
		font-weight: bold;
		color: #321;
	}
	a:hover{
		color: #963;
	}
	</style>
	</head><body>
		<div>
			<h1>403 - Forbidden</h1>
			<p>Zugriff nicht erlaubt</p>
			<p>Kengal 0.9.1</p>
		</div>
	</body>
	</html>`
const errHtml = `
	<html>
	<head>
	<style>
	body{
		font-family: Arial;
	}
	div{
		padding: 15px 30px;
		border: 1px solid black;
		width: 640px;
		margin-top: 40px;
		margin-left: auto;
		margin-right: auto;
		background: white;
		color: #369;
	}
	a {
		font-weight: bold;
		color: #321;
	}
	a:hover{
		color: #963;
	}
	</style>
	</head><body>
		<div>
			<h1>%v - %s</h1>
			<p><a href="http://webmaster.indiana.edu/tool_guide_info/errorcodes.shtml" target="_blank">%s</a></p>
			<p>Kengal 0.9.1</p>
		</div>
	</body>
	</html>`

type KengalWebError struct{
	Code int 
	Msg string
}

func (kw *KengalWebError)getErrorName()string{
	if kw.Code == 404{
		return "Not Found"
	}
	if kw.Code == 403{
		return "Forbidden"
	}
	return ""
}

func (kw *KengalWebError)Write(w http.ResponseWriter){
	w.WriteHeader(kw.Code)
	msg := fmt.Sprintf(errHtml,kw.Code,kw.getErrorName(),kw.Msg)
	w.Write([]byte(msg))
	w.Flush()
}