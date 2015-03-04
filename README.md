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

Install MongoDB
Download Flanders
Extract Flanders
Run Flanders


## Development Setup

### Prerequisites

- [NodeJS](http://nodejs.org) (for user interface and build tools)
- [VirtualBox](https://www.virtualbox.org/)
- [Vagrant](http://vagrantup.com) - Spins up a virtual machine with all prerequisites ready to go

### Instructions

Checkout the code base

```
$ git clone https://github.com/weave-lab/flanders.git
```

Go into the directory

```
$ cd flanders
```

#### Server

Spin up the development virtual machine

```
$ vagrant up flanders
```

SSH into the dev VM

```
$ vagrant ssh flanders
```

Inside the VM, run the app

```
vagrant$ cd /opt/go/src/github.com/weave-lab/flanders
vagrant$ go run main/main.go
```

When you change the Go code in flanders, it is automatically synced to your virtual machine, so you just have to restart your app to see changes.

#### User Interface

In a different terminal window, change to the web directory

```
$ cd web
```

Install all the front-end dependencies

```
$ npm install
$ npm install -g bower
$ bower install
$ npm install -g grunt-cli
```

Start the dev server for user interface

```
$ grunt serve
```

A browser window will popup and will show the flanders ui and will automatically connect to the flanders service running in the virtual machine.
If you make changes to the front-end code, grunt will automatically update your browser window


