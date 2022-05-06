#!/bin/sh

nohup php8 -S 0.0.0.0:80 -t /var/www/localhost/htdocs/ > /dev/null 2>&1 &

[ ! -z "$TZ" ] && echo 'Setting up the time zone'
[ ! -z "$TZ" ] && ln -sf /usr/share/zoneinfo/$TZ /etc/localtime

echo 'Container date and time:'
date

[ ! -z "$CRON_FIELD" ] && echo 'Setting up the cron'
[ ! -z "$CRON_FIELD" ] && echo "$CRON_FIELD php8 /var/www/localhost/htdocs/speedtest.php" | crontab -
[ ! -z "$CRON_FIELD" ] && crond -b > /dev/null 2>&1

[ ! -z "$PING_IP" ] && echo 'Checking internet connection to public IP'
[ ! -z "$PING_IP" ] && ping -c 4 $PING_IP | grep icmp* | wc -l > /dev/null 2>&1
[ ! -z "$PING_DOMAIN" ] && echo 'Checking internet connection to public domain'
[ ! -z "$PING_DOMAIN" ] && ping -c 4 $PING_DOMAIN | grep icmp* | wc -l > /dev/null 2>&1

echo 'Active public IPv4 addresses:'
curl -s -4 https://ipv4.icanhazip.com

echo 'Active public IPv6 addresses:'
curl -s -6 https://ipv6.icanhazip.com

echo 'Active local IP addresses:'
ip address | grep inet

echo 'Listening ports:'
netstat -an | grep LISTEN

echo 'You need to be connected to the internet,'
echo 'make sure your date and time settings are correct'
echo 'and ensure that there are no problems'
echo 'with your DNS and routing from container to the internet.'

echo 'Enjoy speedtest :)'

exec top > /dev/null 2>&1
