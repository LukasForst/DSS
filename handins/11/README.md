# Hand-in 11: Exercise 16.13
___
* *Hannah Eliza Schaible, Lukas Forst*
* Github Repository - [LukasForst/DSS](https://github.com/LukasForst/DSS/tree/master/handins/)
___

Note about the following text - we didn't really understand why they are talking about the statistical data, because they're not part of the original assignment/specification.

## Security policy
This security policy gives the security objectives to protect the voting system of _Danmarksdebatten_.
The system should satisfy the following security objectives:

- In order to ensure only eligible votes, the vote is qualified to vote in a poll if he/she fullfils the following properties.
    1. A voter is a user of the public digital signature system in order to ensure the authenticity and the non-repudiation of the votes.
    2. A voter satisfies all extra demands defined by the party who sets up the vote.
- In order to ensure the validity of the votes and the integrity of the system, the voters can cast votes only during the period, when the poll is opened (meaning, after the poll was created until the poll's deadline).
- In order to ensure the authenticity for the votes, the voter can cast only one vote per poll. Thus, once a voter voted, he/she can't vote in the same poll anymore.
- The partial results, nor the votes themselves, are not public, when the poll is opened (meaning, after the poll was created until the poll's deadline). This security policy ensures that the votes are genuine and unbiased by the partial results.
- In order to ensure the integrity of the system, the voter has to be able to verify that his/her vote was counted. In order to ensure the confidentiality, the voter can check only the status of the vote, that belongs to the said voter.
