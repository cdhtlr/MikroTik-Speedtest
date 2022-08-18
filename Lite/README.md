<b>DockerHub :</b> https://hub.docker.com/r/cdhtlr/mikrotik-speedtest

![](https://raw.githubusercontent.com/cdhtlr/MikroTik-Speedtest/main/Lite/demo/Main%20Page.png "Homepage")
![](https://raw.githubusercontent.com/cdhtlr/MikroTik-Speedtest/main/Lite/demo/Graph.png "Ookla Speedtest.net Result Graph")

<b>System requirements:</b>
- <b>Any linux PC</b> or <b>MikroTik RouterOS v7.1rc3 (and above)</b> (currently support Container for <b>ARM</b>, <b>ARM64</b>, <b>X86</b> and <b>X86-64</b> based machine, like the <i>CCR2004-16G-2S+</i> device with ARM64 architecture)
- <b>50MB of disk space remaining</b> (for PC) or <b>100MB of disk space remaining</b> (for MikroTik RouterBoard)
- <b>30MB of RAM space remaining</b> (could be more, depending on the workload)

===============================================================

<b>For use with PC (outside MikroTik):</b>

    docker run --restart=unless-stopped \ 
    --name speedtest -d \ #set name "speedtest" and Run container in background and print container ID
    -p 8080:80 \ #Host-port:Container-port to access the app inside this container via port 8080
    -e 'PING_IP=8.8.8.8' \ #set the IP address on the internet used to test internet configuration (optional)
    -e 'PING_DOMAIN=google.com' \ #set the domain on the internet used to test internet configuration (optional)
    -e 'TZ=Asia/Jakarta' \ #set Container Timezone (optional)
    -e 'CRON_FIELD=0 * * * *' \ #set the scheduler to run Ookla Speedtest.net and log the results automatically according to cron format (optional)
    -v '/home/data:/var/www/localhost/htdocs/results/data:rw' \ #Host-dir:Container-dir to make persistent database and records with read-write permission (optional)
    cdhtlr/mikrotik-speedtest:amd64 #Image for amd64 architecture

You can use the above example on <b>docker-compose</b>.

===============================================================

<b>For use as Container (inside MikroTik):</b>

![](https://raw.githubusercontent.com/cdhtlr/MikroTik-Speedtest/main/Lite/demo/Perform_Ookla_speedtest_directly_from_MikroTik_Container.png "Perform Ookla speedtest directly from MikroTik Container")

Pull this image to your computer (You can use any computer with any cpu architecture).

Example to pull image for ARM64 based MikroTik Router:

    docker pull cdhtlr/mikrotik-speedtest:arm64-lite

Save to TAR:

    docker save cdhtlr/mikrotik-speedtest:arm64-lite > speedtest.tar

then upload to your MikroTik Router.


Example configuration for MikroTik:

    /interface veth add name=veth1-speedtest address=192.168.1.2/24 gateway=192.168.1.1
    /interface bridge add name=bridge-docker
    /interface bridge port add interface=veth1-speedtest bridge=bridge-docker
    /ip address add address=192.168.1.1/24 interface=bridge-docker
    /ip firewall nat add chain=srcnat action=masquerade src-address=192.168.1.0/24
    /ip firewall nat add chain=dstnat action=dst-nat protocol=tcp to-addresses=192.168.1.2 to-ports=80 dst-port=8080
    /container envs add list=speedtest name=PING_IP value="8.8.8.8"
    /container envs add list=speedtest name=PING_DOMAIN value="google.com"
    /container envs add list=speedtest name=TZ value="Asia/Jakarta"
    /container envs add list=speedtest name=CRON_FIELD value="0 * * * *"
    /container mounts add name=speedtest dst=/var/www/localhost/htdocs/results/data src=disk1/speedtest
    /container add file=speedtest.tar interface=veth1-speedtest envlist=speedtest mounts=speedtest hostname=speedtest logging=yes

Check your container list in MikroTik Router:

    /container/print

    /container/start 0

0 is your container ID, please see the list of containers you got from the <i>/container/print</i> command.

===============================================================

<b>How to use:</b>

Now access speedtest from your web browser  (IP Router:Port).

For example http://192.168.1.2:8080 (<i>this is the main page containing the Ookla Speedtest.net links</i>)

You can go to http://192.168.1.2:8080/results/ to access Ookla Speedtest.net test results (8 most recent results).

Finally, you can now get your actual Ookla Speedtest.net bandwidth (in Megabits per second) using MikroTik Terminal and do some scripting.

For example:

    /tool fetch mode=http url="http://192.168.1.2:8080/download.php" keep-result=yes
    :delay 3
    :local result [/file get download.php contents]
    :log info "Your actual download bandwidth is $result Mbps"
the example output from the script above is: Your actual download bandwidth is 3 Mbps

or

    :local result [tool fetch mode=http url="http://192.168.1.2:8080/upload.php?serverid=37744&threshold=5" output=user as-value]
    :local result [:put ($result->"data")]
    :log info "Your actual upload bandwidth is currently $result"
the example output from the script above is: Your actual upload bandwidth is currently bad

Don't worry if MikroTik fails to execute the script due to a timeout, the process is already running in the Container.

All you need to do is run the following script which functions to retrieve the latest speedtest data obtained from the previous script. Run this script approximately 90 seconds after you run the previous script.

    :local result [tool fetch mode=http url="http://192.168.1.2:8080/result.php?testtype=upload" output=user as-value]
    :local result [:put ($result->"data")]
    :log info "Your actual upload bandwidth is $result Mbps"
the example output from the script above is: Your actual upload bandwidth is 3 Mbps

Filename list:
- <b>speedtest.php</b> (for full speedtest output). Optional parameter: <i>serverid</i>
- <b>upload.php</b> (upload only). Optional parameter: <i>serverid</i> and <i>threshold</i>
- <b>download.php</b> (download only). Optional parameter: <i>serverid</i> and <i>threshold</i>
- <b>result.php</b> (to retrieve the latest upload/download). Optional parameter: <i>threshold</i>

Parameter explanation:
- <b>testtype</b> is the parameter to specify the Ookla Speedtest.net type (upload or download)
- <b>serverid</b> is the parameter to specify the Ookla Speedtest.net server
- <b>threshold</b> is the parameter to specify the threshold value of the bandwidth generated by Ookla Speedtest.net. If the speedtest result is above the threshold then the output of this program is "good". However, if the speedtest result is below the threshold then the output of this program is "bad"

<p><b>Copyright notice:</b><br>The command-line speedtest used in this Docker Image is copyright of <a href="https://github.com/showwin/speedtest-go">ITO Shogo</a> (MIT License) and the chart graph used in this Docker Image is copyright of <a href="https://github.com/chartjs/Chart.js">ChartJS</a> (MIT License).</p>
