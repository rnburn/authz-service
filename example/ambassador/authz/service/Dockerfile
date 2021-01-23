FROM ubuntu:18.04


RUN apt-get update && \
    apt-get install --no-install-recommends --no-install-suggests -y \
                ca-certificates \
                git \
                wget && \
    wget https://golang.org/dl/go1.15.6.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.15.6.linux-amd64.tar.gz
ADD config.yaml /
ENV HT_CONFIG_FILE=/config.yaml
RUN /usr/local/go/bin/go get github.com/rnburn/authz-service
CMD ["/root/go/bin/authz-service"]
