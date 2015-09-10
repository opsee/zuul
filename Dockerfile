FROM gliderlabs/alpine:3.2
MAINTAINER Greg Poirier <greg@opsee.co>

ENV PATH=/zuul/bin:/opt/bin:/bin:/usr/bin:/usr/local/bin:/sbin:/usr/sbin:/usr/local/sbin

RUN mkdir -p /zuul/bin /opt/bin && \
    apk add --update ca-certificates curl bash && \
    curl -Lo /opt/bin/ec2-env https://s3-us-west-2.amazonaws.com/opsee-releases/go/ec2-env/ec2-env && \
    chmod 755 /opt/bin/ec2-env

COPY register.sh /zuul/bin/register.sh
COPY bin/ /zuul/bin/
