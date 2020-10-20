⊥


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
