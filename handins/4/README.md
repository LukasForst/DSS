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

The second specification of the security policy, which states that users can not get information about other users requests could be potentially broken as well. If the RSA implementation, used in the protocol, is not using paddings/salting or other means of nonces, adversary can simply try to "bruteforce" the original request. That can be done by generating the request and encrypting it with the public key of the database (which is by definition public) and then comparing the ciphertexts with the ciphertext of the original request. Using this bruteforce attack, adversary can determine which user asked for what data thus violating the second security policy.

> are there any general conclusions one could draw from this example?

The example shows that it is important to sign and then encrypt the requests, since the signatures could otherwise be manipulated.
This also applies generally for the use of RSA encryption to ensure secure communication. So to sum it up, it's always important to sing the plaintext instead of ciphertext.
