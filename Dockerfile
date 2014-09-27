FROM google/golang:latest

RUN go get -u github.com/tools/godep

RUN export PATH=$PATH:/gopath/bin:/goroot/bin

#ADD . $GOPATH/src/lab.getweave.com/weave/flanders
#RUN godep get lab.getweave.com/weave/flanders
#RUN cd $GOPATH/src/lab.getweave.com/weave/flanders && godep restore

