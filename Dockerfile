FROM concourse/busyboxplus:git

ADD built-check /opt/resource/check
ADD built-out /opt/resource/out
