package main

import (
	"net/http"
)

var Home = []byte(`<!DOCTYPE html>
<html>
<head>
	<title>4407 | Editor</title>
</head>
<body>
<script>
function level(i) {
	var xhr = new XMLHttpRequest()
	xhr.open('GET', '/level/' + i, true)
	xhr.addEventListener('load', function() {
		var l = JSON.parse(xhr.responseText)
		console.log(l)
	}, false)
	xhr.send()
}
</script>
</body>
</html>
`)

func home(w http.ResponseWriter) {
	_, err := w.Write(Home)
	handle(err)
}
