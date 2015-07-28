FROM ubuntu:12.04
MAINTAINER Alok Kumar Singh "mail.alok87@gmail.com"

# Needed for "go get" to work - Mercurial
RUN echo 'deb http://ppa.launchpad.net/mercurial-ppa/releases/ubuntu precise main' > /etc/apt/sources.list.d/mercurial.list
RUN echo 'deb-src http://ppa.launchpad.net/mercurial-ppa/releases/ubuntu precise main' >> /etc/apt/sources.list.d/mercurial.list
RUN apt-key adv --keyserver keyserver.ubuntu.com --recv-keys 323293EE

# Install Go
RUN apt-get update
RUN apt-get install -y curl git bzr mercurial
RUN curl -s https://storage.googleapis.com/golang/go1.4.2.linux-amd64.tar.gz | tar -v -C /usr/local/ -xz

# Set Environment 
ENV PATH /usr/local/go/bin:/usr/local/bin:/usr/local/sbin:/usr/bin:/usr/sbin:/bin:/sbin
ENV GOPATH /go
ENV GOROOT /usr/local/go

# Get your project
RUN go get github.com/alok87/sysAdmin
WORKDIR /go/src/github.com/alok87/sysAdmin
ADD . /go/src/github.com/alok87/sysAdmin
RUN go get -u github.com/gorilla/websocket

# Install application
RUN cd /go/src/github.com/alok87/sysAdmin/main
RUN go install

# Expose application port
EXPOSE 8005

# Entrypoint, the command which will run when application launches
ENTRYPOINT /go/bin/main
