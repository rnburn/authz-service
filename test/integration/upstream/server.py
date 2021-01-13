from flask import Flask, request

app = Flask(__name__)


@app.route('/service')
def hello():
  print(request.headers)
  return 'Hello, Hello'


if __name__ == "__main__":
  app.run(host='0.0.0.0', port=8080, debug=False)
