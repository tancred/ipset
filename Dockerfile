# Build with
#
#   docker build --file Dockerfile --tag golang-ipset:1.23 .
#
# Run with CWD mounted at /work with
#
#   docker run -it --rm -v "$PWD":/work -w /work golang-ipset:1.23 bash
#
# Happy coding!
#
FROM golang:1.23

RUN export DEBIAN_FRONTEND=noninteractive DEBCONF_NONINTERACTIVE_SEEN=true \
    && apt-get -q update \
    && apt-get -q dist-upgrade -y \
    && apt-get install -y \
      libipset-dev \
      ipset \
    && rm -r /var/lib/apt/lists/*
