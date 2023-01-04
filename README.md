<b>DockerHub :</b> https://hub.docker.com/r/cdhtlr/mikrotik-speedtest

![](https://raw.githubusercontent.com/cdhtlr/MikroTik-Speedtest/main/Web_Demo.png "Web Demo")

<b>System requirements:</b>
- <b>Any linux PC</b> or <b>MikroTik RouterOS v7.6 (and above)</b> (Currently support Container for <b>ARM</b>, <b>ARM64</b>, <b>X86</b> and <b>X86-64</b> based machine, like the <i>CCR2004-16G-2S+</i> device with ARM64 architecture)
- <b>11MB of disk space left</b> (3MB for compressed TAR image, 8MB for uncompressed running image. You can delete the compressed TAR image once the image is successfully decompressed)
- <b>6MB RAM space left</b> (Could be more, depending on the workload)

===============================================================

<b>For use with PC (outside MikroTik):</b>

    docker run --restart=unless-stopped \ 
    --name speedtest -d \ #set name "speedtest" and Run container in background and print container ID
    -p 8080:80 \ #Host-port:Container-port to access the app inside this container via port 8080
    -e 'URL=https://jakarta.speedtest.telkom.net.id.prod.hosts.ooklaserver.net:8080/download?size=25000000' \ #url to download (optional)
    -e 'MAX_DLSIZE=1.0' \ #maximum size in MB (Megabytes) to download (optional)
    -e 'MIN_THRESHOLD=1.0' \ #download threshold in Mbps (Mbits per sec), to check for download speed condition (optional)
    cdhtlr/mikrotik-speedtest:amd64 #Image for amd64 architecture

You can use the above example on <b>docker-compose</b>.

===============================================================

<b>For use as Container (inside MikroTik):</b>

Pull this image to your computer (You can use any computer with any cpu architecture).

Example to pull image for ARM64 based MikroTik Router:

    docker pull cdhtlr/mikrotik-speedtest:arm64

Save to TAR:

    docker save cdhtlr/mikrotik-speedtest:arm64 > speedtest.tar

then upload to your MikroTik Router.


Example configuration for MikroTik:

    /interface veth add name=veth1-speedtest address=192.168.1.2/24 gateway=192.168.1.1
    /interface bridge add name=bridge-docker
    /interface bridge port add interface=veth1-speedtest bridge=bridge-docker
    /ip address add address=192.168.1.1/24 interface=bridge-docker
    /ip firewall nat add chain=srcnat action=masquerade src-address=192.168.1.0/24
    /ip firewall nat add chain=dstnat action=dst-nat protocol=tcp to-addresses=192.168.1.2 to-ports=80 dst-port=8080
    /container envs add list=speedtest name=URL value="https://jakarta.speedtest.telkom.net.id.prod.hosts.ooklaserver.net:8080/download?size=25000000"
    /container envs add list=speedtest name=MAX_DLSIZE value="1.0"
    /container envs add list=speedtest name=MIN_THRESHOLD value="1.0"
    /container add file=speedtest.tar interface=veth1-speedtest envlist=speedtest hostname=speedtest logging=yes

Check your container list in MikroTik Router:

    /container/print

    /container/start 0

0 is your container ID, please see the list of containers you got from the <i>/container/print</i> command.

===============================================================

<b>How to use:</b>

Now access speedtest from your web browser  (IP Router:Port).

You can go to http://192.168.1.2:8080 to do download test, go to http://192.168.1.2:8080/condition to check for threshold based download speed condition or http://192.168.1.2:8080/chart to check for speedtest history chart.

Finally, you can now get your actual download speedtest (in Megabits per second) using MikroTik Terminal and do some scripting like speedtest based failover.

For example:

    :local result [tool fetch mode=http url="http://192.168.1.2:8080" output=user as-value]
    :local result [:put ($result->"data")]
    :log info "Your actual download bandwidth is $result Mbps"
the example output from the script above is: Your actual download bandwidth is 3.14 Mbps

or

    :local result [tool fetch mode=http url="http://192.168.1.2:8080/condition" output=user as-value]
    :local result [:put ($result->"data")]
    :log info "Your actual download bandwidth currently $result"
the example output from the script above is: Your actual download bandwidth is currently Good

You can also use MikroTik Netwatch for easier scripting like this:

![](https://raw.githubusercontent.com/cdhtlr/MikroTik-Speedtest/main/Netwatch.png "Netwatch")

	Response code 200 means the current download bandwidth is Good and the network status will be "up"
	Response code 201 means the current download bandwidth is Bad and the network status will be "down"

<b>For performance and speedtest accuracy:</b>

Larger MAX_DLSIZE can give better download speedtest results, but MAX_DLSIZE setting that is too large can cause MikroTik to fail to execute scripts due to timeout.

The command line application in this container is made using Golang which is well known for its performance but it is more difficult to do manual memory management.

If the memory usage in the container continues to grow and you are not comfortable with this, you can set the memory limit on the container.

Memory limit that is too small can reduce CPU performance. So please set the memory limit wisely.

<b>Copyright notice:</b>

The command-line speedtest used in this Docker Image is modified from <a href="https://github.com/raviraa/speedtest">Raviraa Speedtest</a> and the chart used in this Docker Image is made by <a href="https://github.com/go-echarts/go-echarts">go-echarts</a>.
