function print(e, m) {
	document.getElementById(e).innerHTML += `<div>${m}</div>`;
}

{{range $c, $v := .CallBridge}}
	function {{$c}}(method) {
		var xhttp = new XMLHttpRequest();
		xhttp.onreadystatechange = function() {
			if (this.readyState == 4 && this.status == 200) {
				console.log(this.responseText);
				print("event_log", this.responseText);
			}
		};
		switch(method) {
		case "GET":
			xhttp.open("GET", "{{$c}}", true);
			xhttp.send();
			break;
		case "POST":
			xhttp.open("POST", "{{$c}}", true);
			xhttp.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
			xhttp.send("c={{$c}}&{{$v}}=v");
			break;
		}
	}
{{end}}
