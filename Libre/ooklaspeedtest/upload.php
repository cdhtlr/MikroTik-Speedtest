<?php
header("Cache-control: no-cache, no-store, must-revalidate, max-age=0, s-maxage=0");
header("Pragma: no-cache");
header("Expires: 0");

$url = "http://$_SERVER[HTTP_HOST]$_SERVER[REQUEST_URI]";
$url_components = parse_url($url);

parse_str($url_components['query'], $params);
$params = preg_replace("/[^0-9]+/", "", $params);

$speedtest = "sudo /speedtest-cli -u";

if (array_key_exists("serverid", $params))
{
	if($params['serverid'])
	{
		$speedtest = shell_exec($speedtest." --server ".$params['serverid']);
	}
} else {
	$speedtest = shell_exec($speedtest);
}

preg_match('/Upload:\s+(\d+\.\d+)/', $speedtest, $result);

if ( $result[1] != '' ){
    $filename = "/var/www/localhost/htdocs/results/data/ooklaspeedtest_upload.txt";

    $count = shell_exec("wc -l < ".$filename);
    if( $count == 8 ) {
        shell_exec("sudo sed -i 1d ".$filename);
    } else if ($count > 8) {
		shell_exec("sudo rm -rf ".$filename);
	}

	if (array_key_exists("threshold", $params))
	{
		if($params['threshold'])
		{
			if($result[1] > $param['threshold']){
				echo "good";
			}
			else{
				echo "bad";
			}
		}
	}
	else
	{
		echo $result[1];
	}
    system("echo '".$result[1]."' | sudo tee -a ".$filename." > /dev/null 2>&1");
} else {
	echo "Error: Upload failed.";
}
?>
