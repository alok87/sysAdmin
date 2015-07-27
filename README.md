System Administrator App
=========================
This is an open source application for performing regular system admin tasks.

Installation
=============
1. Running the application in a container using docker to manage the container.
 
 * Install docker in the machine where this application will run in a container.
 
   	[Install Golang](https://docs.docker.com/installation/)

 * Clone the application's repository from github.

   	```git clone https://github.com/alok87/sysadminApp```

 * Build the Dockerfile, to create an image.
 	```
   	cd sysadminApp
   	docker build -t <yourusername>/sysadminApp .
   	eg: docker build -t alok87/sysadminApp
	```
 * Run the docker image to spwan the container.
   
	```docker run --name sysadmin01 alok87/sysadminApp```
 
 * Visit http://localhost:8005/users
 

2. Running the application by manually setting up the machine to manage the machine.

 * Setup the password less sudo access for the user who will run the application. 
 
   	```sudo visudo```	or open the ```/etc/sudoerrs``` file (it wil prompt for the password please give it)

   	Add the below lines in it. (user here is aks for example )
   	
		## Allow root to run any commands anywhere
		
		User_Alias SUDOUSERS = aks
		
		SUDOUSERS       ALL = (ALL) NOPASSWD: ALL

	Uncomment the wheel line to have wheel group password less sudo access
	## Same thing without a password
	%wheel  ALL=(ALL)       NOPASSWD: ALL

 
 * Disable requiretty 
 
   	```sudo visudo or open the /etc/sudoerrs file```
   	
	Comment the "Defaults requiretty" line

 * Clone the application's repository from github.
 
   	```git clone https://github.com/alok87/sysadminApp```

 * Install go language in your system

 * Set GOPATH and GOROOT

 * Add the current application in your GOPATH

 * Install additional libraries 
  
 	```go get -u github.com/gorilla/websocket```
  
 * Go to the sysadminApp directory and build the go project.
 
	 ```
	cd sysadminApp
  
	go build -o sysAdmin 
   
	./sysAdmin
	```
	
 * Visit http://localhost:8005/users
