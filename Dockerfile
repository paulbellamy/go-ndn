FROM ubuntu
MAINTAINER Paul Bellamy <paul@paulbellamy.com>

RUN apt-get update
RUN apt-get install -y build-essential pkg-config git-core python-software-properties python libssl-dev libpcap-dev libboost-all-dev libcrypto++-dev libsqlite3-dev
RUN mkdir /nfd
RUN cd /nfd
RUN git clone --depth 1 git://github.com/named-data/ndn-cxx ndn-cxx
RUN cd ndn-cxx && ./waf configure && ./waf -j1 && ./waf install && ldconfig && cd ..
RUN ndnsec-keygen /tmp/key | ndnsec-install-cert -
RUN git clone --depth 1 git://github.com/named-data/NFD NFD
RUN cd NFD && ./waf configure --with-tests && ./waf -j1 && ./waf install && cd ..

# use supervisor to start the two processes, because they tell you to have it
# self-daemonize.
RUN apt-get install -y supervisor
RUN mkdir -p /var/log/supervisor
RUN echo "[supervisord]\nnodaemon=true" > /etc/supervisor/conf.d/supervisor.conf
RUN echo "[program:nfd]\ncommand=nfd" > /etc/supervisor/conf.d/nfd.conf
RUN echo "[program:nrd]\ncommand=nrd" > /etc/supervisor/conf.d/nrd.conf
CMD ["/usr/bin/supervisord"]
