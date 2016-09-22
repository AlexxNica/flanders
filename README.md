# Flanders

A sip-capture server written in Go.

## Inspiration

The open source project [Homer](http://www.sipcapture.org/) is a great tool for your VoIP arsenal. I would say it is necessary for easy
diagnosing of SIP related issues in your VoIP stack. Homer has saved the day for me many times over when trying to dianose issues. 
Flanders is being designed to be a drop in replacement for Homer with some different goals in mind:

## Goals

- Easy Installation - We bundled the sip capture server into the app for one single binary to install
- Clean and modern UI - We programmed the user interface as a nice single-page angular app
- Improved Data Store [up for debate :-)] - We opted for MongoDB as the default storage engine because of its ability to handle so many inserts out of the box, and built-in map reduce functions for complex queries
- Sharing Call History - Call details have unique urls for easy sharing with co-workers. No popup hell.
- Real time SIP packet filters - THIS IS AWESOME! We want to be able to see calls progress in real time based on filters. Screw you ngrep...

This project is super young and isn't even close to production ready, and doesn't have nearly the features of Homer... YET. It is actively being developed here at [Weave](http://getweave.com) and so expect big changes and more stability soon.

## Installation 

- Install MongoDB
- Download Flanders
- Extract Flanders
- Run Flanders

