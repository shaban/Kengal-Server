package main

const Error = `
	<html>
	<head>
	<style>
	body{
		font-family: Arial;
	}
	div{
		padding: 15px;
		border: 1px solid black;
		width: 640px;
		margin-top: 40px;
		margin-left: auto;
		margin-right: auto;background: #aaa;
	}
	a {
		text-decoration: none;
		font-weight: bold;
		color: #333;
	}
	a:hover{
		text-decoration: underline;color: black;
	}
	</style>
	</head><body>
		<div>
			<h3>Kengal Webserver v.0.9</h3>
			<h1>Fehler: 404</h1>
			<p>Möglicherweise haben Sie die falsche Url eingegeben.</p>
			<br/>Zurück zur
			<a href="/">Startseite</a>
		</div>
	</body>
	</html>`
