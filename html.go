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
</style>
<body>
<div align="center">
<div class="main">
`

	footer = `
</div>
</div>
</body>
</html>
`
)
