package main

var (
	header = `
<html>
<head>
<style type="text/css">
.main {
  color: #333;
  font-family: Helvetica, arial, freesans, clean, sans-serif, "Segoe UI Emoji", "Segoe UI Symbol";
  size: 1.4em;
  text-align: left;
  width: 790px;
  display: block;
  border: 1px solid #ddd;
  border-top-left-radius: 7px;
  border-top-right-radius: 7px;
  padding: 20px;
}
code {
  padding: 0.2em 0;
  margin: 0;
  border-radius: 3px;
  background-color: #eee;
  font: 12px Consolas, "Liberation Mono", Menlo, Courier, monospace;
  font-size: 100%;
}
code:before, code:after {
  letter-spacing: -0.2em;
  content: "\00a0";
}
pre {
  padding: 16px;
  overflow: auto;
  font-size: 85%;
  line-height: 1.45;
  border-radius: 3px;
  margin-bottom: 16px;
  background-color: #eee;
}
pre code:before, pre code:after {
  content: normal;
}
li {
  line-height: 1.6;
}
#notify {
  color: #f00;
}
</style>
<script type="text/javascript">
function longpoll() {
  var req = new XMLHttpRequest();
  req.open('GET', document.URL, true);
  req.onreadystatechange = function() {
    if (req.readyState == 4) {
      if (req.status == 200) {
        replace(req.responseText);
        longpoll();
      } else {
        document.getElementById('notify').innerHTML = 'not connected';
      }
    }
  };
  req.send(null);
}
function replace(content) {
  document.write(content);
  document.close();
}
</script>
<body onload="longpoll()">
<div align="center">
<div class="main">
<div id="notify"></div>
`

	footer = `
</div>
</div>
</body>
</html>
`
)
