# Microservices with a Game On! Room
[Game On!](https://game-on.org/) is both a sample microservices application, and a throwback text adventure brought to you by the wasdev team at IBM. It aims to do a few things:

- Provide a sample microservices-based application.
- Demonstrate how microservice architectures work from two points of view:
 - As a Player: Naviagate through a network/maze of rooms, where each room is a unique implementation of a common API. Each room supports chat, and interaction with items (some of which may be in the room, some of which may be separately defined services as well).
 - As a Developer: Learn about microservice architectures and supporting infrastructure by extending the game with your own services. Write additional rooms or items and see how they interact with the rest of the system.


##Introduction

This walkthrough will guide you through adding a room to a running Game On! server.  You will be shown how to setup a container based room that is implemented in the Go programming language.  There are also instructions for deploying a similar room written in Node.js that is deployed as a Cloud Foundry application in Bluemix.  Both of these rooms take a microservice approach to adding a room to a running Game On! text adventure game server.  

### Installation prerequisites

Gameon-room-go when deployed using containers requires:

 - [Docker](https://docs.docker.com/engine/installation/) - You will need Docker on your host machine to create a docker image to be pushed to IBM Bluemix Containers
 - [IBM Containers CLI](https://www.ng.bluemix.net/docs/containers/container_cli_ov.html#container_cli_cfic_install)

## Create Bluemix accounts and log in
To build a Game On! room in Bluemix, you will first need a Bluemix account. 

### Sign up and log into Bluemix and DevOps
Sign up for Bluemix at https://console.ng.bluemix.net and DevOps Services at https://hub.jazz.net. When you sign up, you'll create an IBM ID, create an alias, and register with Bluemix.


## Get Game On! API Key and User ID
For a new room to register with the Game On! server, you must first log into game-on.org and sign in using one of several methods to get your Game On! user ID and ApiKey.

1.	Go to https://game-on.org/ and click **Play**
2.	Select any authentication method to log in with your username and password for that type.
3.	Click the **Edit Profile** button(the person icon) at the top right.
4.	You should now see the ID and API key at the bottom of the page.  You may need to refresh the page to generate the API key.  You will need to make note of your API key for later in the walkthrough.

## Getting the source code

The source code is located on [GitHub](https://github.com/cfsworkload/).  Download the ZIP and unzip the code, or using the GitHub CLI clone the repository with

    git clone https://github.com/cfsworkload/

## Updating Code with Game On! credentials

### Configure Node.js room
For the Node.js room, the server.js file contains information about the new room and your user credentials. The user credentials must be set manually. 

- **gameonAPIKey** - Use the ApiKey value from the game-on.org user settings page.
- **gameonUID** : Use the Id value from the game-on.org user settings page.
- **endpointip** : Use the IP requested or retrieved from earlier steps.
- **port** : The default value is 3000. This port must be opened when running the container.

### Configure Go room



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

  `cf ic login`  
2. List your current external IP addresses.

  `cf ic ip list`
3. Release an IP address.

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
To start the room for your specific language run:
	- Node.js: `node server.js &`
	- Go: ./gameon-room-go -c <containerIP> -g <containerIP> -p <port> -r <RoomName>
