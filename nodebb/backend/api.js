const { indexOf } = require('./events.js');
var events = require('./events.js');
const dataSandbox = '/data/';

exports.events = function (req, res) {
  res.json(events);
};

exports.event = function (req, res) {
  res.json(events[req.param.eventId]);
};


exports.hello = function (req, res) {
  var hello = "hello";
  res.json(hello);
}

exports.eventlist = function(req, res) {
  var f = req.params.eventFile;
  var e = require(dataSandbox + 'people.js');
    // res.json(e);
  if (!f) {
    res.json(e);
  } else {
    // Limit retrieved file to dataSandbox (/data/) for security
    f = dataSandbox + f;
    try {
      e = require(f);
      res.json(e);
    } catch {
      // Some data files aren't JSON so try simple read if we hit an exception above
      const fs = require('fs');
      fs.readFile(f, 'utf8' , (err, data) => {
        if (err) {
          res.json(err);
        }
        res.json(data);
      })
    }
  }
}