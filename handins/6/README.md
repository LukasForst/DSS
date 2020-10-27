# Hand-in 6 Report
___
* *Hannah Eliza Schaible, Lukas Forst*
* Github Repository - [LukasForst/DSS](https://github.com/LukasForst/DSS/tree/master/handins/6)
___

### How we tested the program
Established the peer network, run the `test.go` transactions and observed
whether the transactions are executed and checked whether the final sums check 
in the debugger. 

### How to the peer
1. Run as many peers as you want with command `make run`
    * when the first peer asks for the IP and port, simply don't put anything
    and just press enter

### How to run test stuff
Please see `test.go` file.

1. comment function `func main()` in  `main.go` and uncomment it in `test.go`
1. run `make run` again, the test should start, again you can run as many processes 
of test as you want to, just give it unique Peer Id (can be just one letter) on prompt:
```
> Enter unique test peer ID
```


### How to test signatures
In the `test.go` on line 84, we're signing the request like that:
```go
transaction.ComputeAndSetSignature(accounts[from].PK)		
```
to verify that the program can recognize signatures, one can simply remove
this line or with some probability sign it with `to` secret key instead of `from` one.
The program output should be then following:
```
> - Transaction B7 received.
> - Could not verify signature: crypto/rsa: verification error
> - Transaction B7 has incorrect signature!
``` 
