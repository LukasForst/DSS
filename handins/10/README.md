# Hand-in 10: 15.2
___
* *Hannah Eliza Schaible, Lukas Forst*
* Github Repository - [LukasForst/DSS](https://github.com/LukasForst/DSS/tree/master/handins/)
___

## Question 1 
If the URL-check is not perfromed, and the adversary can, as stated in the exercise, completely control the connection, 
then a possible attack could look like this: the adversary takes over the connection from the very beginning, and pretends to be the target server
for the first sent packet. He responds with his own certificate, however this will not be checked. Now the user is securely connected to an adversary's website _W<sub>b</sub>_, although he intended to securely connect to website _W<sub>a</sub>_.

## Question 2
Since the adversary controls the whole connection and network, and the URL is not present in the certificate, the browser has no possibility of successfully solving the above issue.

## Question 3
To construct an AKE protocol π satisfying definition D from an AKE protocol π' that satisfies D', one can run protocoö π' first, and then use the output identity as an input to run the protocol π afterwards. Then definition D will be satisfied.

The TLS AKE without URL check satisfies the definition D', since this protocol takes no identity as input. Meanwhile, the protocol D takes the identity ad input, and would therefore be equivalent to a TLS AKE with URL check.

