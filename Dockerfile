FROM tanglicai.xyz:5555/crhome-alpine:1.0

ENV ENV="/etc/profile"
WORKDIR /application

COPY docs/* /docs/
COPY bin/ /application/
COPY resource/* /resource/
RUN chmod +x /application/ark-robot

STOPSIGNAL SIGTERM
CMD /application/ark-robot -logtostderr=true -v=4