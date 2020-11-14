FROM ubuntu:focal

ADD https://pkg.tarsnap.com/tarsnap-deb-packaging-key.asc /
RUN apt-get update && apt-get install -y gnupg2 && apt-key add tarsnap-deb-packaging-key.asc
RUN echo "deb http://pkg.tarsnap.com/deb/focal ./" | tee -a /etc/apt/sources.list.d/tarsnap.list
RUN apt-get update && apt-get install -y tarsnap-archive-keyring
RUN echo "deb-src http://pkg.tarsnap.com/deb-src ./" | tee -a /etc/apt/sources.list.d/tarsnap.list
RUN apt-get update && gpg --no-default-keyring --keyring trustedkeys.gpg --import /usr/share/keyrings/tarsnap-archive-keyring.gpg
RUN env DEBIAN_FRONTEND="noninteractive" apt-get build-dep -y tarsnap && cd $(mktemp -d) && apt-get source -b tarsnap && dpkg -i tarsnap_*.deb
