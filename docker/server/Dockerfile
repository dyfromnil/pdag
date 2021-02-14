FROM alpine
RUN apk update \
    && apk add iputils \
    && apk add iperf
COPY ./main /home/dy/server/