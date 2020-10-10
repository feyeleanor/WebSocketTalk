function read_var(v) {
	return document.getElementById(v).value;
}

function print(e, m) {
	document.getElementById(e).innerHTML += "<div>" + m + "</div>";
}

function format_message(t) {
	var m = t.split("\n");
	var author = m.shift();
	var timestamp = m.shift();
	var message = m.shift();
	return `<div><h2>${author}</h2><div>${timestamp}</div><div>${message}</div></div>`;
}

function post_comment() {
	var xhttp = new XMLHttpRequest();
	xhttp.onreadystatechange = function() {
		if (this.readyState == 4 && this.status == 200) {
			console.log(this.responseText);
			var f = document.forms["addMessage"];
			f.author.value = "";
			f.message.value = "";
		}
	};

	xhttp.open("POST", "message", true);
	xhttp.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
	xhttp.send(`a=${read_var('author')}&m=${read_var('message')}`);
}

var comments_seen = 0;
var polling_timer = setInterval(next_comment, 1000);

function next_comment() {
	var xhttp = new XMLHttpRequest();
	xhttp.onreadystatechange = function() {
		if (this.readyState == 4 && this.status == 200) {
			var m = this.responseText;
			console.log(m);
			print("message_list", format_message(m));
			comments_seen += 1;
		}
	};

	xhttp.open("GET", `/message?i=${comments_seen}`, true);
	xhttp.send();
}