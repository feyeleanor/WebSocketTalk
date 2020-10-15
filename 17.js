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
	move_message(e, 1, 2);
	move_message(e, 0, 1);
	update(`${e}_0`, m);
	update(`${e}_count`, c);
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
var public_total = 0;
var private_seen = 0;
var private_total = 0;
var events_total = 0;

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
		ajax_get(`/message?r=0&i=${public_seen}`, response => {
			if (response.length > 0) {
				public_seen++;
				update_message_buffer("public_list", public_total, format_message(JSON.parse(response)));
			}
		})
	}
});

server_link(1000, () => {
	if (private_seen < private_total) {
		ajax_get(`/message?r=${client_id}&i=${private_seen}`, response => {
			if (response.length > 0) {
				private_seen++;
				update_message_buffer("private_list", private_total, format_message(JSON.parse(response)));
			}
		})
	}
});

function server_socket(url, onMessage) {
	var socket = new WebSocket(url);
	socket.onerror = function(error) {
		console.log(`error for ${url}: ${error.message}`);
	};
	socket.onmessage = onMessage;
	return socket;
}

var monitor = null;

window.onload = function() {
	monitor = server_socket("ws://localhost:3000/register", m => {
		client_id = JSON.parse(m.data);
		update("id_banner", client_id);
		update("public_list_count", public_total);
		update("private_list_count", private_total);
		update("event_list_count", events_total);

		monitor.onmessage = function(m) {
			var d = JSON.parse(m.data);
			events_total = d[0];
			if (d[1] != "") {
				update_message_buffer("event_list", events_total, `${events_total}: ${d[1]}`);
			}
			public_total = d[2];
			private_total = d[3];
		}
	})
}