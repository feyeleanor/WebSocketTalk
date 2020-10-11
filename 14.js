function read_var(v) {
	return document.getElementById(v).value;
}

function print(e, m) {
	document.getElementById(e).innerHTML += `<div>${m}</div>`;
}

function update(e, m) {
	document.getElementById(e).innerHTML = m;
}

function format_message(t) {
	var m = t.split("\n");
	var author = m.shift();
	var timestamp = m.shift();
	var message = m.shift();
	return `<h3>${author}</h3><div>${timestamp}</div><div>${message}</div>`;
}

function ajax_setup(f) {
	var xhttp = new XMLHttpRequest();
	xhttp.onreadystatechange = function() {
		if (this.readyState == 4 && this.status == 200) {
			f(xhttp);
		}
	};
	return xhttp;	
}

var comments_seen = 0;

function post_comment() {
	var xhttp = ajax_setup(x => {
		var f = document.forms["addMessage"];
		f.author.value = "";
		f.message.value = "";
	});

	xhttp.open("POST", "message", true);
	xhttp.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
	xhttp.send(`a=${read_var('author')}&m=${read_var('message')}`);
}

function ajax_get(url, response_handler) {
	var xhttp = ajax_setup(x => response_handler(x.responseText));
	xhttp.open("GET", url, true);
	xhttp.send();
}

function server_link(interval, f) {
	setInterval(f, interval)
}

server_link(1000, () =>
	ajax_get(`/message?i=${comments_seen}`, response => {
		print("message_list", format_message(response));
		comments_seen += 1;
	})
);

server_link(250, () =>
	ajax_get("/messages", response =>
		update("message_count", `messages on server: ${response}`)),
);