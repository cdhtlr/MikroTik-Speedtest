<?php
error_reporting(0);

header('Content-Type: text/html; charset=utf-8');
header('Cache-Control: no-store, no-cache, must-revalidate, max-age=0, s-maxage=0');
header('Cache-Control: post-check=0, pre-check=0', false);
header('Pragma: no-cache');
?>
<!DOCTYPE html>
<html>
    <head>
        <title>Speedtest Results</title>
        <style type="text/css">
            html,body{
                margin:0;
                padding:0;
                border:none;
                width:100%; min-height:100%;
            }
            html{
                background-color: hsl(198,72%,35%);
                font-family: "Segoe UI","Roboto",sans-serif;
            }
            body{
                background-color:#FFFFFF;
                box-sizing:border-box;
                width:100%;
                max-width:70em;
                margin:4em auto;
                box-shadow:0 1em 6em #00000080;
                padding:1em 1em 4em 1em;
                border-radius:0.4em;
            }
            h1,h2,h3,h4,h5,h6{
                font-weight:300;
                margin-bottom: 0.1em;
            }
            h1{
                text-align:center;
            }
            table{
                margin:2em 0;
                width:100%;
            }
            table, tr, th, td {
                border: 1px solid #AAAAAA;
            }
            th {
                width: 6em;
            }
            td {
                word-break: break-all;
            }
            div {
                margin: 1em 0;
            }
        </style>
    </head>
    <body>
		<p><a href="../index.html">Home</a></p>
		<p>Latest <b>Ookla Speedtest.net</b> results (in Megabits per second):</p>
		<ul>
		  <li>Full test results (<a href="results/data/ooklaspeedtest.txt">text</a>|<a href="results/graph.php">graph</a>)</li>
		  <li>Upload test results (<a href="results/data/ooklaspeedtest_upload.txt">text</a>|<a href="results/graph.php?testtype=upload">graph</a>)</li>
		  <li>Download test results (<a href="results/data/ooklaspeedtest_download.txt">text</a>|<a href="results/graph.php?testtype=download">graph</a>)</li>
		</ul>
    </body>
</html>
