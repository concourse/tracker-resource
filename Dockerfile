FROM gliderlabs/alpine

RUN apk-install ca-certificates
RUN apk-install git

ADD built-check /opt/resource/check
ADD built-out /opt/resource/out
