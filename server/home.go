package main

import (
	"net/http"
)

var Home = []byte(`<!DOCTYPE html>
<html>
<head>
<title>4407</title>
<style>
body {
	padding: 0;
	margin: 0;
	overflow: hidden;
}
</style>
</head>
<body>
<canvas></canvas>
<script>
var ctx;
function resize() {
	var c = document.querySelector('canvas');
	c.width = window.innerWidth;
	c.height = window.innerHeight;
	ctx = c.getContext('2d');
	ctx.webkitImageSmoothingEnabled = false;
	ctx.mozImageSmoothingEnabled = false;
}
resize();
window.onresize = resize;
window.onkeydown = function(e) {
	switch (e.which) {
	case 38: // up
		--offsetY;
		return;
	case 40: // down
		++offsetY;
		return;
	case 37: // left
		--offsetX;
		return;
	case 39: // right
		++offsetX;
		return;
	}
};

var tile = [new Image(), new Image(), new Image(), new Image(), new Image(), new Image()];
for (var i = 0; i < tile.length; i++) {
	tile[i].src = '/tile/' + i + '.png';
}

var requestAnimationFrame = window.requestAnimationFrame ||
	window.mozRequestAnimationFrame ||
	window.webkitRequestAnimationFrame ||
	window.msRequestAnimationFrame ||
	function(f){ setTimeout(f, 33); };

var map = [];

var xhr = new XMLHttpRequest();
xhr.open('GET', '/level/0', true);
xhr.onload = function() {
	map = JSON.parse(xhr.responseText);
	map.sort(function(a, b) {
		if (a[1] == b[1])
			return a[0] - b[0];
		return a[1] - b[1];
	});
};
xhr.send();

function round(f) {
	for (var i = 1; ; i <<= 1)
		if (f <= i)
			return i;
}

var offsetX = 0, offsetY = 0;

function paint() {
	requestAnimationFrame(paint);

	var w = ctx.canvas.width, h = ctx.canvas.height;

	ctx.fillStyle = '#000';
	ctx.fillRect(0, 0, w, h);

	var size = round(h / 16);

	map.forEach(function(t) {
		ctx.drawImage(tile[t[2]], t[0]*size - size/2 - offsetX*size + w/2, t[1]*size - size/2 - offsetY*size + h/2, size, size);
	});
}
requestAnimationFrame(paint);
</script>
</body>
</html>`)

func init() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		w.Write(Home)
	})
}
