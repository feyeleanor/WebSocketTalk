function print(e, m) {
	document.getElementById(e).innerHTML += `<div>${m}</div>`;
}

{{range $c, $v := .Commands}}
	function {{$c}}() {
		var xhttp = new XMLHttpRequest();
		xhttp.onreadystatechange = function() {
			if (this.readyState == 4 && this.status == 200) {
				console.log(this.responseText);
				print("event_log", this.responseText);
			}
		};
		xhttp.open("GET", "{{$c}}", true);
		xhttp.send();
	}
{{end}}
