FROM phusion/baseimage:latest

ENV HOME /root
ENV DEBIAN_FRONTEND nonintaractive

CMD ["/sbin/my_init", "--enable-insecure-key"]

ENV DEBIAN_FRONTEND nonintaractive
RUN apt-get update -y -qq
RUN apt-get install -y -qq curl

RUN apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*