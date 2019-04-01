
var isLoggedIn = false;
var jwtToken = "";
var xsrf_token = "" ;
var g_auth_key = "1234";

function renderError ( title, msg ) {
	var form = [ ''
		,'<div style="margin-bottom:24px;">'
			,'<div class="row">'
				,'<div class="col-sm-12">'
					,'<div class="card bg-danger">'
						,'<div class="card-header"><h4 style="color:white;">'+title+'</h4></div>'
						,'<div class="card-body bg-light">'
							,'<div>'+msg+'</div>'
						,'</div>'
					,'</div>'
				,'</div>'
			,'</div>'
		,'</div>'
	].join("\n");
	$("#msg").html(form);
}
function renderMessage ( title, msg ) {
	var form = [ ''
		,'<div style="margin-bottom:24px;">'
			,'<div class="row">'
				,'<div class="col-sm-12">'
					,'<div class="card bg-success">'
						,'<div class="card-header"><h4 style="color:white;">'+title+'</h4></div>'
						,'<div class="card-body bg-light">'
							,'<div>'+msg+'</div>'
						,'</div>'
					,'</div>'
				,'</div>'
			,'</div>'
		,'</div>'
	].join("\n");
	$("#msg").html(form);
}
function renderClearMessage ( ) {
	$("#msg").html("");
}
function render5SecClearMessage ( ) {
	setTimeout(function() {
		$("#msg").html("");
	},5000);
}

function submitItData ( event, data, action, succ, erro ) {
	if ( event ) {
		event.preventDefault();
	}

	console.log ( "form data (passed): ", data);
	$.ajax({
		type: 'GET',
		url: action,
		data: data,
		success: function (data) {
			console.log ( "success AJAX", data );
			$("#output").text( JSON.stringify(data, null, 4) );
			if ( succ ) {
				succ(data);
			}
		},
		error: function(resp) {
			$("#output").text( "Error!"+JSON.stringify(resp) );
			alert("got error status="+resp.status+" "+resp.statusText);
			if ( erro ) {
				erro(resp);
			}
		}
	});
}

function fetchRefData ( url, succ, erro ) {
	var data = { "auth_key": g_auth_key, "_ran_": ( Math.random() * 10000000 ) % 10000000 };
	$.ajax({
		type: 'GET',
		url: url,
		data: data,
		success: function (data) {
			$("#output").text( JSON.stringify(data, null, 4) );
			if ( succ ) {
				succ(data.data);
			}
		},
		error: function(resp) {
			$("#output").text( "Error!"+JSON.stringify(resp) );
			// alert("got error status="+resp.status+" "+resp.statusText);
			if ( erro ) {
				erro(resp);
			}
		}
	});
}

var curDocID = "";

var LoggInDone;	// call function on successful login.
var LoggOut;	// call function on logout clicked.

function doLogin ( event ) {
	if ( event ) {
		event.preventDefault();
	}
	console.log ( "doLogin");
	renderForm25 ( event ) ;
}
$("#login").click(doLogin);
$(".show-logged-in").hide();

LoggInDone = function ( JWTToken ) {
	console.log ( "LoggInDone");
	$("#body").html("<span></span>");
	// <a class="nav-link" href="#" id="login">Login</a>
	isLoggedIn = true;
	jwtToken = JWTToken;
	SetupJWTBerrer();
	$("#login").html("Logout");
	$("#login").click(LoggOut);
}

LoggOut = function ( event ) {
	console.log ( "LoggOut");
	if ( event ) {
		event.preventDefault();
	}
	isLoggedIn = false;
	jwtToken = "";
	xsrf_token = "" ;
	$("#login").html("Login");
	$("#login").click(doLogin);
	$(".show-anon").show();
	$(".show-logged-in").hide();
	renderMessage ( "Logged Out", "You are now logged out.<br>");
renderForm09(null);
	render5SecClearMessage();
}


function submitIt ( event, data, action, succ, erro ) {
	event.preventDefault();

	// xyzzy - add in _ran_ to data

	$.ajax({
		type: 'GET',
		url: action,
		data: data,
		success: function (data) {
			if ( succ ) {
				succ(data);
			}
			$("#output").text( JSON.stringify(data, null, 4) );
		},
		error: function(resp) {
			$("#output").text( "Error!"+JSON.stringify(resp) );
			if ( erro ) {
				erro(data);
			}
			// alert("got error status="+resp.status+" "+resp.statusText);
		}
	});
}
$("#getStatus").click(function(event){
	submitIt ( event, {}, "/api/v1/status",
		function(data) {
			$("#body").html( "<pre>"+JSON.stringify(data)+"</pre>" );
		}
	);
});


