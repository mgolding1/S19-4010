
function submitForm09 ( event ) {
    event.preventDefault(); // Totally stop stuff happening

	console.log ( "Click of submit button" );
}

function renderForm09 ( event ) {
	renderClearMessage();
	var form = [ ''
				,'<div>'
					,'<div class="row">'
						,'<div class="col-sm-1"></div>'
						,'<div class="col-sm-4">'
						  ,'<div class="card mb-4 box-shadow bg-primary" style="height:300px;">'
								,'<div class="card-body">'
								  ,'<p class="card-text"><h2>Welcome to document signing</h2></p>'
								  ,'<div class="d-flex justify-content-between align-items-center">'
								  // ,'<button type="button" class="btn btn-sm btn-outline-secondary" id="form00-render-alt">Start Now</button>'
								  ,'</div>'
								,'</div>'
						  ,'</div>'
						,'</div>'
					,'</div>'
				  ,'</div>'
	].join("\n");
	$("#body").html(form);
}

