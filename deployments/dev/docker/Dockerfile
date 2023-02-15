FROM alpine:3.10

MAINTAINER stccon <1013124327@qq.com>

#换源 安装软件
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories && \
apk update && apk upgrade && \
apk add --no-cache \
php7 \
php7-pdo_sqlite \
php7-bcmath \
php7-common \
php7-ctype \
php7-curl \
php7-dom \
php7-exif \
php7-fpm \
php7-gd \
php7-iconv \
php7-fileinfo \
php7-intl \
php7-json \
php7-mysqli \
php7-mysqlnd \
php7-mbstring \
php7-mcrypt \
php7-openssl \
php7-opcache \
php7-pdo \
php7-pdo_mysql \
php7-pdo_pgsql \
php7-pgsql \
php7-phar \
php7-redis \
php7-session \
php7-simplexml \
php7-soap \
php7-sqlite3 \
php7-tokenizer \
php7-xdebug \
php7-xml \
php7-xmlrpc \
php7-xmlreader \
php7-xmlwriter \
php7-xsl \
php7-zlib \
php7-zip \
composer \
curl \
nginx \
git \
htop \
openssl \
supervisor \
tzdata \
wget \
yasm \
ffmpeg \
exiftool \
go musl-dev gcc build-base unrar \
vim && \
cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
apk del tzdata && \
mkdir -p /var/www/html

#复制n配置文件
ADD etc/ /etc/

RUN touch /var/run/nginx.pid

ADD start.sh /start.sh
RUN chmod 755 /start.sh && dos2unix /start.sh

VOLUME /var/www/html
WORKDIR /var/www/html

EXPOSE 80

CMD ["/start.sh"]
