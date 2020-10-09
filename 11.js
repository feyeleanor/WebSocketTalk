{{range $c, $v := .Commands}}
	function {{$c}}() {
		var xhttp = new XMLHttpRequest();
		xhttp.onreadystatechange = function() {
			if (this.readyState == 4 && this.status == 200) {
				console.log(this.responseText);
				document.getElementById("event_log").innerHTML += "<span>" + this.responseText + "</span>";
			}
		};
		xhttp.open("GET", "{{$c}}", true);
		xhttp.send();
	}
{{end}}
