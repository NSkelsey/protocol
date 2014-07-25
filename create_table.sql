CREATE TABLE blocks (
    hash TEXT NOT NULL, 
    prevhash TEXT, 
    height INT, -- The number of blocks between this one and the genesis block.
    PRIMARY KEY (hash)
    FOREIGN KEY(prevhash) REFERENCES blocks(hash)
);

CREATE TABLE bulletins (
    author TEXT NOT NULL, -- From the address of the first OutPoint used.
    txid TEXT NOT NULL, 
    topic TEXT,   -- UTF-8
    message TEXT, -- UTF-8
    block TEXT, 
    PRIMARY KEY (txid), 
    FOREIGN KEY(block) REFERENCES blocks (hash)
);

