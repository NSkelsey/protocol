

Usernames
----

All bitcoin addresses that are treated as usernames are interpreted as follows:
-Everything but the first 8 characters are removed
-The last 4 characters of the address are stored as a fingerprint
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

Bulletins
----

For a bitcoin transaction to be interpreted as a Bulletin it must:
- Contain an OP_RETURN that has the following format:
[ 0x2D2D20 app-version 0x202D2D20 [ epochint lat,lon ] ]

###Example
```
-- ahimsa-0.0.1 -- 2131232 (123.3123, 54.3232)
```
- The first output in the tx must have this format: 
```
val [ 0x00000000 num_topics len_data data ]
```
The rest follow:
```
[val [ topic_script ]]....
[val [ data_script ]]....
```
