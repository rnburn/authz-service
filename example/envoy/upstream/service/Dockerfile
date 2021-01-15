FROM python:3-alpine


RUN apk update && apk add build-base
RUN pip3 install -q Flask==1.1.2 opentelemetry-api==0.16b1 \
                                  opentelemetry-sdk==0.16b1 \
                                  opentelemetry-exporter-zipkin==0.16b1 \
                                  opentelemetry-instrumentation-flask==0.16b1
COPY . ./app
CMD ["python3", "/app/service/server.py"]
