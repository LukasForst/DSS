# Hand-in 9 Report
___
* *Hannah Eliza Schaible, Lukas Forst*
* Github Repository - [LukasForst/DSS](https://github.com/LukasForst/DSS/tree/master/handins/9)
___

## Run the program

To run it
1. go to `zz-main.go`,
1. uncomment the `main`
1. run `make run`
1. hit enter
1. copy the IP address and the port (this is the peer #1)
1. run `make run`
1. put there the IP and port for peer #1
1. comment `main` in `zz-main.go`
1. go to `zz-test.go`
1. comment `main` in `zz-test.go`
1. run `make run`
1. enter `A`
1. enter IP of the peer #1
1. enter IP of the peer #2 (the output of the second peer)
1. hit enter

## Testing

For testing, an amount of accounts are initialized with the initial states and the private keys, and the initial model is set up. The connections between the peers are established and the test transactions are run.
The tests were conducted with the given ten peers, in different scenarios under varying conditions, such as differing account balances and transactions, including valid and invalid transaction attempts. The resulting account balances were checked.
