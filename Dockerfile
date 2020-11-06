FROM alpine:3.12.1

ENV TILE38_VERSION 1.22.3
ENV TILE38_DOWNLOAD_URL https://github.com/tidwall/tile38/releases/download/$TILE38_VERSION/tile38-$TILE38_VERSION-linux-amd64.tar.gz

ENV TILE38_WEBSERVER_VERSION 0.1.1
ENV TILE38_WEBSERVER_DOWNLOAD_URL https://github.com/lucascalion/tile38-webserver/releases/download/$TILE38_WEBSERVER_VERSION/tile38-webserver-$TILE38_WEBSERVER_VERSION-linux-amd64.tar.gz
RUN addgroup -S tile38 && adduser -S -G tile38 tile38

RUN apk update \
    && apk add ca-certificates \
    && update-ca-certificates \
    && apk add openssl bash \
    && wget -O tile38.tar.gz "$TILE38_DOWNLOAD_URL" \
    && tar -xzvf tile38.tar.gz \
    && rm -f tile38.tar.gz \
    && mv tile38-$TILE38_VERSION-linux-amd64/tile38-server /usr/local/bin \
    && mv tile38-$TILE38_VERSION-linux-amd64/tile38-cli /usr/local/bin \
    && rm -fR tile38-$TILE38_VERSION-linux-amd64 \
    && wget -O tile38-webserver.tar.gz "$TILE38_WEBSERVER_DOWNLOAD_URL" \
    && tar -xzvf tile38-webserver.tar.gz \
    && rm -f tile38-webserver.tar.gz \
    && mv tile38-webserver-$TILE38_WEBSERVER_VERSION-linux-amd64/tile38-webserver /usr/local/bin \
    && mv tile38-webserver-$TILE38_WEBSERVER_VERSION-linux-amd64/start_server /usr/local/bin \
    && rm -fR tile38-webserver-$TILE38_WEBSERVER_VERSION-linux-amd64

RUN mkdir /data && chown tile38:tile38 /data

VOLUME /data
WORKDIR /data

#Note that although only tile38 default port is being exposed here
#the webser port defined at the .env file will also be available and must be mapped accordingly
EXPOSE 9851
CMD ["start_server"]
