# Hand-in 6 Report
___
* *Hannah Eliza Schaible, Lukas Forst*
* Github Repository - [LukasForst/DSS](https://github.com/LukasForst/DSS/tree/master/handins/6)
___


Steps to reproduce:

1. Run as many peers as you want with command `make run`
    * when the first peer asks for the IP and port, simply don't put anything
    and just press enter
1. comment function `func main()` in  `main.go` and uncomment it in `test.go`
1. run `make run` again, the test should start, again you can run as many processes 
of test as you want to, just give it unique Peer Id (can be just one letter) on prompt:
```
> Enter unique test peer ID
```
