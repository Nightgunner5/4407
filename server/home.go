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
var moveUp = false, moveDown = false, moveLeft = false, moveRight = false;
window.onkeydown = function(e) {
	switch (e.which) {
	case 38: // up
		moveUp = true;
		return;
	case 40: // down
		moveDown = true;
		return;
	case 37: // left
		moveLeft = true;
		return;
	case 39: // right
		moveRight = true;
		return;
	}
};
window.onkeyup = function(e) {
	switch (e.which) {
	case 38: // up
		moveUp = false;
		return;
	case 40: // down
		moveDown = false;
		return;
	case 37: // left
		moveLeft = false;
		return;
	case 39: // right
		moveRight = false;
		return;
	}
};

var lastMove = new Date().getTime();
setInterval(function() {
	var now = new Date().getTime();
	var delta = (now - lastMove) / 1000;
	lastMove = now;

	var dx = 0, dy = 0;
	if (moveLeft) {
		dx -= delta;
	}
	if (moveRight) {
		dx += delta;
	}
	if (moveUp) {
		dy -= delta;
	}
	if (moveDown) {
		dy += delta;
	}
	if (dx < -1) dx = -1;
	if (dx > 1) dx = 1;
	if (dy < -1) dy = -1;
	if (dy > 1) dy = 1;
	if (Math.round(offsetX + dx) != Math.round(offsetX) ||
		Math.round(offsetY + dy) != Math.round(offsetY)) {
		if (!open(Math.round(offsetX + dx), Math.round(offsetY + dy)))
			return;

		ws.send(JSON.stringify({Position: {
			X: Math.round(offsetX + dx),
			Y: Math.round(offsetX + dx)
		}}));
	}
	offsetX += dx;
	offsetY += dy;
}, 25);

var tileSize = 32, statusCond = new Image();
statusCond.src = '/icon/status-cond.png';
var tile = [new Image(), new Image(), new Image(), new Image(), new Image(), new Image()];
for (var i = 0; i < tile.length; i++) {
	tile[i].src = '/tile/' + i + '.png';
}

var requestAnimationFrame = window.requestAnimationFrame ||
	window.mozRequestAnimationFrame ||
	window.webkitRequestAnimationFrame ||
	window.msRequestAnimationFrame ||
	setTimeout;

var currentLevel = 0, map = [], atmos = [];

function open(x, y) {
	var o = true;
	map.forEach(function(t) {
		if (t[0] == x && t[1] == y && (t[2] == 1 || t[2] == 3)) {
			o = false;
		}
	});
	return o;
}

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
		var merge = function(t) {
			if (t == 5) return 2;
			return t;
		};
		map.forEach(function(t) {
			var icon = 0;
			map.forEach(function(tt) {
				if (merge(tt[2]) != merge(t[2])) return;
				if (tt[0] == t[0] && tt[1] == t[1]-1) icon |= 1;
				if (tt[0] == t[0] && tt[1] == t[1]+1) icon |= 2;
				if (tt[0] == t[0]+1 && tt[1] == t[1]) icon |= 4;
				if (tt[0] == t[0]-1 && tt[1] == t[1]) icon |= 8;
			});
			t[3] = icon;
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
	requestAnimationFrame(paint, 33);

	var w = ctx.canvas.width, h = ctx.canvas.height;
	var centerX = Math.round(w/2);
	var centerY = Math.round(h/2);

	ctx.fillStyle = '#000';
	ctx.fillRect(0, 0, w, h);

	var size = round(h / 16);
	var s2 = size/2;

	var currentTile = 0;
	map.forEach(function(t) {
		var x = Math.round(t[0]*size - s2 - offsetX*size + centerX);
		var y = Math.round(t[1]*size - s2 - offsetY*size + centerY);
		ctx.drawImage(tile[t[2]], tileSize*t[3], 0, tileSize, tileSize, x, y, size, size);
		if (t[0] == Math.round(offsetX) && t[1] == Math.round(offsetY)) {
			currentTile = t[2];
		}
	});

	var currentAtmos;
	atmos.forEach(function(t) {
		var x = Math.round(t.X*size - s2 - offsetX*size + centerX);
		var y = Math.round(t.Y*size - s2 - offsetY*size + centerY);
		ctx.fillStyle = 'rgba(' + Math.round(Math.min(Math.max(t.Temp - 100, 0), 255)) + ', 128, ' + Math.round(Math.min(Math.max(300 - t.Temp, 0), 255)) + ', 0.2)';
		ctx.fillRect(x, y, size, size);
		if (t.X == Math.round(offsetX) && t.Y == Math.round(offsetY)) {
			currentAtmos = t;
		}
	});

	ctx.fillStyle = '#000';
	/*map.forEach(function(t) {
		if (t[2] == 1) {
			var x = Math.round(t[0]*size - offsetX*size + centerX);
			var y = Math.round(t[1]*size - offsetY*size + centerY);
			var dx = centerX-x, dy = centerY-y;
			if (x < centerX) {
				ctx.beginPath();
				ctx.moveTo(x+s2, y-s2);
				ctx.lineTo(x+s2, y+s2);
				ctx.lineTo(x+s2+(s2-dx)*1000, y+s2+(s2-dy)*1000);
				ctx.lineTo(x+s2+(s2-dx)*1000, y-s2+(-s2-dy)*1000);
				ctx.fill();
			}
			if (x > centerX) {
				ctx.beginPath();
				ctx.moveTo(x-s2, y-s2);
				ctx.lineTo(x-s2, y+s2);
				ctx.lineTo(x-s2+(-s2-dx)*1000, y+s2+(s2-dy)*1000);
				ctx.lineTo(x-s2+(-s2-dx)*1000, y-s2+(-s2-dy)*1000);
				ctx.fill();
			}
			if (y < centerY) {
				ctx.beginPath();
				ctx.moveTo(x-s2, y+s2);
				ctx.lineTo(x+s2, y+s2);
				ctx.lineTo(x+s2+(s2-dx)*1000, y+s2+(s2-dy)*1000);
				ctx.lineTo(x-s2+(-s2-dx)*1000, y+s2+(s2-dy)*1000);
				ctx.fill();
			}
			if (y > centerY) {
				ctx.beginPath();
				ctx.moveTo(x-s2, y-s2);
				ctx.lineTo(x+s2, y-s2);
				ctx.lineTo(x+s2+(s2-dx)*1000, y-s2+(-s2-dy)*1000);
				ctx.lineTo(x-s2+(-s2-dx)*1000, y-s2+(-s2-dy)*1000);
				ctx.fill();
			}
		}
	});*/

	if (currentAtmos) {
		var x = Math.round(w - size * 1.1);
		var y = Math.round(size * 0.1);
		if (currentAtmos.Temp > 325) {
			ctx.drawImage(statusCond, tileSize*1, 0, tileSize, tileSize, x, y, size, size);
			x -= size;
		}
		if (currentAtmos.Temp < 270) {
			ctx.drawImage(statusCond, tileSize*2, 0, tileSize, tileSize, x, y, size, size);
			x -= size;
		}
		if (currentAtmos.Oxygen - currentAtmos.CarbonDioxide < 5 && currentTile != 4) {
			ctx.drawImage(statusCond, tileSize*0, 0, tileSize, tileSize, x, y, size, size);
			x -= size;
		}
		if (currentAtmos.Plasma + currentAtmos.NitrousOxide > 1) {
			ctx.drawImage(statusCond, tileSize*3, 0, tileSize, tileSize, x, y, size, size);
			x -= size;
		}
	}

	ctx.fillRect(centerX-s2, centerY-s2, size, size);
}
requestAnimationFrame(paint, 33);
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
