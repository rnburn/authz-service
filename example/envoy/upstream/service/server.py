from flask import Flask, request
import sys
from opentelemetry import trace
from opentelemetry.exporter import zipkin
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.sdk.trace.export import BatchExportSpanProcessor
from opentelemetry.instrumentation.flask import FlaskInstrumentor

trace.set_tracer_provider(TracerProvider())
tracer = trace.get_tracer(__name__)

# create a ZipkinSpanExporter
zipkin_exporter = zipkin.ZipkinSpanExporter(
    service_name="my-helloworld-service",
    url = "http://zipkin:9411/api/v2/spans",
)
span_processor = BatchExportSpanProcessor(zipkin_exporter)
trace.get_tracer_provider().add_span_processor(span_processor)

app = Flask(__name__)
FlaskInstrumentor().instrument_app(app)


@app.route('/service')
def hello():
  sys.stderr.write(str(request.headers))
  print(request.headers)
  return 'Hello from behind Envoy!'


if __name__ == "__main__":
  # zipkin = Zipkin(app, sample_rate=100)
  # app.config['ZIPKIN_DISABLE'] = False
  # app.config['ZIPKIN_DSN'] = "http://zipkin:9411/api/v1/spans"
  app.run(host='0.0.0.0', port=8080, debug=False)
