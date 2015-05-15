FROM concourse/busyboxplus:git

# satisfy go crypto/x509
RUN cat /etc/ssl/certs/*.pem > /etc/ssl/certs/ca-certificates.crt

ADD built-check /opt/resource/check
ADD built-in /opt/resource/in
ADD built-out /opt/resource/out
