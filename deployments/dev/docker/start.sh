#!/bin/sh

#mkdir -p /data/session && chmod -R 777 /data/session
#mkdir -p /data/cache && chmod -R 777 /data/cache
#mkdir -p /data/log && chmod -R 777 /data/log
#mkdir -p /data/upload && chmod -R 777 /data/upload

chown -Rf nginx.nginx /var/www/html

# 清除给定连接组件的数据库表结构缓存
#/var/www/html/yii cache/flush-schema db --interactive=0

# Start supervisord and services
/usr/bin/supervisord -n -c /etc/supervisord.conf
