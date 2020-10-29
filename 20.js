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

function move_message(e, o, d) {
	document.getElementById(`${e}_${d}`).innerHTML = document.getElementById(`${e}_${o}`).innerHTML;
}

function update_message_buffer(e, c, m) {
	move_message(e, 2, 3);
	move_message(e, 1, 2);
	move_message(e, 0, 1);
	update(`${e}_0`, m);
	update(`${e}_count`, c);
}

var server = null;
var client_id = 0;
var public_seen = 0;
var public_total = 0;
var private_seen = 0;
var private_total = 0;

function post_comment() {
	var m = {
		recipient: read_var('recipient'),
		content: read_var('message'),
	}
	var f = document.forms["addMessage"];
	f.recipient.value = "";
	f.message.value = "";
	console.log(`m = ${JSON.stringify(m)}`);
	server.send(JSON.stringify(m));
}

function server_socket(url, onMessage) {
	var socket = new WebSocket(url);
	socket.onerror = function(error) {
		console.log(`error for ${url}: ${error.message}`);
	};
	socket.onmessage = onMessage;
	return socket;
}

window.onload = function() {
	server = server_socket("ws://localhost:3000/register", m => {
		client_id = JSON.parse(m.data);
		update("id_banner", client_id);
		update("public_list_count", public_total);
		update("private_list_count", private_total);

		server.onmessage = function(m) {
			var d = JSON.parse(m.data);

			switch (d.shift().toLowerCase()) {
			case "status":
				public_total = d[0];
				update("public_list_count", public_total);
				private_total = d[1];
				update("private_list_count", private_total);
				break;

			case "broadcast":
				console.log(`m.data = ${m.data}`);
				update_message_buffer("public_list", public_total, format_message(d[0]));
				break;

			case "private":
				console.log(`m.data = ${m.data}`);
				update_message_buffer("private_list", private_total, format_message(d[0]));
				break;
			}
		}
	})
}