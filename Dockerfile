FROM centos:7.1.1503
MAINTAINER Alok Kumar Singh "mail.alok87@gmail.com"

# Install basic utilities
RUN yum -y update \
    && yum -y install epel-release \
    && yum -y install gcc make git tar mariadb-devel libffi-devel openssl-devel \
    && yum -y clean all 

# Install Go
RUN yum install -y curl
RUN curl -s https://storage.googleapis.com/golang/go1.4.2.linux-amd64.tar.gz | tar -v -C /usr/local/ -xz

# Set Environment 
ENV PATH /usr/local/go/bin:/usr/local/bin:/usr/local/sbin:/usr/bin:/usr/sbin:/bin:/sbin
ENV GOPATH /go
ENV GOROOT /usr/local/go

# Get your project
RUN mkdir -p /go/src/github.com/alok87/sysAdmin
RUN mkdir -p /go/{bin,pkg}
WORKDIR /go/src/github.com/alok87/sysAdmin
ADD . /go/src/github.com/alok87/sysAdmin
RUN go get -u github.com/gorilla/websocket

# Install application
WORKDIR /go/src/github.com/alok87/sysAdmin/src/main
RUN yum install -y sudo
RUN go install

# Expose application port
WORKDIR /go/bin
EXPOSE 3500 

# Entrypoint, the command which will run when application launches
#ENTRYPOINT ./main
