#!/bin/bash

set -e
apt-get update 
apt-get install --no-install-recommends --no-install-suggests -y \
                build-essential \
                apt-transport-https \
                software-properties-common \
                ca-certificates \
                gnupg-agent \
                python python3 python3-distutils python3-dev \
                git \
                curl \
                wget

# Python
curl https://bootstrap.pypa.io/get-pip.py -o get-pip.py
python3 get-pip.py
pip3 install Flask==1.1.2

# Go
wget https://golang.org/dl/go1.15.6.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.15.6.linux-amd64.tar.gz

# Envoy
curl -sL 'https://getenvoy.io/gpg' | apt-key add -
add-apt-repository \
    "deb [arch=amd64] https://dl.bintray.com/tetrate/getenvoy-deb \
     $(lsb_release -cs) \
     stable"
apt-get update && apt-get install -y getenvoy-envoy
