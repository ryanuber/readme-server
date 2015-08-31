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
  border-radius: 3px;
  margin-bottom: 16px;
  background-color: #eee;
}
pre code:before, pre code:after {
  content: normal;
}
li {
  padding-bottom: 10px;
}
h1 {
  width: 100%;
  border-bottom: 1px solid #ddd;
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
        document.open();
        document.write(req.responseText);
        document.close();
        longpoll();
      } else {
        document.getElementById('notify').innerHTML = 'not connected';
      }
    }
  };
  req.send(null);
}
document.onload = longpoll();
</script>
<body>
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
