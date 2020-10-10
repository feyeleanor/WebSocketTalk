function read_var(v) {
	return document.getElementById(v).value;
}

function print(e, m) {
	document.getElementById(e).innerHTML += "<div>" + m + "</div>";
}

function format_message(t) {
	var m = t.split("\n");
	return `<div><h2>${m[0]}</h2><div>${m[1]}</div><div>${m[2]}</div></div>`;
}

function post_comment() {
	var xhttp = new XMLHttpRequest();
	xhttp.onreadystatechange = function() {
		if (this.readyState == 4 && this.status == 200) {
			console.log(this.responseText);
			print("message_list", format_message(this.responseText));
			var f = document.forms["addMessage"];
			f.author.value = "";
			f.message.value = "";
		}
	};

	xhttp.open("POST", "", true);
	xhttp.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
	xhttp.send(`a=${read_var('author')}&m=${read_var('message')}`);
}