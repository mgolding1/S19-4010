

// Login Form

function submitForm25 ( event ) {
    event.preventDefault(); // Totally stop stuff happening

	console.log ( "Click of submit button for form25 - login" );

	var data = {
	   	  "username"		: $("#username").val()
		, "password"		: $("#password").val()
		, "__method__"		: "POST"
		, "_ran_" 			: ( Math.random() * 10000000 ) % 10000000
	};
	submitItData ( event, data, "/login", function(data){
		console.log ( "data=", data );
		if ( data && data.status && data.status == "success" ) {
			user_id = data.user_id; // sample: -- see bottom of file: www/js/pdoc-form02.js
			auth_token = data.auth_token;
			LoggInDone ( auth_token );
			$(".show-anon").hide();
			$(".show-logged-in").show();
			renderMessage ( "Successful Login", "You are now logged in<br>");
console.log ( "AAA", data );
		} else {
			console.log ( "ERROR: ", data );
			renderError ( "Failed to Login", data.msg );
		}
	}, function(data) {
		console.log ( "ERROR: ", data );
		renderError ( "Failed to Login - Network communication failed.", "Failed to communicate with the server." );
	}
	);
}

function renderForm25 ( event ) {
	var form = [ ''
		,'<div>'
			,'<div class="row">'
				,'<div class="col-sm-6">'
					,'<div class="card bg-default">'
						,'<div class="card-header"><h2>Login</h2></div>'
						,'<div class="card-body">'
							,'<form id="form01">'
								,'<input name="app"    		                	type="hidden" 	value="app.beefchain.com">'
								,'<input name="auth_key"               			type="hidden" 	value="1234">'
								,'<div class="form-group">'
									,'<label for="username">Email</label>'
									,'<input type="text" class="form-control" id="username" name="username"/>'
								,'</div>'
								,'<div class="form-group">'
									,'<label for="password">Password</label>'
									,'<input type="password" class="form-control" id="password" name="password"/>'
								,'</div>'
								,'<button type="button" class="btn btn-primary" id="form25-submit">Log In</button>'
							,'</form>'
						,'</div>'
					,'</div>'
				,'</div>'
			,'</div>'
		,'</div>'
	].join("\n");
	$("#body").html(form);
	// Add events
	$("#form25-submit").click(submitForm25);
	// xyzzy - additional click events forgot-pass, forgot-acct
}
$("#form25-render").click(renderForm25); 	// Attach to link to paint the partial


