# Socket IP Sub-Network Calculator

This is a university project made to test and understand the use of **sockets** for **Network** class.
Our objective is to make a **calculator for IP sub-networks** using a **client and server socket architecture**.

The project works as follows:
+ **Client** connects to **Server**
+ **Client** sends an **IP**, **Mask** and the **number of sub-networks**
+ **Server** process the information and returns the **sub-network address** and **mask**, **starting address** and **end address**.

## Tools and Dependencies

The project uses [golang](https://go.dev) as the main language for building the server inner functionality.
[Docker](https://www.docker.com) was also used to run the server in a separate environment. 

### Dependencies

Before using the project ensure you have installed:
- [Go Programming Language](https://go.dev)
- [Docker](https://www.docker.com)
## Running

To run the project, run the following commands:

- **Install** the project and **go into** the root folder
```bash
git clone "https://github.com/Thales3006/socket-ip-calc"
cd socket-ip-calc
```
- **Run** and **Build** the server
  you can repeat this command if you wish to rebuild the server (no need to shutdown before it).
```bash
# running server process:
docker-compose up --build -d 
# if you wish to stop the server:
# docker-compose down
```

- **Run** the client
```bash
# running client:
go run cmd/client/main.go
```
