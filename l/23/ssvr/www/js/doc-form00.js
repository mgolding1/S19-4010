
function submitForm00 ( event ) {
    event.preventDefault(); 

	console.log ( "Click of submit button for form00 - paper form registration" );

	var data = {
	   	  "email"			: $("#email").val()
	   	, "real_name"		: $("#real_name").val()
	   	, "phone_number"	: $("#phone_number").val()
	   	, "address_usps"	: $("#address_usps").val()
	   	, "note"			: $("#note").val()
		, "auth_key"		: g_auth_key
		, "__method__"		: "POST"
		, "_ran_" 			: ( Math.random() * 10000000 ) % 10000000
	};
	submitItData ( event, data, "/api/v1/documents", function(data){
		console.log ( "data=", data );
		curDocID = data.id;
		UploadTheFile();
	});
}

function renderForm00 ( event ) {
	render5SecClearMessage();
	var form = [ ''
		,'<div>'
			,'<div class="row">'
				,'<div class="col-sm-10">'
					,'<div class="card bg-default">'
						,'<div class="card-header"><h2>Save Document</h2></div>'
						,'<div class="card-body">'
							,'<form id="form00">'
								,'<div class="form-group">'
									,'<label for="ranch_name">Persons Name</label>'
									,'<input type="text" class="form-control" id="real_name" name="real_name" >'
								,'</div>'
								,'<div class="form-group">'
									,'<label for="email">Email Address</label>'
									,'<input type="text" class="form-control" id="email" name="email" >'
								,'</div>'

								,'<div class="form-group">'
									,'<label for="ranch_name">Phone Number</label>'
									,'<input name="phone_number" id="phone_number" 	type="text" 	class="form-control"  />'
								,'</div>'
								,'<div class="form-group">'
									,'<label for="ranch_name">Mailing Address</label>'
									,'<input name="address_usps" id="address_usps" 	type="text" 	class="form-control"  />'
								,'</div>'
								,'<div class="form-group">'
									,'<label for="ranch_name">Other/Notes</label>'
									,'<input name="note" id="note" 					type="text" 	class="form-control"  />'
								,'</div>'

								,'<div class="form-group">'
									,'<div class="custom-file">'
										,'<label id="file-label" for="file-id" class="custom-file-label">Upload Document</label>'
										,'<input type="file" class="custom-file-input" id="file-id" name="file">'				
									,'</div>'
								,'</div>'
								,'<button type="button" class="btn btn-primary" id="form00-submit">Submit</button>'
							,'</form>'
						,'</div>'
					,'</div>'
				,'</div>'
			,'</div>'
		,'</div>'
  		,'<p id="upload-status"></p>'
  		,'<p id="progress"></p>'
  		,'<pre id="result"></pre>'
	].join("\n");
	$("#body").html(form);
	$("#form00-submit").click(submitForm00);
	$('#file-id').on('change',function(){
		var fileName = $(this).val(); 
		if ( fileName.startsWith('C:\\fakepath\\') ) {
			fileName = fileName.substring(12);
		}
		$('#file-label').html(fileName); 
	});
}
$("#form00-render").click(renderForm00); 	

