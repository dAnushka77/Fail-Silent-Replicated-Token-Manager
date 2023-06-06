# Fail-Silent-Replicated-Token-Manager

Fail-Silent Replicated Token Manager with Atomic Semantics
Project 3, CMSC 621, Spring 2023
Project done by – Anushka Dhekne (VD19739)

In this fail silent replicated token manager, I have considered a single client and 3 server system. This allows in easy replication of data, thus making the system fail-silent.  The steps I followed for the setup and implementation of the project are given below.

Step 1) Protobuf File Setup – to create grpc.pb.go and pb.go files using the proto3 module.

Step 2) Command used – go mod tidy – to complete setup and configuration of the individual go files as well as the protobuf files.
 
Step 3) Creating a config.yaml file to specify the reader and writer ports. The location of these ports is accessed through the code.

Step 4) Start all servers as mentioned by their port numbers in the config.yaml file. For example, in my case, the writer ports were at localhost 9001 and 9002. In this case, whatever gets written to 9001 is replicated on port 9002 and vice-a-versa. 
Hence, start the connection by initializing ports 9001 and 9002 on the servers.

Step 5) Once all the required servers are up and running, open a new terminal to start the client  and implement the first write operation using the following command on id 006 –
go run client/client.go -write -id 006 -name abc -low 0 -mid 10 -high 100

Server 3, running on port 9003 does not have read/write from config.yaml so there will be no output on server 3 whatsoever.

Step 6) Trying the same for id = 007. 

Writer is only activated for the 9002 server port, therefore server 2 with port 9002 gets updated. The server adds 007 to its pre-existing list of tokens it holds, and it now becomes [006, 007]. Thus, client gets updated.

Step 7) If we check read operation for non-existing token ids – 001 and 002, we see that the result is blank, and we do not get any desired output.

Step 8) Similar to write operation, the same can be applied to the read operation, update operation and drop operation.


