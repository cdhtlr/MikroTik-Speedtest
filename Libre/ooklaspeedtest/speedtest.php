<?php
header("Cache-control: no-cache, no-store, must-revalidate, max-age=0, s-maxage=0");
header("Pragma: no-cache");
header("Expires: 0");

$url = "http://$_SERVER[HTTP_HOST]$_SERVER[REQUEST_URI]";
$url_components = parse_url($url);

parse_str($url_components['query'], $params);
$params = preg_replace("/[^0-9]+/", "", $params);

$speedtest = "sudo /speedtest-cli -ud";

if (array_key_exists("serverid", $params))
{
	if($params['serverid'])
	{
		$speedtest = shell_exec($speedtest." --server ".$params['serverid']);
	} 
} else {
	$speedtest = shell_exec($speedtest);
}

preg_match('/Download:\s+(\d+\.\d+)/', $speedtest, $result);

if ( $result[1] != '' ){
    $filename = "/var/www/localhost/htdocs/results/data/ooklaspeedtest.txt";
    $filename_u = "/var/www/localhost/htdocs/results/data/ooklaspeedtest_u.txt";
    $filename_d = "/var/www/localhost/htdocs/results/data/ooklaspeedtest_d.txt";

    $count = shell_exec("wc -l < ".$filename);
    if( $count == 48 ) {
        shell_exec("sudo sed -i 1,6d ".$filename);
    } else if ($count > 48) {
		shell_exec("sudo rm -rf ".$filename);
	}

    $count_u = shell_exec("wc -l < ".$filename_u);
    if( $count_u == 8 ) {
        shell_exec("sudo sed -i 1,10d ".$filename_u);
    } else if ($count_u > 8) {
		shell_exec("sudo rm -rf ".$filename_u);
	}

    $count_d = shell_exec("wc -l < ".$filename_d);
    if( $count_d == 8 ) {
        shell_exec("sudo sed -i 1,10d ".$filename_d);
    } else if ($count_d > 8) {
		shell_exec("sudo rm -rf ".$filename_d);
	}
	
	preg_match('/Upload:\s+(\d+\.\d+)/', $speedtest, $result_u);
	preg_match('/Download:\s+(\d+\.\d+)/', $speedtest, $result_d);
	
	shell_exec("echo '".$result_u[1]."' | sudo tee -a ".$filename_u." > /dev/null 2>&1");
	shell_exec("echo '".$result_d[1]."' | sudo tee -a ".$filename_d." > /dev/null 2>&1");
	
    echo "<pre>";
    system("echo '".$speedtest."' | sudo tee -a ".$filename);
    echo "</pre>";
} else {
	echo "Error: Speedtest failed.";
}
?>