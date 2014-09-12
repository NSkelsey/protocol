Definitions of 

###Storage Format
The storage format is slightly complicated. We encode bulletins in the 20 byte slices
used for bitcoin addresses in Pay2PubKeyHash transactions. The Tx indicates that it is
a public bulletin by making the first 8 bytes of that first 20 byte slice equal to `0x
425245544852454e`. The actual bulletin itself is then encoded in a protocol buffer for 
effeciency!?

###Database Schema
As of version 0.0.0, the database consists of two tables. Blocks and bulletins are
the only objects whose existence we track.



Prior Work
======

There are quite a few projects from which we have taken ideas, concepts and source
code. The big ones are:
- [btcd](https://github.com/conformal/btcd) 
    - The developers at Conformal have developed some awesome bitcoin libraries. 
    We have used them extensively.
- [twister](https://github.com/miguelfreitas/twister-core)
    - The intended use case of this project and ours is excatly the same. We just
    elected to use bitcoin's infrastructure, not set up our own.
- [bitmessage](https://github.com/Bitmessage/PyBitmessage)
    - A similar messaging tool that encrypts messages preflight. 
- [bitchirp](https://bitchirp.org/)
    - The original distributed version of twitter
- [Proof of Existence](http://www.proofofexistence.com/)
    - The first tool we were aware of that uses the blockchain for its distrbuted 
    timestamp and data storage.
- [Counter Party](https://www.counterparty.co/)
    - Our bitcoin daemon and database schema was modeled after the one created by
    PhantomPhreak.
