## Original Dolev-String
Taken from an anonymous student from the [Blackboard Forum](https://blackboard.au.dk/webapps/discussionboard/do/message?action=list_messages&course_id=_138620_1&nav=discussion_board_entry&conf_id=_273509_1&forum_id=_189724_1&message_id=_291047_1).


### Initialize:
Initialize for node _P<sub>i</sub>_
- Query toyPKI for the private key _sk<sub>i</sub>_, the private key for _P<sub>i</sub>_
- Queery toyPKI for the public key of all nodes _n_, _(vk<sub>1</sub>,...,vk<sub>n</sub>)_
- Create a map _Relayed_<sub>i</sub>, where for all messages _m_ and all broadcast-ids _bid_:
	- _Relayed<sub>i</sub>(bid,m) = ⊥_

### Broadcast:
On input _(bid, P<sub>i</sub>, m)_ and _Cast<sub>i</sub>_ do as follows:
- Compute _o<sub>j</sub> <- Sig<sub>ski</sub> (bid, m)_
- Compute _SigSet_ = {_o<sub>i</sub>_}
- Set _Relayed<sub>i</sub>(bid,m) = T_
- Send _(bid, P<sub>i</sub>, m, SigSet)_ to all parties

### Relay:
In round _r_ for node P<sub>j</sub> with input _(bid, P<sub>i</sub>, m, SigSet)_ do as follows:
- if _P<sub>i</sub> = P<sub>j</sub>_ do nothing, you have received your own message
- If _Relayed<sub>j</sub>(bid, m) = T_, do nothing, you have already relayed the message
- Check if _SigSet_ is valid
- _SigSet_ is valid if it has:
	- _(r - 1)_ distinct signatures
	- 1 signature from the original sender, _R<sub>i</sub>_
	- No signature from itself, _R<sub>j</sub>_
	- and for all signatures: _Ver<sub>vkk</sub> = T_, meaning all signatures verify to true
- if _SigSet_ is invalid do nothing
- otherwise compute _o<sub>j</sub> <- Sig<sub>skk</sub> (bid, m)_
- compute _SigSet' <- SigSet U {o<sub>j</sub>}_
- send _(bid, P<sub>i</sub>, m, SigSet')_ to all parties
- set _Relayed<sub>j</sub> (bid, m) = T_

### Output:
In round _n+2_ with input _(bid, P<sub>i</sub>, m, SigSet)_ do as follows:
- if there is one and only message _m_ such that _Relayed<sub>i</sub>(bid, m) = T_
- then output _(bid, P<sub>i</sub>, m)_ on _Cast<sub>j</sub>_
- else output _(bid, P<sub>i</sub>, NoMsg)_ on _Cast<sub>j</sub>_

## Modification to prevent DDoS attack by Byzantine corrupted sender
We assume that the sender is byzantine corrupted and is sending a different message _m_ with the same broadcast id _bid_ and with a valid signature.

Now, the original Dolev-Strong uses _Relayed<sub>i</sub>_, which stores _(bid, m)_. As the broadcast id is by definition unique for an unique message, we don't need to store both of the values - the _bid_ and the _m_, but instead we would store just the _bid_. Thus when the sender sends a different message _m_ with the same _bid_, the parties won't accept the message and wouldn't do anything.

The modified version is then following:
___
### Initialize:
Initialize for node _P<sub>i</sub>_
- Query toyPKI for the private key _sk<sub>i</sub>_, the private key for _P<sub>i</sub>_
- Queery toyPKI for the public key of all nodes _n_, _(vk<sub>1</sub>,...,vk<sub>n</sub>)_
- Create a map _Relayed_<sub>i</sub>, where for all messages _m_ and all broadcast-ids _bid_:
	- _Relayed<sub>i</sub>(bid,m) = ⊥_
- **Create an empty map _Encountered<sub>i</sub>_ where we store a _bid_ for _m_**

### Broadcast:
On input _(bid, P<sub>i</sub>, m)_ and _Cast<sub>i</sub>_ do as follows:
- Compute _o<sub>j</sub> <- Sig<sub>ski</sub> (bid, m)_
- Compute _SigSet_ = {_o<sub>i</sub>_}
- Set _Relayed<sub>i</sub>(bid,m) = T_
- Send _(bid, P<sub>i</sub>, m, SigSet)_ to all parties

### Relay:
In round _r_ for node P<sub>j</sub> with input _(bid, P<sub>i</sub>, m, SigSet)_ do as follows:
- if _P<sub>i</sub> = P<sub>j</sub>_ do nothing, you have received your own message
- If _Relayed<sub>j</sub>(bid, m) = T_, do nothing, you have already relayed the message
- **if _Encountered<sub>i</sub>[bid] = m_ or does not exist, continue**
- **however, if _Encountered<sub>i</sub>[bid] != m_ do nothing, you received invalid combination of _bid_ and _m_**
- Check if _SigSet_ is valid
- _SigSet_ is valid if it has:
	- _(r - 1)_ distinct signatures
	- 1 signature from the original sender, _R<sub>i</sub>_
	- No signature from itself, _R<sub>j</sub>_
	- and for all signatures: _Ver<sub>vkk</sub> = T_, meaning all signatures verify to true
- if _SigSet_ is invalid do nothing
- otherwise compute _o<sub>j</sub> <- Sig<sub>skk</sub> (bid, m)_
- compute _SigSet' <- SigSet U {o<sub>j</sub>}_
- send _(bid, P<sub>i</sub>, m, SigSet')_ to all parties
- set _Relayed<sub>j</sub> (bid, m) = T_

### Output:
In round _n+2_ with input _(bid, P<sub>i</sub>, m, SigSet)_ do as follows:
- if there is one and only message _m_ such that _Relayed<sub>i</sub>(bid, m) = T_
- then output _(bid, P<sub>i</sub>, m)_ on _Cast<sub>j</sub>_
- else output _(bid, P<sub>i</sub>, NoMsg)_ on _Cast<sub>j</sub>_
