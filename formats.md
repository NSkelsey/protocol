Data
----
We have defined a format that describes how to store messages in bitcoin transactions. The first output of a message must have a header in its data field.

The rest of the message is encoded as a block of txouts. The header gives enough information for the object to be decoded from the following txouts. Since data encoding formats are likely to change all datatypes are versioned with corresponding headers that provide extra information about that object. Allowing us to reconstruct it elsewhere.

NOTE var:int means variable var requires int bytes of space.

The format for the message header is unclear.....:


| Value | Scriptlen | Script |
|-------|:---------:|--------|
| 0     |   len     | OP_RETURN 0xdeadbeef data_type:8 verison:2 type_headers:? |


Currently supported data_types are:


| Data Type | Version | Type Headers |
|-----------|:-------:|--------------|
| bulletin  |    1    | msglen:4 numTopics:1 numNames:1 userAgent:10 |

For encoding a v1 bulletin in a txout the format is:

| Value | Scriptlen | Script |
|-------|-----------|--------|
| 546   |     25    | OP_DUP OP_HASH160 data:20 OP_EQUALVERIFY OP_CHECKSIG |

Where the order of the txouts is: topics; usernames; msg


An example message header follows:
```

6a1edeadbeef62756c6c6574696e00010000000e03016168696d73612d776562
    |      |               |   |       | | |                    |
OPS   Mag       datatype     ver  len  nT nN    user-agent

Decodes as:
('bulletin', 1, (14, 3, 1, 'ahimsa-web'))
  datatype  ver  len nT nN   user-agent
```



Usernames
----

All bitcoin addresses that are treated as usernames are interpreted as follows:

- Everything but the characters between [1:9] are discarded
- The last 4 characters of the address can be used as a fingerprint
 * Optionally that fingerprint can be converted into an rgb value
- A valid address must contain a capital letter at postion [0]
 * It must not be an L or an X (this keeps usernames and topics mutually exclusive)


###Examples
####Valid
```
1NSkeLseycHqhW8QR9jya1TaZeXYwyy7s --> NSkelseyc yy7s
1AskuckTe53qxHdPqXEvdu8WdCXWn6Cmb --> Askuckte5 6Cmb
```
####Invalid
```
1nskeLseycHqhW8QR9jya1TaZeXYw5jdt
16Frogsw3CHVqiBfmC3vCLrXPob6wLQSP4
1LaRRy1dhB6rneeZazy9hQg9iUpo961Mt9
```
Topics
----

For a bitcoin address to be treated as a topic the address will abide by:
- The address will have at least one trailing X before the last four characters 
- Every character in the topic must be lowercase except for L

To transform it into a topic the following operations will occur:
- The version byte will be removed from the topic
- The X train and checksum will be removed from the tail
- All capital X's within the remaining string will be interpreted as spaces
- Any o's adjacent to numbers will be intrepreted as 0's
- All words will have their first character captialized

###Examples
####Valid
```
1hedgefundsXdrainXsocietyXXXYUebE --> Hedge Funds Drain Society
1thisXisXaXpubLicXLedgerXXXbMieDH --> This Is A Public Ledger
1protestsXinX2o14XXXXXXXXXXh6JThb --> Protests In 2014
```
####Invalid
```
1ProTestSXiNX2o14XXXXXXXXXXXZVQQq6
```
