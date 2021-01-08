fs = require('fs');
fs.writeFile('/data/flag.txt', 'Flag GOES HERE!', function (err) {
  if (err) return console.log(err);
  console.log('Created file in /data/flag.txt');
});