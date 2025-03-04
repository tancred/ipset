# Build with
#
#   docker build --file Dockerfile --tag golang-ipset:1.23 .
#
# Run with the NET_ADMIN capability set and CWD mounted at /work with
#
#   docker run -it --rm --cap-add NET_ADMIN -v "$PWD":/work -w /work golang-ipset:1.23 bash
#
# Happy coding!
#
FROM golang:1.24

RUN export DEBIAN_FRONTEND=noninteractive DEBCONF_NONINTERACTIVE_SEEN=true \
    && apt-get -q update \
    && apt-get -q dist-upgrade -y \
    && apt-get install -y \
      libipset-dev \
      ipset \
    && rm -r /var/lib/apt/lists/*
