<?php
header("Cache-control: no-cache, no-store, must-revalidate, max-age=0, s-maxage=0");
header("Pragma: no-cache");
header("Expires: 0");
?>

<html lang="en-US">
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width,initial-scale=1">
    <title>Speedtest Result Graph</title>
	<script src="chart.js"></script>
  </head>
  <body>
    <div>
	  <canvas id="myChart"></canvas>
	</div>
  </body>
  <script>
  const data = {
    labels: [1, 2, 3, 4, 5, 6, 7, 8],
    datasets: [

<?php
$url = "http://$_SERVER[HTTP_HOST]$_SERVER[REQUEST_URI]";
$url_components = parse_url($url);

parse_str($url_components['query'], $params);
$params = preg_replace("/[^a-z]+/", "", $params);

if (array_key_exists("testtype", $params))
{
	if($params['testtype'] == "upload")
	{
		$data1 = preg_replace('#\s+#',',',trim(file_get_contents("/var/www/localhost/htdocs/results/data/ooklaspeedtest_upload.txt")));
		?>
		{
		  label: 'Upload test',
		  backgroundColor: 'rgba(255, 99, 132, 0.5)',
		  borderColor: 'rgba(255, 99, 132, 0.5)',
		  data: [<?php echo $data1; ?>],
		  pointStyle: 'rect',
		  pointRadius: 10,
		  pointHoverRadius: 15,
		  tension: 0.4,
		  fill: 'start'
		}
		<?php
	} else if($params['testtype'] == "download")
	{
		$data1 = preg_replace('#\s+#',',',trim(file_get_contents("/var/www/localhost/htdocs/results/data/ooklaspeedtest_download.txt")));
		?>
		{
		  label: 'Download test',
		  backgroundColor: 'rgba(0, 180, 255, 0.5)',
		  borderColor: 'rgba(0, 180, 255, 0.5)',
		  data: [<?php echo $data1; ?>],
		  pointStyle: 'circle',
		  pointRadius: 10,
		  pointHoverRadius: 15,
		  tension: 0.4,
		  fill: 'start'
		}
		<?php
	} else
	{
		?>
		{
		  label: 'No data selected',
		  data: [0],
		}
		<?php
	}
} else {
	$data1 = preg_replace('#\s+#',',',trim(file_get_contents("/var/www/localhost/htdocs/results/data/ooklaspeedtest_u.txt")));
	$data2 = preg_replace('#\s+#',',',trim(file_get_contents("/var/www/localhost/htdocs/results/data/ooklaspeedtest_d.txt")));
	?>
	{
      label: 'Upload',
      backgroundColor: 'rgba(255, 99, 132, 0.5)',
      borderColor: 'rgba(255, 99, 132, 0.5)',
      data: [<?php echo $data1; ?>],
      pointStyle: 'rect',
      pointRadius: 10,
      pointHoverRadius: 15,
	  tension: 0.4
    },
	{
      label: 'Download',
      backgroundColor: 'rgba(0, 180, 255, 0.5)',
      borderColor: 'rgba(0, 180, 255, 0.5)',
      data: [<?php echo $data2; ?>],
      pointStyle: 'circle',
      pointRadius: 10,
      pointHoverRadius: 15,
	  tension: 0.4,
    }<?php
}

?>

]
  };

	const plugin = {
	  id: 'custom_canvas_background_color',
	  beforeDraw: (chart) => {
		const ctx = chart.canvas.getContext('2d');
		ctx.save();
		ctx.globalCompositeOperation = 'destination-over';
		ctx.fillStyle = 'white';
		ctx.fillRect(0, 0, chart.width, chart.height);
		ctx.restore();
	  }
	};

  let delayed;
	const config = {
	  type: 'line',
	  data: data,
	  plugins: [plugin],
	  options: {
		animation: {
		  onComplete: () => {
			delayed = true;
		  },
		  delay: (context) => {
			let delay = 0;
			if (context.type === 'data' && context.mode === 'default' && !delayed) {
			  delay = context.dataIndex * 75 + context.datasetIndex * 75;
			}
			return delay;
		  },
		},
		scales: {
		  x: {
			stacked: true,
			title: {
				display: true,
				text: "Test order (right is most recent)"
			  },
		  },
		  y: {
			title: {
				display: true,
				text: "Mbps"
			  },
		  }
		},
		hoverBackgroundColor: 'green'
	  }
	};
  
  const myChart = new Chart(
    document.getElementById('myChart'),
    config
  );
  </script>
</html>
