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
	var m = t.split("\t");
	var author = m.shift();
	var timestamp = m.shift();
	var message = m.shift();
	return `<hr/><h3>From: ${author}</h3><div>Date: ${timestamp}</div><div>${message}</div>`;
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

function ajax_get(url, response_handler) {
	var xhttp = ajax_setup(x => response_handler(x.responseText));
	xhttp.open("GET", url, true);
	xhttp.send();
}

function ajax_post(xhttp, url, params) {
	xhttp.open("POST", url, true);
	xhttp.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
	xhttp.send(params);
}

const client_id = {{.Clients}};
var public_seen = 0;
var private_seen = 0;

function post_comment() {
	var xhttp = ajax_setup(x => {
		var f = document.forms["addMessage"];
		f.recipient.value = "";
		f.message.value = "";
	});
	ajax_post(xhttp, "message", `a=${client_id}&m=${read_var('message')}&r=${read_var('recipient')}`);
}

function server_link(interval, f) {
	setInterval(f, interval)
}

server_link(1000, () =>
	ajax_get(`/message?r=public&i=${public_seen}`, response => {
		print("public_list", format_message(response));
		public_seen++;
	})
);

server_link(1000, () =>
	ajax_get(`/message?r=${client_id}&i=${private_seen}`, response => {
		print("private_list", format_message(response));
		private_seen++;
	})
);

server_link(250, () =>
	ajax_get("/messages?r=public", response =>
		update("public_count", `messages on server: ${response}`)));

server_link(250, () =>
	ajax_get(`/messages?r=private&a=${client_id}`, response =>
		update("private_count", `messages on server: ${response}`)));

window.onload = function() {
	update("id_banner", `contact ID: ${client_id}`);
}