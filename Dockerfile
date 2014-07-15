# Dockerfile to run main.go 
# VERSION 2 - EDITION 1

# Base image for sandbox enviroment

FROM mike/ubuntu13.04:update

#Maintainer: ggaaooppeenngg <gaopeg01 at gmail com>

MAINTAINER ggaaooppeenngg,gaopeg01@gmail.com


#install newest golang
ADD  go1.3.linux-amd64.tar.gz  /home/
RUN  mkdir  /home/GoPath
ENV  GOPATH /home/GoPath
ENV  GOROOT /home/go
ENV  PATH   /home/go/bin:$PATH
RUN  mkdir  /usr/local/go

#RUN  cp -r /home/go/* /usr/local/go

RUN  apt-get install -y --force-yes git
RUN  apt-get install -y --force-yes  gcc
RUN  apt-get install -y --force-yes  g++
RUN  go get -v github.com/codegangsta/cli
RUN  go get -v github.com/ggaaooppeenngg/sandbox
RUN  go install github.com/ggaaooppeenngg/sandbox/sandbox


