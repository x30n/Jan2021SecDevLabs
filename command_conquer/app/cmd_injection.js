var http = require('http');

var child_process = require('child_process')
var express = require('express');
var app = express();

app.get('/', function(req, res){
    child_process.exec(
	'ping -c 1 ' + req.query.host,
	function(err, data) {
	    console.log('err: ', err)
	    res.send(data);
	});
});

app.listen(3000);
