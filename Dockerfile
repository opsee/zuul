FROM quay.io/opsee/vinz

MAINTAINER Greg Poirier <greg@opsee.co>

USER zuul

ENV KEY_ALIAS=""
ENV KEY_BUCKET=""
ENV SERVER_PRIVATE_KEY_OBJECT=""
ENV SERVER_PUBLIC_KEY_OBJECT=""
ENV CLIENT_PRIVATE_KEY_OBJECT=""
ENV CLIENT_PUBLIC_KEY_OBJECT=""
ENV AWS_DEFAULT_REGION=""
ENV AWS_SECRET_ACCESS_KEY=""
ENV AWS_ACCESS_KEY_ID=""

ENV COMMON="/zuul/common"
ENV STATE="/zuul/state"
ENV PATH="/zuul/bin:/bin:/sbin:/usr/bin:/usr/sbin:/opt/bin"

RUN apk --update add openssh bash && \
    mkdir -p /zuul/bin && \
    mkdir -p /zuul/socket && \
    ln -sf /zuul/connector/connector.sh /zuul/bin/connector && \
    ln -sf /zuul/multiplexer/multiplexer.sh /zuul/bin/multiplexer && \
    ln -sf /zuul/server/server.sh /zuul/bin/server && \
    adduser --disabled-password --gecos '' --home /zuul zuul && \
    chown -R zuul:zuul /zuul

COPY . /zuul

VOLUME /zuul/socket

CMD ["/bin/bash"]
