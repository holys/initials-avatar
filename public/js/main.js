/*Uses https://github.com/stewartlord/identicon.js and http://caligatio.github.io/jsSHA/  */
function get_identicon (text) {
	if (text == "") {
		text = "A";
	}
    $("#show_identicon")[0].src='https://initials.herokuapp.com/' + text;
}

