FROM ubuntu:14.04
MAINTAINER Paul Bellamy <paul@paulbellamy.com>

RUN apt-get update
RUN apt-get install -y build-essential pkg-config git-core python-software-properties python libssl-dev libpcap-dev libboost-all-dev libcrypto++-dev libsqlite3-dev
RUN mkdir /nfd
RUN cd /nfd
RUN git clone --depth 1 git://github.com/named-data/ndn-cxx ndn-cxx
RUN cd ndn-cxx && ./waf configure && ./waf -j1 && ./waf install && ldconfig
RUN ndnsec-keygen /tmp/key | ndnsec-install-cert -
RUN git clone git://github.com/named-data/NFD NFD
RUN cd NFD && git submodule init && git submodule update
RUN cd NFD && ./waf configure --with-tests && ./waf -j1 && ./waf install && cd ..
RUN cp /usr/local/etc/ndn/nfd.conf.sample /usr/local/etc/ndn/nfd.conf

# Get nfd-start-fg script to boot it in the foreground
RUN apt-get install -y wget
RUN wget https://raw.githubusercontent.com/felixrabe/NFD/ndn-start-fg/tools/nfd-start-fg.sh -O /usr/local/bin/nfd-start-fg
RUN sed -i'' -e 's|@BASH@|/bin/bash|' /usr/local/bin/nfd-start-fg
RUN sed -i'' -e 's|@VERSION@|master|' /usr/local/bin/nfd-start-fg
RUN chmod +x /usr/local/bin/nfd-start-fg

EXPOSE 6363

CMD nfd-start-fg
