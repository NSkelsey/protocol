###Implementation notes

####ahimsa app

work order:
- an app that can trasmit a normal message to the network
   * test on mainent
- encode messages in address outs of transaction
   * test on mainnet
- add an op_return with metadata
   * test on mainnet
- Let user specify what they want to send
   * go live


####ahimsad

Be aware:
    - by querying bitcoind getaddr, your connected node will leak localhost:8888. That is to say an attacker can discover you are running ahimsad easily 
        * todo add config to bitcoind

