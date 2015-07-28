System Administrator App
=========================
This is an open source application for performing regular system admin tasks.

Installation
=============
1. Docker Container running this applciation to manage itself. 
 
 * Install docker in the machine where this application will run in a container.
 
   	[Install Docker](https://docs.docker.com/installation/)
	Centos7:	
	```sudo yum install -y docker```
	```sudo systemctl start docker```

 * Clone the source from github.
 
   	```git clone https://github.com/alok87/sysAdmin```

 * Build the Dockerfile, to create an image.
 	
	```cd sysAdmin; sudo docker build -t alok87/centos01 . ```
	
 * Run the docker image to spawn the container.
   
	```sudo docker run -i -t -p 3500:3500 --name centos01 alok87/centos01 /bin/bash```
	```/go/bin/main```
 
 * Visit ```http://localhost:3500/users``` from your docker host machine.
 

2. Installation in Physical machine/VM to manage itself.

 * Setup the password less sudo access for the user who will run the application. 
 
   	```sudo visudo```	or open the ```/etc/sudoerrs``` file (it wil prompt for the password)

   	Add the below lines in it (user here is aks for example).
   	
	```User_Alias SUDOUSERS = aks```
	```SUDOUSERS       ALL = (ALL) NOPASSWD: ALL```

	Uncomment the wheel line to have wheel group password less sudo access(optional)

		 ```%wheel  ALL=(ALL)       NOPASSWD: ALL```

 
 * Disable requiretty 
 
   	```sudo visudo or open the /etc/sudoerrs file```
   	
	Comment the "Defaults requiretty" line

 * Set GOPATH, GOROOT and PATH (put it inside .bashrc of the user that will run the application)

	```export GOROOT=/home/vic/code/golang/go```

	```export GOPATH=/home/vic/code/golang/workspace```

	```export PATH=$PATH:$GOROOT/bin:$GOPATH/bin```	

 * Download go and put it in GOROOT directory
	
	[Install Golang](https://golang.org/doc/install)

 * Clone the application's repository from github inside $GOPATH/src/github.com/alok87
 
	```mkdir -p $GOPATH/src/github.com/alok87```

	```git clone https://github.com/alok87/sysAdmin```


 * Install Gorilla's Websocket library and verify it got installed under $GOPATH/bin/ 
  
 	```go get -u github.com/gorilla/websocket```
  
 * Build the code and run it.
 
	 ```
	cd $GOPATH/src/github.com/alok87/sysAdmin/main/
  	
	go install
   
	./main
	```
	
 * Visit ```http://localhost:3500/users```
