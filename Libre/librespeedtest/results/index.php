<?php
error_reporting(0);

require 'telemetry_settings.php';
require_once 'telemetry_db.php';

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
		  <li>Full test results (<a href="data/ooklaspeedtest.txt">text</a>|<a href="graph.php">graph</a>)</li>
		  <li>Upload test results (<a href="data/ooklaspeedtest_upload.txt">text</a>|<a href="graph.php?testtype=upload">graph</a>)</li>
		  <li>Download test results (<a href="data/ooklaspeedtest_download.txt">text</a>|<a href="graph.php?testtype=download">graph</a>)</li>
		</ul>
		<br>
		<p>Latest <b>Librespeed</b> test results:</p>
		<?php
		$speedtests = getLatestSpeedtestUsers();
		if (false === $speedtests) {
			echo '<div>There was an error trying to fetch latest speedtest results.</div>';
		} elseif (empty($speedtests)) {
			echo '<div>Could not find any speedtest results in database.</div>';
		}
		foreach ($speedtests as $speedtest) {
			?>
			<table>
				<tr>
					<th>Date and time</th>
					<td><?= htmlspecialchars($speedtest['timestamp'], ENT_HTML5, 'UTF-8') ?></td>
				</tr>
				<tr>
					<th>IP Info</th>
					<td>
						<?= htmlspecialchars($speedtest['ip'], ENT_HTML5, 'UTF-8') ?><br/>
					</td>
				</tr>
				<tr>
					<th>User agent and locale</th>
					<td><?= htmlspecialchars($speedtest['ua'], ENT_HTML5, 'UTF-8') ?><br/>
						<?= htmlspecialchars($speedtest['lang'], ENT_HTML5, 'UTF-8') ?>
					</td>
				</tr>
				<tr>
					<th>Download speed</th>
					<td><?= htmlspecialchars($speedtest['dl'], ENT_HTML5, 'UTF-8') ?></td>
				</tr>
				<tr>
					<th>Upload speed</th>
					<td><?= htmlspecialchars($speedtest['ul'], ENT_HTML5, 'UTF-8') ?></td>
				</tr>
				<tr>
					<th>Ping</th>
					<td><?= htmlspecialchars($speedtest['ping'], ENT_HTML5, 'UTF-8') ?></td>
				</tr>
				<tr>
					<th>Jitter</th>
					<td><?= htmlspecialchars($speedtest['jitter'], ENT_HTML5, 'UTF-8') ?></td>
				</tr>
			</table>
			<?php
		}
        ?>
    </body>
</html>
