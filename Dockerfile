FROM phusion/baseimage:latest

ENV HOME /root
ENV DEBIAN_FRONTEND nonintaractive

CMD ["/sbin/my_init", "--enable-insecure-key"]

ENV DEBIAN_FRONTEND nonintaractive
RUN apt-get update -y -qq
RUN apt-get install -y -qq curl

RUN echo 'deb http://files.freeswitch.org/repo/deb/debian/ wheezy main' >> /etc/apt/sources.list.d/freeswitch.list
RUN curl -s http://files.freeswitch.org/repo/deb/debian/freeswitch_archive_g0.pub | apt-key add -
RUN apt-get update -y -qq

RUN apt-get install -y -qq freeswitch-meta-vanilla
RUN cp -a /usr/share/freeswitch/conf/vanilla /etc/freeswitch


RUN apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

ENV DEBIAN_FRONTEND dialog
RUN apt-get clean