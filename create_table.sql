CREATE TABLE blocks (
	hash VARCHAR NOT NULL, 
	prevhash VARCHAR, 
	PRIMARY KEY (hash)
);
CREATE TABLE addresses (
	addr VARCHAR NOT NULL, 
	PRIMARY KEY (addr)
);
CREATE TABLE transactions (
	txid VARCHAR NOT NULL, 
	block VARCHAR, 
	raw BLOB, 
	PRIMARY KEY (txid)
);
CREATE TABLE bulletins (
	txid VARCHAR NOT NULL, 
	data BLOB, 
	PRIMARY KEY (txid), 
	FOREIGN KEY(txid) REFERENCES transactions (txid)
);
CREATE TABLE "references" (
	txid VARCHAR NOT NULL, 
	addr VARCHAR NOT NULL, 
	PRIMARY KEY (txid, addr), 
	FOREIGN KEY(txid) REFERENCES bulletins (txid), 
	FOREIGN KEY(addr) REFERENCES addresses (addr)
);
