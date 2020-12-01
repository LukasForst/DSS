# Hand-in 10: 15.2
___
* *Hannah Eliza Schaible, Lukas Forst*
* Github Repository - [LukasForst/DSS](https://github.com/LukasForst/DSS/tree/master/handins/)
___

## Question 1 
If the URL-check is not perfromed, and the adversary can, as stated in the exercise, completely control the connection, 
then a possible attack could look like this: the adversary takes over the connection from the very beginning, and pretends to be the target server
for the first sent packet. He responds with his own certificate, however this will not be checked. Now the user is securely connected to an adversary's website _W<sub>b</sub>_, although he intended to securely connect to website _W<sub>a</sub>_.
