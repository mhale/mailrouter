// Shows or hides the username and password fields when authentication type is changed.
$("#authentication").change(function() {
	if ($(this).val() == "crammd5") {
		$("#password").attr("placeholder", "secret");
		$("#password-label").text("Secret");
		$("#password-group").removeClass("hidden").addClass("show");
		$("#username-group").removeClass("hidden").addClass("show");
	} else if ($(this).val() == "plain") {
		$("#password").attr("placeholder", "password");
		$("#password-label").text("Password");
		$("#password-group").removeClass("hidden").addClass("show");
		$("#username-group").removeClass("hidden").addClass("show");
	} else {
		$("#password-group").removeClass("show").addClass("hidden");
		$("#username-group").removeClass("show").addClass("hidden");
	}
});

// Handles "data-method" on links such as:
// <a href="/routes/b25f7ee5-b755-11e3-8126-4a5b3b8c74a2" data-method="delete" rel="nofollow" data-confirm="Are you sure?">Delete</a>
$('[data-method]').click(function() {
	if (confirm($(this).attr('data-confirm'))) {
		var form = $('<form method="post" action="' + $(this).attr('href') + '"></form>');
		var metadataInput = '<input name="_method" value="' + $(this).attr('data-method') + '" type="hidden" />';
		form.hide().append(metadataInput).appendTo('body');
		form.submit();
	}
	return false;
});
