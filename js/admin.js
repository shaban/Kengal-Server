$.ajaxSetup( {
	cache : false
});
function retrieveTextAndSubmit() {
	document.getElementById("editor1").value = CKEDITOR.instances.editor1
			.getData();
	document.getElementById("editor2").value = CKEDITOR.instances.editor2
			.getData();

	$.post('/admin/article/save', $('#UpdateEditor').serialize(),
			function(data) {
				$('#console').prepend(data + '<hr/>');
			});
}

function spawnEditors() {
	CKEDITOR.config.enterMode =  CKEDITOR.ENTER_P;
	CKEDITOR.config.forcePasteAsPlainText = true;
	if (CKEDITOR.instances['editor1'])
		delete CKEDITOR.instances['editor1'];
	CKEDITOR.replace('editor1', {
		toolbar : [[ 'Source', '-', 'Undo', 'Redo', '-', 'Find', 'Replace', '-','Bold', 'Italic', 'Underline', 'Strike', '-','Subscript', 'Superscript' ],
				[ 'NumberedList', 'BulletedList', '-', 'Styles','Blockquote' ],
				[ 'JustifyLeft', 'JustifyCenter', 'JustifyRight','JustifyBlock' ], 
				[ 'Link', 'Unlink','Image', 'Table' ] ],
		height : 130
	});

	if (CKEDITOR.instances['editor2'])
		delete CKEDITOR.instances['editor2'];
	CKEDITOR.replace('editor2', {
		toolbar : [[ 'Source', '-', 'Undo', 'Redo', '-', 'Find', 'Replace', '-','Bold', 'Italic', 'Underline', 'Strike', '-','Subscript', 'Superscript' ],
				[ 'NumberedList', 'BulletedList', '-', 'Outdent', 'Indent',	'Blockquote' ],[ 'JustifyLeft', 'JustifyCenter', 'JustifyRight','JustifyBlock' ],
				[ 'Link', 'Unlink' ,'Image', 'Table' ]],
		height : 270
	});
}

function getMarkupEditor() {
	var editor = CodeMirror.fromTextArea('htmlcode', {
		height : "180px",
		parserfile : "parsexml.js",
		stylesheet : "css/xmlcolors.css",
		path : "js/",
		continuousScanning : 500,
		lineNumbers : true
	});

	var editor2 = CodeMirror.fromTextArea('csscode', {
		height : "250px",
		parserfile : "parsecss.js",
		stylesheet : "css/csscolors.css",
		path : "js/",
		continuousScanning : 500,
		lineNumbers : true
	});
}
function saveNewElem(action, elem) {
	$.post(action, $(elem).serialize(), function(data) {
		$('#console').prepend(data + '<hr/>');
		$("#detail1").empty();
		$("#detail2").empty();
		$("#detail3").empty();
		$("#detail4").empty();
	});
}
function delElem(action){
	$.get(action, function(data) {
		$('#console').prepend(data + '<hr/>');
		$("#detail1").empty();
		$("#detail2").empty();
		$("#detail3").empty();
		$("#detail4").empty();
	});
}
function loadDetail(pane, url, reset) {

	if (pane == 1) {
		$("#detail1").load(url);
		if (reset) {
			$("#detail2").empty();
			$("#detail3").empty();
			$("#detail4").empty();
		} else {
			$("#detail2").empty();
			$("#detail3").empty();
			$("#detail4").empty();
		}
		return;
	}
	if (pane == 2) {
		$("#detail2").load(url);
		if (reset) {
			$("#detail3").empty();
			$("#detail4").empty();
		} else {
			$("#detail3").empty();
			$("#detail4").empty();
			$("#detailBlog").toggle(false);
		}
		return;
	}
	if (pane == 3) {
		$("#detail3").load(url);
		if (reset) {
			$("#detail4").empty();
		} else {
			$("#detail4").empty();
			$("#detailBlog").toggle(false);
			$("#detailRubric").toggle(false);
		}
		return;
	}
	if (pane == 4) {
		$("#detail4").load(url);
		if (reset === undefined) {
			$("#detailBlog").toggle(false);
			$("#detailRubric").toggle(false);
			$("#detailArticle").toggle(false);
		}
		return;
	}
}