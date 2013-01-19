package main

import (
	"net/http"
)

var Home = []byte(`<!DOCTYPE html>
<html>
<head>
<title>4407 | Editor</title>
<style>
canvas {
	background: #000;
}
body {
	padding-bottom: 40px;
}
#toolbar {
	position: fixed;
	bottom: 8px;
	left: 8px;
	right: 8px;
}
</style>
</head>
<body>
<canvas width="1" height="1"></canvas>
<div id="toolbar"></div>
<script>
var c = document.querySelector('canvas');
var ctx = c.getContext('2d');

var tileSize = 16;
var tileColor = ['#000', '#444', '#ccc', '#448', '#c44', '#aaa'];
var chosenTile = 1;

function levels() {
	var xhr = new XMLHttpRequest();
	xhr.open('GET', '/levels', true);
	xhr.addEventListener('load', function() {
		var toolbar = document.querySelector('#toolbar');
		toolbar.innerHTML = '';
		for (var i = 0, l = parseInt(xhr.responseText); i < l; i++) {
			var b = document.createElement('button');
			b.innerHTML = 'Level ' + i;
			(function(i) {
				b.onclick = function() {
					level(i);
				};
			})(i);
			toolbar.appendChild(b);
		}
		var b = document.createElement('button');
		b.innerHTML = 'New Level';
		b.onclick = function() {
			newLevel();
		};
		toolbar.appendChild(b);
		tileColor.forEach(function(color, i) {
			var b = document.createElement('button');
			b.innerHTML = 'C';
			b.style.color = '#fff';
			b.style.textShadow = '0 1px 2px #000';
			b.style.background = color;
			b.onclick = function() {
				chosenTile = i;
			};
			toolbar.appendChild(b);
		});

		var b = document.createElement('button');
		b.innerHTML = 'Compile and Save';
		b.onclick = function() {
			compileAndSave();
		};
		toolbar.appendChild(b);
	}, false);
	xhr.send();
}
levels();

function paintTile(x, y, t) {
	ctx.fillStyle = tileColor[t];
	ctx.fillRect(x * tileSize, y * tileSize, tileSize, tileSize);
}

function level(i) {
	var xhr = new XMLHttpRequest();
	xhr.open('GET', '/level/' + i, true);
	xhr.addEventListener('load', function() {
		var l = JSON.parse(xhr.responseText);
		var minX = 0, minY = 0, maxX = 0, maxY = 0;
		l.forEach(function(t) {
			if (t[0] < minX)
				minX = t[0];
			if (t[1] < minY)
				minY = t[1];
			if (t[0] > maxX)
				maxX = t[0];
			if (t[1] > maxY)
				maxY = t[1];
		});
		minX -= 8;
		minY -= 8;
		maxX += 8;
		maxY += 8;

		c.width = (maxX - minX + 1) * tileSize;
		c.height = (maxY - minY + 1) * tileSize;
		ctx = c.getContext('2d');
		l.forEach(function(t) {
			paintTile(t[0] - minX, t[1] - minY, t[2]);
		});

		c.onclick = function(e) {
			var x = Math.floor(e.offsetX / tileSize) + minX;
			var y = Math.floor(e.offsetY / tileSize) + minY;

			setTile(i, x, y);
		};
	}, false);
	xhr.send();
}

function setTile(z, x, y) {
	var xhr = new XMLHttpRequest();
	xhr.open('GET', '/set/' + z + '/' + x + '/' + y + '/' + chosenTile, true);
	xhr.addEventListener('load', function() {
		level(z);
	}, false);
	xhr.send();
}

function newLevel() {
	var xhr = new XMLHttpRequest();
	xhr.open('GET', '/levels/new', true);
	xhr.addEventListener('load', function() {
		levels();
		level(parseInt(xhr.responseText));
	}, false);
	xhr.send();
}

function compileAndSave() {
	var xhr = new XMLHttpRequest();
	xhr.open('GET', '/save', true);
	xhr.send();
}
</script>
</body>
</html>
`)

func home(w http.ResponseWriter) {
	_, err := w.Write(Home)
	handle(err)
}
