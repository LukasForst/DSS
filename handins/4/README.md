# Hand-in 4 Report
___
* *Hannah Eliza Schaibe, Lukas Forst*
* Github Repository - [LukasForst/DSS](https://github.com/LukasForst/DSS/tree/master/handins/4)
___

Answers to the questions in the [test.go](test.go).

## 6.1
> One of these solutions fails to satisfy the security policy. Which one, and why? 

The first suggestion, which proposes encrypt and then sign, fails the given security policy.
In contrast to the second suggestion, which proposes signing and then encryption, in the first suggestion the signature is not being encrypted. With this procedure it would be possibe for another user B to manipulate a request send by user A such that user B replaces A's signature by his own on the request originally sent by user A. This violates the first specification of the security policy, which states that the database D has to be able to determine which user sent the request. If a user B manipulates the request originally sent by user A, the database can not correctly determine the sender of the request anymore, as it would determine user B as the sender. This would not be possible in the second suggestion (sign-then-enc). 
Furthermore, as in the first suggestion (enc-then-sign) user B can see the signature A provided with his private key, user B could now create messages, encrypt them with A's public key and append user A's signature that he knows to them, as well as A's username, which again fails the first specification of the security property, the explicit determination of the sender of a request by database D.
Like this, the database could treat the request as a valid request and return an answer with the requested data, encrypted under A's public key. Since B doesn't know A's private key, he can not decrypt this data though.
The second specification of the security policy, which states that users can not get information about other users requests, is therefore not broken by the first suggestion, as user B does not know user A's secret key, which makes it impossible to decrypt answers the database returns to A, and also it is not possible for B to decrypt the specific request that A sent.

> are there any general conclusions one could draw from this example?

The example shows that it is important to sign and then encrypt the requests, since the signatures could otherwise be manipulated.
This also applies generally for the use of RSA encryption to ensure secure communication.
