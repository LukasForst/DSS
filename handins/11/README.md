# Hand-in 11: Exercise 16.13
___
* *Hannah Eliza Schaible, Lukas Forst*
* Github Repository - [LukasForst/DSS](https://github.com/LukasForst/DSS/tree/master/handins/)
___

## Security policy

This security policy gives the security objectives to protect the voting system of _Danmarksdebatten_.
The system should satisfy the following security objectives:

* A user is qualified to vote if he is user of the public digital signature system and satisfies the demands the vote creator requests.
* Only qualified voters can participate in a vote.
* Voting is possible for all qualified users once until the voting deadline ends.
* After voting users can see the tally of votes so far.
* Only the final result must be published.
* Each voters individual vote has to be kept secret.
* Anonymized data collection of demographical and statistical data can be performed.

## Solution
- The following security policies aim to protect and specify the voting system for the _Danmarksdebatten_ web.
- In order to ensure only eligible votes, the user is qualified to vote in a poll if he/she fullfils the following properties.
    1. A user is part of the public digital signature system.
    2. A user satisfies all extra demands defined by the party who sets up the vote. 
