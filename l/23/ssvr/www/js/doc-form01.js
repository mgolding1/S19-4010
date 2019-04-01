
var search_data = [];

function submitForm01 ( event ) {
	if ( event ) {
		event.preventDefault(); // Totally stop stuff happening
	}

	$("#msg").html("");

	var data = {
		  "auth_key"		: g_auth_key
		, "_ran_" 			: ( Math.random() * 10000000 ) % 10000000
	};
	submitItData ( event, data, "/api/v1/documents", function(data){
			if ( data && data.status && data.status == "success" && data.data && data.data.length > 0 ) {
				search_data = data.data;
				renderForm01 ( event ); // next !!!!
			} else {
				console.log ( "ERROR: ", data );
				// renderError ( "Data Error", "Check that you have create at least 1 document.");
			}
		}, function(data) {
				console.log ( "ERROR: ", data );
				renderError ( "Failed to Register - Network communication failed.", "Failed to communicate with the server." );
		}
	);
}

function renderForm01 ( event ) {
	render5SecClearMessage();
	var thead = [ ''
		,'<tr>'
			,'<th>'
				,'Name'
			,'</th>'
			,'<th>'
				,'Email Address'
			,'</th>'
		,'</tr>'
	].join("\n");

	var rows = [];
	for ( var ii = 0, mx = search_data.length; ii < mx; ii++ ) {
		rows.push ( [ ''
			,'<tr>'
				,'<td>'
					,search_data[ii].real_name
				,'</td>'
				,'<td>'
					,search_data[ii].email
				,'</td>'
				,'<td>'
					,'<a href="#" class="btn btn-primary" onClick=\'xyzzy("'+ search_data[ii].id +'")\'>Show Document</a>'
				,'</td>'
			,'</tr>'
		].join("\n") );
	}
	var tbody = rows.join("\n");

	var formData = [ ''
		,'<div>'
			,'<div class="row">'
				,'<div class="col-sm-12">'
					,'<table class="table">'
						,'<thead id="thead">'
							,thead
						,'</thead>'
						,'<tbody id="tbody">'
							,tbody
						,'</tbody>'
					,'</table>'
				,'</div>'
			,'</div>'
		,'</div>'
	];

	var form = formData.join("\n");
	$("#body").html(form);
}

submitForm01(null);

