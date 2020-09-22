# Hand-in 2 Report
___
* *Hannah Eliza Schaibe, Lukas Forst*
* Github Repository - [LukasForst/DSS](https://github.com/LukasForst/DSS/tree/master/handins/2)
___

> Test you system and describe how you tested it

We established the network as in the [first assignment](https://github.com/LukasForst/DSS/tree/master/handins/1#hand-in-1-report)
and then executed `test.go` with different parameters, connecting to different peers,
running the different amount of transactions.

> Discuss whether connection to the next ten peers is a good strategy with respect
> to connectivity. In particular, if the network has 1000 peers, 
> how many connections need to break to partition the network?

As each peer is connecting to the next ten peers, the total connections to one peer
is 20 peers (ten before connecting to this peer, the peer then connects to ten after),
for that reason if we break 20 connections (ten in one segment/line, 
and ten in the second segment/line), the network will be disconnected.

As for the question, whether this is a good strategy, the number of connections
is reasonable (20 at one time). However, this strategy is not really robust, as 
the peers are connected to each other in alphabetical order (at least in our implementation).
The main concern here is that if multiple peers share the geolocation and IP addresses,
it could happen, that this particular location would fall down due to 
electricity/internet problems. Which could potentially result in disconnecting the network.


> Argue that your system has eventual consistency if all processes are correct
> and the system is run in two-phase mode

During the first phase od building the network, every peer eventually receives 
the whole list of all peers in the network. That is because when it joins the network,
it receives whole connection table from the network peer. After joining, it announces
its presence, subsequently other peers receive a presence message and if they 
already have the record, they ignore it. However, if they don't know the peer,
they update the peers list and broadcast the presence message to all 
other connections they have. This process ensures that the network is connected.

If the transaction is received, peer checks whether it was already executed by 
checking the `ID` of the transaction, and the transaction log.
If the transaction was executed, it is ignored and nothing happens next.
If the transaction was not received before, it is executed, put to the transaction log
and broadcast to all other connected peers. This process ensures that each transaction
was executed exactly once and moreover, as the network is connected, that each peer
eventually receives the transaction message. Thus, the transaction logs should be consistent
at the end of the execution.

The system wouldn't be consistent if it does not run in the two-phase mode,
mainly because latecomers wouldn't receive transactions, that were distributed 
on the network before they joined.


> Assume we made the following change to the system:
> When a transaction arrives, it is rejected if the sending account goes below 0. 
> Does your system still have eventual consistency? Why or why not?

Thanks to the associative property of adding and subtracting numbers,
it is not necessary to use Vector clock in the basic implementation of the protocol.
However, in order to be able to properly track `0` balance on the account 
and possibly rejecting transaction, we would need to ensure that the transactions were processed
in the correct order. For that reason, some mechanism, such as Vector Clocks must be used.
Our system does not contain Vector clock, so it shouldn't be used like that,
because it could end up in different state, depending on the transactions order. 
