FROM ubuntu:18.04 as build-env

WORKDIR /
RUN apt-get update && \
    apt-get install --no-install-recommends --no-install-suggests -y \
                ca-certificates \
                git \
                wget && \
    wget https://golang.org/dl/go1.15.6.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.15.6.linux-amd64.tar.gz

RUN git clone https://github.com/rnburn/authz-service && \
    cd authz-service && \
    /usr/local/go/bin/go build

FROM gcr.io/distroless/base
ADD config.yaml /
COPY --from=build-env /authz-service/authz-service /
ENV HT_CONFIG_FILE=/config.yaml
CMD ["/authz-service"]
