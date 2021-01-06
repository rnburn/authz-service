#!/bin/bash

set -e
apt-get update 
apt-get install --no-install-recommends --no-install-suggests -y \
                build-essential \
                ca-certificates \
                python python3 python3-distutils python3-dev \
                git \
                wget
curl https://bootstrap.pypa.io/get-pip.py -o get-pip.py
python3 get-pip.py