// Function that will allow us to know if Ajax uploads are supported
function supportAjaxUploadWithProgress() {
	return supportFileAPI() && supportAjaxUploadProgressEvents() && supportFormData();

	// Is the File API supported?
	function supportFileAPI() {
		var fi = document.createElement('INPUT');
		fi.type = 'file';
		return 'files' in fi;
	};

	// Are progress events supported?
	function supportAjaxUploadProgressEvents() {
		var xhr = new XMLHttpRequest();
		return !! (xhr && ('upload' in xhr) && ('onprogress' in xhr.upload));
	};

	// Is FormData supported?
	function supportFormData() {
		return !! window.FormData;
	}
}

function displayFileUploadSupported() {

	// Actually confirm support
	if (supportAjaxUploadWithProgress()) {
		// Ajax uploads are supported!
		// Change the support message and enable the upload button
		var notice = document.getElementById('support-notice');
		var uploadBtn = document.getElementById('upload-button-id');
		notice.innerHTML = "Your browser supports HTML uploads. Go try me! :-)";
		uploadBtn.removeAttribute('disabled');

		// Init the Ajax form submission
		initFullFormAjaxUpload();

		// Init the single-field file upload
		initFileOnlyAjaxUpload();
	}

}

function initFullFormAjaxUpload() {
	var form = document.getElementById('form-id');
	form.onsubmit = function() {
		// FormData receives the whole form
		var formData = new FormData(form);

		// We send the data where the form wanted
		var action = form.getAttribute('action');

		// Code common to both variants
		sendXHRequest(formData, action);

		// Avoid normal form submission
		return false;
	}
}

function UploadTheFile() {
	var formData = new FormData();

	// Since this is the file only, we send it to a specific location
	var action = '/upload';

	// FormData only has the file
	var fileInput = document.getElementById('file-id');
	var file = fileInput.files[0];
	formData.append('file', file);
	formData.append('id', curDocID);
	formData.append("auth_key", g_auth_key);
	formData.append("app", "app.beefchain.com");

	// Code common to both variants
	sendXHRequest(formData, action);
}

function initFileOnlyAjaxUpload() {
	var uploadBtn = document.getElementById('upload-button-id');
	uploadBtn.onclick = function (evt) {
		var formData = new FormData();

		// Since this is the file only, we send it to a specific location
		var action = '/upload';

		// FormData only has the file
		var fileInput = document.getElementById('file-id');
		var file = fileInput.files[0];
		formData.append('file', file);
		formData.append('id', "15e42502-e7a5-44e2-6920-b410b9308412");

		// Code common to both variants
		sendXHRequest(formData, action);
	}
}

// Once the FormData instance is ready and we know
// where to send the data, the code is the same
// for both variants of this technique
function sendXHRequest(formData, uri) {
	// Get an XMLHttpRequest instance
	var xhr = new XMLHttpRequest();

	// Set up events
	xhr.upload.addEventListener('loadstart', onloadstartHandler, false);
	xhr.upload.addEventListener('progress', onprogressHandler, false);
	xhr.upload.addEventListener('load', onloadHandler, false);
	xhr.addEventListener('readystatechange', onreadystatechangeHandler, false);

	// Set up request
	xhr.open('POST', uri, true);

	// Fire!
	xhr.send(formData);
}

// Handle the start of the transmission
function onloadstartHandler(evt) {
	var div = document.getElementById('upload-status');
	div.innerHTML = 'Upload started.';
}

// Handle the end of the transmission
function onloadHandler(evt) {
	var div = document.getElementById('upload-status');
	div.innerHTML += '<' + 'br>File uploaded. Waiting for response.';
}

// Handle the progress
function onprogressHandler(evt) {
	var div = document.getElementById('progress');
	var percent = evt.loaded/evt.total*100;
	div.innerHTML = 'Progress: ' + percent + '%';
}

// Handle the response from the server
function onreadystatechangeHandler(evt) {
	var status, text, readyState;

	try {
		readyState = evt.target.readyState;
		text = evt.target.responseText;
		status = evt.target.status;
	}
	catch(e) {
		return;
	}

	if (readyState == 4 && status == '200' && evt.target.responseText) {
		var status = document.getElementById('upload-status');
		status.innerHTML += '<' + 'br>Success!';

		var result = document.getElementById('result');
		result.innerHTML = '<p>The server saw it as:</p><pre>' + evt.target.responseText + '</pre>';
	}
}

