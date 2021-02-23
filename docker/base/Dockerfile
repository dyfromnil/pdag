FROM alpine
RUN apk update \
    && apk add iftop \
    && apk add iputils \
    && apk add iperf \
    && apk add iproute2 \
    && ln -s /usr/lib/tc /lib/tc