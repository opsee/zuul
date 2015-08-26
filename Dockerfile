FROM quay.io/opsee/vinz
MAINTAINER Greg Poirier <greg@opsee.co>

ENV PATH="/zuul/bin:/bin:/sbin:/usr/bin:/usr/sbin:/opt/bin"

RUN apk --update add openvpn bash curl && \
    curl -Lo /opt/bin/ec2-env https://s3-us-west-2.amazonaws.com/opsee-releases/go/ec2-env/ec2-env && \
    chmod 755 /opt/bin/ec2-env && \
    mkdir -p /zuul/bin && \
    mkdir -p /zuul/state && \
    ln -sf /zuul/client/client.sh /zuul/bin/client && \
    ln -sf /zuul/server/server.sh /zuul/bin/server && \
    ln -sf /zuul/gozer/auth.sh /zuul/bin/auth && \
    ln -sf /zuul/gozer/bin/router /zuul/bin/router && \
    openssl dhparam -out /zuul/state/dh1024.pem 1024 && \
    adduser -D -g '' -h /zuul -H -s /sbin/nologin zuul && \
    passwd -u zuul

COPY . /zuul

RUN chown -R zuul:zuul /zuul

CMD ["/bin/bash"]
