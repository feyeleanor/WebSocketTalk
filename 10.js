function print(m) {
	document.getElementById("event_log").innerHTML += "<div>" + m + "</div>";
}

{{range $c, $v := .Commands}}
	function {{$c}}() {
		var xhttp = new XMLHttpRequest();
		xhttp.onreadystatechange = function() {
			if (this.readyState == 4 && this.status == 200) {
				console.log(this.responseText);
				print(this.responseText);
			}
		};
		xhttp.open("GET", "{{$c}}", true);
		xhttp.send();
	}
{{end}}
