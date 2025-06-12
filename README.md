# Socket IP Sub-Network Calculator

This is a university project made to test and understand the use of **sockets** for **Network** class.
Our objective is to make a **calculator for IP sub-networks** using a **client and server socket architeture**.

The project work as the following steps:
+ **Client** connects to **Server**
+ **Client** send an **IP**, **Mask** and the **number of sub-networks**
+ **Server** process the information and returns the **sub-network address** and **mask**, **starting address** and **end address**.

## Tools and Dependencies

The project uses [golang](https://go.dev) as the main language for building the server inner functionality.
[Docker](https://www.docker.com) was also used the run the server in a separated environment. 

### Dependencies

Before using the project ensure you have installed:
- [Go Programming Language](https://go.dev)
- [Docker](https://www.docker.com)
## Running

To run the project  run the following commands:

- **Install** the project and **go into** the root folder
```bash
git clone "https://github.com/Thales3006/socket-ip-calc"
cd socket-ip-calc
```
- **Build** the server
```bash
# building
docker build -t socket-ip-calc .
```
- **Run** the server and client
```bash
# running server
docker run -p 8080:8080 socket-ip-calc
# running client
go run cmd/client/main.go
```
