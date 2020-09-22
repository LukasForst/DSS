# Hand-in 2 Report
___
* *Hannah Eliza Schaibe, Lukas Forst*
* Github Repository - https://github.com/LukasForst/DSS
___

> Test you system and describe how you tested it

pass

> Discuss whether connection to the next ten peers is a good strategy with re-
  spect to connectivity. In particular, if the network has 1000 peers, how many
  connections need to break to partition the network?

pass

> Argue that your system has eventual consistency if all processes are correct
> and the system is run in two-phase mode

pass



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
