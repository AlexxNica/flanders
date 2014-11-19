# Flanders

A sip-capture server written in Golang.

## Inspiration

The open source project [Homer](http://www.sipcapture.org/) is a great tool for your VoIP arsenal. I would say almost necessary for easy
diagnosing of SIP related issues in your VoIP stack. Homer has saved the day for me many times over when trying to dianose issues. 
But Homer lacks in certain areas:

- UI - the whole user interface is hideous, with areas for drag and drop that should never be drag and drop. It was obviously created by a systems engineer like most open source VoIP tools.
- Requires Sip Capture Server - In most cases, this is [Kamailio](http://www.kamailio.org/). This means installing 2 tools, the capture server, and Homer, which is just a web interface over a mysql database. (Homer does now provide its own sip capture server as well, but is still a separate install)
- Slow - It's default MySQL setup can get super slow super fast for large call volume clusters. It comes to a screetching halt and crashes sometimes if you query the data the wrong way

## Goals

- Beautiful UI - With Flanders, we set out to make a beautiful UI, that is much more user friendly. 
- Bundled Sip Catpure Server - We also wanted to make it a single install with minimal config. It should be easy to setup.
- Improved Data Store (default) - Lasly, we opted to use a great time series database, [InfluxDB](http://influxdb.com/) which gives us some great features for storing SIP packet data
- HEP Compatible - We want this to be a drop in replacement for most Homer setups. Currently, we have only tested with FreeSWITCH and OpenSIP (HEPv1)
- Real time SIP packet filters - We want to be able to see calls progress in real time based on filters. Screw you ngrep...

This project is super young and isn't even close to production ready, and doesn't have nearly the features of Homer... YET. It is actively being developed here at Weave and so expect big changes and more stability soon.


![Flanders](web/app/images/stupid_sexy_flanders.jpg)
