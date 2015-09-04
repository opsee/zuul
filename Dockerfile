FROM gliderlabs/alpine:3.2
MAINTAINER Greg Poirier <greg@opsee.co>

RUN mkdir -p /zuul/bin /opt/bin && \
    apk add --update ca-certificates curl && \
    curl -Lo /opt/bin/ec2-env https://s3-us-west-2.amazonaws.com/opsee-releases/go/ec2-env/ec2-env && \
    chmod 755 /opt/bin/ec2-env

COPY register.sh /zuul/bin/register.sh
COPY bin/ /zuul/bin/
