Data
----
We have defined a format that describes how to store data in a bitcoin transaction

To store data a transaction must contain an OP_RETURN where all subsequent txouts are used for storage.

Everything above that OP_RETURN is ignored. (So change outputs can be specified above)

NOTE var:int means variable var requires int bytes of space.

The format for the message header is:


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

- Everything but the first 8 characters are removed
- The last 4 characters of the address are stored as a fingerprint
 * Optioinally that fingerprint can be converted into an rgb value


###Examples
```
1nskeLseycHqhW8QR9jya1TaZeXYwyy7s --> nskelseyc yy7s
1askuckTe53qxHdPqXEvdu8WdCXWn6Cmb --> askuckte5 6Cmb
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
