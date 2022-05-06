<?php
header("Cache-control: no-cache, no-store, must-revalidate, max-age=0, s-maxage=0");
header("Pragma: no-cache");
header("Expires: 0");

$url = "http://$_SERVER[HTTP_HOST]$_SERVER[REQUEST_URI]";
$url_components = parse_url($url);

parse_str($url_components['query'], $params);
$params = preg_replace("/[^a-z0-9]+/", "", $params);

if (array_key_exists("testtype", $params))
{
	if($params['testtype'])
	{
		$command = "";

		if($params['testtype'] == "upload"){
			$command = shell_exec("cat /var/www/localhost/htdocs/results/data/ooklaspeedtest_upload.txt | tail -1");
		}else if($params['testtype'] == "download"){
			$command = shell_exec("cat /var/www/localhost/htdocs/results/data/ooklaspeedtest_download.txt | tail -1");
		}else{
			echo "testtype parameter must be either <b>upload</b> or <b>download</b>";
		}
		
		preg_match('/(\d+\.\d+)/', $command, $result);

		if (array_key_exists("threshold", $params))
		{
			if($params['threshold'])
			{
				if($result[1] != ""){
					if($result[1] > $params['threshold']){
						echo "good";
					}
					else{
						echo "bad";
					}
				}
			}
		}
		else
		{
			echo $result[1];
		}
	}
}
else
{
	echo "testtype parameter must be defined";
}
?>
