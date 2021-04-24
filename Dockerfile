FROM alpine:3.6
#FROM alpine:latest

## Installs latest Chromium package.
#RUN echo "http://dl-cdn.alpinelinux.org/alpine/edge/main" > /etc/apk/repositories \
#    && echo "http://dl-cdn.alpinelinux.org/alpine/edge/community" >> /etc/apk/repositories \
#    && echo "http://dl-cdn.alpinelinux.org/alpine/edge/testing" >> /etc/apk/repositories \
#    && echo "http://dl-cdn.alpinelinux.org/alpine/v3.11/main" >> /etc/apk/repositories \
#    && apk upgrade -U -a \
#    && apk add --no-cache \
#    libstdc++ \
#    chromium \
#    harfbuzz \
#    nss \
#    freetype \
#    ttf-freefont \
#    wqy-zenhei \
#    && rm -rf /var/cache/* \
#    && mkdir /var/cache/apk
#
#ENV CHROME_BIN=/usr/bin/chromium-browser \
#    CHROME_PATH=/usr/lib/chromium/

RUN echo "http://mirrors.ustc.edu.cn/alpine/v3.6/main" > /etc/apk/repositories \
    && apk --update add ca-certificates tzdata \
    && rm -f /var/cache/apk/*

ENV ENV="/etc/profile"
WORKDIR /application

COPY docs/* /docs/
COPY bin/ /application/
COPY resource/* /resource/
RUN chmod +x /application/ark-robot

STOPSIGNAL SIGTERM
CMD /application/ark-robot -logtostderr=true -v=4