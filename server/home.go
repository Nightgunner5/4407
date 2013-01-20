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
var tileSize = 32;

var requestAnimationFrame = window.requestAnimationFrame ||
	window.mozRequestAnimationFrame ||
	window.webkitRequestAnimationFrame ||
	window.msRequestAnimationFrame ||
	function(f){ setTimeout(f, 33); };

var currentLevel = 0, map = [], atmos = [];

var ws;
function connect() {
	ws = new WebSocket('ws://' + location.host + '/ws');
	ws.onmessage = function(e) {
		dispatch(JSON.parse(e.data));
	};
	ws.onclose = function() {
		setTimeout(connect, 1000);
	};
}
connect();

function dispatch(p) {
	if ('Atmos' in p) {
		atmos = p.Atmos;
		return;
	}
	if ('Map' in p) {
		map = p.Map;
		map.forEach(function(t) {
			var icon = 0;
			map.forEach(function(tt) {
				if (tt[2] != t[2]) return;
				if (tt[0] == t[0] && tt[1] == t[1]-1) icon |= 1;
				if (tt[0] == t[0] && tt[1] == t[1]+1) icon |= 2;
				if (tt[0] == t[0]+1 && tt[1] == t[1]) icon |= 4;
				if (tt[0] == t[0]-1 && tt[1] == t[1]) icon |= 8;
			});
			t.push(icon);
		});
		return;
	}
	console.log(p);
}

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
		var x = Math.round(t[0]*size - size/2 - offsetX*size + w/2);
		var y = Math.round(t[1]*size - size/2 - offsetY*size + h/2);
		ctx.drawImage(tile[t[2]], tileSize*t[3], 0, tileSize, tileSize, x, y, size, size);
	});
	atmos.forEach(function(t) {
		var x = Math.round(t.X*size - size/2 - offsetX*size + w/2);
		var y = Math.round(t.Y*size - size/2 - offsetY*size + h/2);
		ctx.fillStyle = 'rgba(' + Math.round(Math.min(Math.max(t.Temp - 100, 0), 255)) + ', 128, ' + Math.round(Math.min(Math.max(300 - t.Temp, 0), 255)) + ', 0.2)';
		ctx.fillRect(x, y, size, size);
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
