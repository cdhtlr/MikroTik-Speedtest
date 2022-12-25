<b>DockerHub :</b> https://hub.docker.com/r/cdhtlr/mikrotik-speedtest

![](https://raw.githubusercontent.com/cdhtlr/MikroTik-Speedtest/main/Demo.png "Demo")

<b>System requirements:</b>
- <b>Any linux PC</b> or <b>MikroTik RouterOS v7.1rc3 (and above)</b> (Currently support Container for <b>ARM</b>, <b>ARM64</b>, <b>X86</b> and <b>X86-64</b> based machine, like the <i>CCR2004-16G-2S+</i> device with ARM64 architecture)
- <b>27MB of disk space left</b> (9MB for compressed TAR image, 18MB for uncompressed running image. You can delete the compressed TAR image once the image is successfully decompressed)
- <b>About 6MB to 12MB of RAM space left</b> (Could be more, depending on the workload)

===============================================================

<b>For use with PC (outside MikroTik):</b>

    docker run --restart=unless-stopped \ 
    --name speedtest -d \ #set name "speedtest" and Run container in background and print container ID
    -p 8080:80 \ #Host-port:Container-port to access the app inside this container via port 8080
    -e 'MAX_KB=1000' \ #maximum size in KB to download (optional)
    -e 'THRESHOLD_MBPS=1.0' \ #download threshold in Mbps, to check for download speed condition (optional)
    -e 'URL=https://jakarta.speedtest.telkom.net.id.prod.hosts.ooklaserver.net:8080/download?size=25000000' \ #url to download (optional)
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
    /container envs add list=speedtest name=MAX_KB value="1000"
    /container envs add list=speedtest name=THRESHOLD_MBPS value="1.0"
    /container envs add list=speedtest name=URL value="https://jakarta.speedtest.telkom.net.id.prod.hosts.ooklaserver.net:8080/download?size=25000000"
    /container add file=speedtest.tar interface=veth1-speedtest envlist=speedtest hostname=speedtest logging=yes

Check your container list in MikroTik Router:

    /container/print

    /container/start 0

0 is your container ID, please see the list of containers you got from the <i>/container/print</i> command.

===============================================================

<b>How to use:</b>

Now access speedtest from your web browser  (IP Router:Port).

For example http://192.168.1.2:8080 (<i>this is the main page containing the Speedtest Graph</i>)

You can go to http://192.168.1.2:8080/test to do download speedtest or go to http://192.168.1.2:8080/condition to check for threshold based download speed condition

Finally, you can now get your actual download speedtest (in Megabits per second) using MikroTik Terminal and do some scripting.

For example:

    :local result [tool fetch mode=http url="http://192.168.1.2:8080/test" output=user as-value]
    :local result [:put ($result->"data")]
    :log info "Your actual download bandwidth is $result Mbps"
the example output from the script above is: Your actual download bandwidth is 3.14 Mbps

or

    :local result [tool fetch mode=http url="http://192.168.1.2:8080/condition" output=user as-value]
    :local result [:put ($result->"data")]
    :log info "Your actual download bandwidth currently $result"
the example output from the script above is: Your actual download bandwidth is currently Good

Larger MAX_KB can give better download speedtest result, but MAX_KB settings that are too large can cause MikroTik to fail to execute scripts due to timeouts.

<p><b>Copyright notice:</b><br>The command-line speedtest used in this Docker Image is made by <a href="https://github.com/raviraa/speedtest">Raviraa Speedtest</a> and the chart graph used in this Docker Image is made by <a href="https://github.com/go-echarts/go-echarts">go-echarts</a>.</p>
