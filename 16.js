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

function move_message(e, o, d) {
	document.getElementById(`${e}_${d}`).innerHTML = document.getElementById(`${e}_${o}`).innerHTML;
}

function update_message_buffer(e, c, m) {
	move_message(e, 1, 2);
	move_message(e, 0, 1);
	update(`${e}_0`, m);
	update(`${e}_count`, c);
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
var public_total = 0;
var private_seen = 0;
var private_total = 0;

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

server_link(1000, () => {
	if (public_seen < public_total) {
		ajax_get(`/message?r=public&i=${public_seen}`, response => {
			public_seen++;
			update_message_buffer("public_list", public_total, format_message(response));
		})
	}
});

server_link(1000, () => {
	if (private_seen < private_total) {
		ajax_get(`/message?r=${client_id}&i=${private_seen}`, response => {
			private_seen++;
			update_message_buffer("private_list", private_total, format_message(response));
		})
	}
});

server_link(500, () =>
	ajax_get("/messages?r=public", r => public_total = r))

server_link(500, () =>
	ajax_get(`/messages?r=private&a=${client_id}`, r => private_total = r))

function server_socket(url, onMessage) {
	if ('WebSocket' in window) {
		var socket = new WebSocket(url);
		socket.onopen = function(e) {
			console.log(`opening socket: ${url}`);
		};
		socket.onclose = function(e) {
			console.log(`closing socket: ${url}`);
		};
		socket.onerror = function(error) {
			console.log(`error for ${url}: ${error.message}`);
		};
		socket.onmessage = onMessage;
		return socket;
	} else {
		//	use AJAX polling as an alternative?
		return null;
	}
}

var monitor_feed;

window.onload = function() {
	update("id_banner", `contact ID: ${client_id}`);
	monitor_feed = server_socket("ws://localhost:3000/monitor", m => {
		var d = m.data.split("\t");
		update_message_buffer("event_list", d[0], d[1])
	})
}