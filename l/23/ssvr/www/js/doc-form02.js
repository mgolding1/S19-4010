
// Render Document

function submitForm02 ( event ) {
    event.preventDefault(); // Totally stop stuff happening

	console.log ( "Click of submit button for form02 - render document" );
}

var theHistory ;
var theDocument ;

function renderForm02 ( event ) {
	render5SecClearMessage();
	if ( ! isLoggedIn ) {
		$("#form25-render").click(renderForm25); 	// Attach to link to paint the partial
		return
	}
	var form = [ ''
		,'<div>'
			,'<div class="row">'
				,'<div class="col-sm-12">'
					,'<div class="card bg-default">'
						,'<div class="card-header"><h2>Document Information</h2></div>'
						,'<div class="card-body">'
							,theHistory
						,'</div>'
					,'</div>'
				,'</div>'
			,'</div>'
		,'</div>'
		,'<div>'
			,'<div class="row">'
				,'<div class="col-sm-12">'
					,theDocument
				,'</div>'
			,'</div>'
		,'</div>'
	].join("\n");
	$("#body").html(form);
	// $("#form02-submit").click(submitForm02);
}

var pdf_url ;

function renderDoc(id){
	if ( ! isLoggedIn ) {
		$("#form25-render").click(renderForm25); 	// Attach to link to paint the partial
		return
	}

	// take list of data, pull out "id"
	// if .pdf, if image etc. -- set theDocument for correct data to render.
	// remeber document history! -- When Signed, Digital Signature etc.
	theHistory = "<h1> History of Document </h1>";
	theDocument = "<h1> Placeholder form02 </h1>";
	var found = false;
	var data = {};
	for ( var i = 0, mx = search_data.length; i < mx; i++ ) {
		if ( id === search_data[i].id ) {
			found = true;
			data = search_data[i];
		}
	}
	if ( ! found ) {
		theHistory = "<h1> Document Not Found </h1>";
		theDocument = "";
		renderForm02();
	} else {
		theHistory = [ ''
			,'<table class="doc-history">'
				,'<tr>'
					,'<th>Name</th><td>',data.real_name,'</td>'
				,'</tr>'
				,'<tr>'
					,'<th>Created</th><td>',data.created,'</td>'
				,'</tr>'
				,'<tr>'
					,'<th>Address</th><td>',data.address_usps,'</td>'
				,'</tr>'
				,'<tr>'
					,'<th>File Name</th><td>',data.orig_file_name,'</td>'
				,'</tr>'
				,'<tr>'
					,'<th>Digital Signature</th><td>',data.hash,"/",data.signature,'<br><b>Signature has been Verified on the Ethereum BlockChain</b></td>'	
				,'</tr>'
			,'</table>'
		].join("\n");
		theDocument = [ ''
			,'<iframe src="/viewer.html?file=',data.url_file_name,'" width="100%" height="800px"></iframe>'
		].join("\n");
	}

	// renderForm02();
	// TODO Do Validation of "Signature" on documet at this point!
	// 		mux.Handle("/api/v1/validate-document", http.HandlerFunc(HandleStatus)) // xyzzy - TODO - validate a document before paint of .pdf

	var data = {
		  "auth_key"		: g_auth_key
		, "id"				: id
		, "_ran_" 			: ( Math.random() * 10000000 ) % 10000000
	};
	submitItData ( event, data, "/api/v1/validate-document", function(data){
			if ( data && data.status && data.status == "success" ) {
				renderForm02();
			} else {
				console.log ( "ERROR: ", data );
				theHistory = [ ''
					,'<table class="doc-history">'
						,'<tr>'
							,'<th>Name</th><td>',data.real_name,'</td>'
						,'</tr>'
						,'<tr>'
							,'<th>Created</th><td>',data.created,'</td>'
						,'</tr>'
						,'<tr>'
							,'<th>Address</th><td>',data.address_usps,'</td>'
						,'</tr>'
						,'<tr>'
							,'<th>File Name</th><td>',data.orig_file_name,'</td>'
						,'</tr>'
						,'<tr>'
							,'<th>Digital Signature</th><td><br><b>Signature DID NOT Validate</b></td>'	
						,'</tr>'
					,'</table>'
				].join("\n");
				theDocument = [ ''
					,'<iframe src="/viewer.html?file=',data.url_file_name,'" width="100%" height="800px"></iframe>'
				].join("\n");
				renderForm02();
			}
		}, function(data) {
				console.log ( "ERROR: ", data );
				renderError ( "Network communication failed.", "Failed to communicate with the server." );
		}
	);

}



