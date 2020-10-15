function read_var(v) {
	return document.getElementById(v).value;
}

function print(e, m) {
	document.getElementById(e).innerHTML += `<div>${m}</div>`;
}

function format_message(t) {
	var m = JSON.parse(t)
	return `<hr/><h3>From: ${m.Author}</h3><div>Date: ${m.TimeStamp}</div><div>${m.Content}</div>`;
}

function post_comment() {
	var xhttp = new XMLHttpRequest();
	xhttp.onreadystatechange = function() {
		if (this.readyState == 4 && this.status == 200) {
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

window.onload = function() {
	{{range $c, $v := .Messages}}
		print("message_list", format_message("{{$v.Author}}\t{{$v.TimeStamp}}\t{{$v.Content}}"))
	{{end}}
}