# Dockerfile to run main.go 
# VERSION 2 - EDITION 1

# Base image used is ubuntu

FROM mike/ubuntu13.04:update

#Maintainer: ggaaooppeenngg <gaopeg01 at gmail com>

MAINTAINER ggaaooppeenngg,gaopeg01@gmail.com


#install newest golang
ADD go1.3.linux-amd64.tar.gz  /home/
RUN  mkdir /home/GoPath
ENV  GOPATH /home/GoPath
ENV  GOROOT /home/go
RUN  mkdir /usr/local/go

RUN  cp -r /home/go/* /usr/local/go

RUN  apt-get install -y --force-yes git
RUN  apt-get install -y gcc
RUN  apt-get install -y g++
RUN  /home/go/bin/go get github.com/ggaaooppeenngg/sandbox


