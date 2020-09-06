# Hand-in 1 Report

> Test your system and describe how you tested it.

We tested our p2p network manually by creating following topologies:
* chain - all peers are connected to the previous peer 
* one central peer - all other peers are connected to the central peer
* star topology - multiple "central" peers connected between each other

Then we sent message on every peer and observed if the message was received by all others.


> Argue that your system has eventual consistency in the sense that if all clients
> stop typing, then eventually all clients will print the same set of strings.

Each peer sends all received messages, that are unique (= has not been shared by the peer yet)
 to his network (all other peers connected to this one) exactly just once.
 Moreover, storage for messages (map which records whether the message has been already seen in the past),
 uses mutex for read and for write operations. 
 For that reason just one thread can access the storage.
 This results in the atomicity of the operation "Distribute message if not seen previously".
 
If all peers are connected, then they eventually receive set,
 but not necessarily in the same order.  
