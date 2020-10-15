function read_var(v) {
	return document.getElementById(v).value;
}

function print(e, m) {
	document.getElementById(e).innerHTML += `<div>${m}</div>`;
}

function update(e, m) {
	document.getElementById(e).innerHTML = m;
}

function format_message(m) {
	return `<hr/><h3>From: ${m.Author}</h3><div>Date: ${m.TimeStamp}</div><div>${m.Content}</div>`;
}

function ajax_setup(f) {
	var xhttp = new XMLHttpRequest();
	xhttp.onreadystatechange = function() {
		if (this.readyState == 4 && this.status == 200) {
			f(this);
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

var client_id = 0;
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
	ajax_get(`/message?r=0&i=${public_seen}`, response => {
		if (response.length > 0) {
			print("public_list", format_message(JSON.parse(response)));
			public_seen++;
		}
	})
);

server_link(1000, () =>
	ajax_get(`/message?r=${client_id}&i=${private_seen}`, response => {
		if (response.length > 0) {
			print("private_list", format_message(JSON.parse(response)));
			private_seen++;
		}
	})
);

server_link(1000, () =>
	ajax_get("/messages?r=0", response =>
		update("public_count", JSON.parse(response))));

server_link(1000, () =>
	ajax_get(`/messages?r=0&a=${client_id}`, response =>
		update("private_count", JSON.parse(response))));

window.onload = function() {
	ajax_get("/register", response => {
		client_id = JSON.parse(response);
		update("id_banner", client_id);
	});
}