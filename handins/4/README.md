# Hand-in 4 Report
___
* *Hannah Eliza Schaibe, Lukas Forst*
* Github Repository - [LukasForst/DSS](https://github.com/LukasForst/DSS/tree/master/handins/4)
___

Answers to the questions in the [test.go](test.go).

## 6.1
> One of these solutions fails to satisfy the security policy. Which one, and why? 

The solution where you first encrypt and then sign the encrypted data. 
The signature does not prove that the sender was aware of the context of the plaintext, thus it violates non-repudiation.

> are there any general conclusions one could draw from this example?

Yes, the following reasoning:
> The signature does not prove that the sender was aware of the context of the plaintext, thus it violates non-repudiation.

is general. As the main reason why we sign the data is to ensure non-repudiation, the signing itself is useless.
