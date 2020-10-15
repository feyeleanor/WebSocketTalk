function print(e, m) {
	document.getElementById(e).innerHTML += `<div>${m}</div>`;
}

{{range $c, $v := .CallBridge}}
	function {{$c}}() {
		var xhttp = new XMLHttpRequest();
		xhttp.onreadystatechange = function() {
			if (this.readyState == 4 && this.status == 200) {
				print("event_log", JSON.parse(this.responseText));
			}
		};
		xhttp.open("GET", "{{$c}}", true);
		xhttp.send();
	}
{{end}}
