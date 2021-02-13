FROM alpine:3.6

RUN echo "http://mirrors.ustc.edu.cn/alpine/v3.6/main" > /etc/apk/repositories \
    && apk --update add ca-certificates tzdata \
    && rm -f /var/cache/apk/*

ENV ENV="/etc/profile"
WORKDIR /application

COPY bin/ /application/
RUN chmod +x /application/ark-robot

STOPSIGNAL SIGTERM
CMD /application/ark-robot -logtostderr=true -v=4