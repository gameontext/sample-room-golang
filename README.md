# Microservices with a Game On! Room
[Game On!](https://game-on.org/) is both a sample microservices application, and a throwback text adventure brought to you by the wasdev team at IBM. It aims to do a few things:

- Provide a sample microservices-based application.
- Demonstrate how microservice architectures work from two points of view:
 - As a Player: Naviagate through a network/maze of rooms, where each room is a unique implementation of a common API. Each room supports chat, and interaction with items (some of which may be in the room, some of which may be separately defined services as well).
 - As a Developer: Learn about microservice architectures and supporting infrastructure by extending the game with your own services. Write additional rooms or items and see how they interact with the rest of the system.


##Introduction

This walkthrough will guide you through adding a room to a running Game On! server.  You will be shown how to setup a container based room that is implemented in the Go programming language.  There are also instructions for deploying a similar room written in Node.js that is deployed as a Cloud Foundry application in Bluemix.  Both of these rooms take a microservice approach to adding a room to a running Game On! text adventure game server.  

### Installation prerequisites

Game On! when deployed using containers requires:

- [Bluemix account](https://console.ng.bluemix.net)
- [GitHub account](https://github.com/)
- [Docker](https://docs.docker.com/engine/installation/) - You will need Docker on your host machine to create a docker image to be pushed to IBM Bluemix Containers
- [IBM Containers CLI](https://www.ng.bluemix.net/docs/containers/container_cli_ov.html#container_cli_cfic_install) - This is required to manage IBM Containers from your local machine.


## Create Bluemix accounts and log in
To build a Game On! room in Bluemix, you will first need a Bluemix account. 

### Sign up and log into Bluemix and DevOps
Sign up for Bluemix at https://console.ng.bluemix.net and DevOps Services at https://hub.jazz.net. When you sign up, you'll create an IBM ID, create an alias, and register with Bluemix.

## Get Game On! API Key and User ID
For a new room to register with the Game-On server, you must first log into game-on.org and sign in using one of several methods to get your game-on user ID and ApiKey.

1.	Go to [https://game-on.org/](https://game-on.org/) and click **Play**
2.	Select any authentication method to log in with your username and password for that type.
3.	Click the **Edit Profile** button(the person icon) at the top right.
4.	You should now see the ID and API key at the bottom of the page.  You may need to refresh the page to generate the API key.  You will need to make note of your API key for later in the walkthrough.

## Getting the source code

The source code is located in GitHub, navigate to our [source code](https://github.com/cfsworkload/gameon-room-go.git) and download the ZIP and unzip the code on your local machine.
Alternatively you can use the GitHub CLI to clone the repository with

`git clone https://github.com/cfsworkload/gameon-room-go.git`

## Updating code with Game On! credentials
To create a room with Game On!, some credentials may need to be added to configuration files.

### Configure Go room
All the parameters and credentials are passed in via command line, so there are no pre-configuration steps needed. A room is initialized with parameters passed in to the command:

`$GOPATH/bin/gameon-room-go -c <CONTAINER_IP> -g game-on.org -cp <PORT> -lp <PORT> -r <ROOM_NAME> -id <GAMEON_ID> -apikey <GAMEON_API_KEY>`

The default port as defined in the Dockerfile is the exposed port 3000. If you would like to use another port, or have multiple ports, you would need to expose those in the Dockerfile before creatng your container.


## Make sure a public IP is available in your Bluemix space
This solution requires a free public IP. In order to determine if a public IP is available, you need to find the number of used and your max quota of IP addresses allowed for your space.

To find this information:

1. Log into Bluemix at https://console.ng.bluemix.net.
2. Select **DASHBOARD**.
3. Select the space you where you would like your Wish List app to run.
4. In the Containers tile, information about your IP addresses is listed.
5. Check that the **Public IPs Requested** field has at least one available IP address.

If you have an IP address available, you're ready to start building your Wish List app. If all of your IP addresses have been used, you will need to release one. In order to release a public IP, install the [IBM Containers CLI](https://www.ng.bluemix.net/docs/containers/container_cli_ov.html#container_cli_cfic_installs) plugin, which can be found at the website below.

https://www.ng.bluemix.net/docs/containers/container_cli_ov.html#container_cli_cfic_installs

Once installed:

1. Log into your Bluemix account and space.

  `cf login`  
2. Log into Cloud Foundry for IBM Containers.

  `cf ic login`  
3. List your current external IP addresses.

  `cf ic ip list`
4. Release an IP address.

  `cf ic ip release <public IP>`


## Game On! room on Containers
To build a Game On! room on containers, you create the container locally in Docker using the provided Dockerfile set ups, and then push the container onto Bluemix.

###Build the Docker container locally
1.	Start Docker to run Docker commands. e.g. with the Docker QuickStart Terminal.
2.	In the folder, build the container using run `docker build -t registry.ng.bluemix.net/<your_namespace>/<imageName> .`
3.	To verify the image is now created, you can run `docker images`

###Push and run on Bluemix
1.	Make sure you are logged in with the `cf ic login` command.
2.	Push the container to Bluemix using the command `docker push registry.ng.bluemix.net/<your_namespace>/<imageName>` using the same image name you used to create the local container.

###Run the container on Bluemix
1.	Use the following command to run the container, using the same image name you used to create the local container and port 3000 as the open port as set for the port variable in the server.js. 
	`cf ic run -p <port>:<port> -m 256 -d registry.ng.bluemix.net/<your_namespace>/<imageName>`
2.	With the resulting container ID, use that value to bind to the IP you got earlier with the command: 
	`cf ic ip bind <IP> <Container ID>`

###Accessing a running container
To edit a file on a running container, you can log in to the container with the command 
	`cf ic exec -it <Container ID> /bin/bash`
Note: The container has been created with vim to allow editing. You may need this to edit the files in /usr/src/app.

### Starting a room
To start the room run:
`$GOPATH/bin/gameon-room-go -c <CONTAINER_IP> -g game-on.org -cp <PORT> -lp <PORT> -r <ROOM_NAME> -id <GAMEON_ID> -apikey <GAMEON_API_KEY>`

The parameters are as follows:
- **CONTAINER_IP** - The public IP address to which your container is bound.
- **GAMEON_ID** - The ID given in Game On!
- **GAMEON_API_KEY** - The API key given in Game On!
- **PORT** - The default value is 3000 as it is exposed in the Dockerfile.
- **ROOM_NAME** - Your name for the room.


