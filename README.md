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

![this is a cat](http://upload.wikimedia.org/wikipedia/commons/2/22/Turkish_Van_Cat.jpg)
