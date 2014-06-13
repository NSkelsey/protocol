from sqlalchemy import (Column, String, LargeBinary,
                        Enum, ForeignKey, create_engine,
                        Integer, Table, Binary)

from sqlalchemy import orm
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import sessionmaker, relationship


engine = create_engine("sqlite:///db.db", echo=True)
Base = declarative_base()


class Tx(Base):
    __tablename__ = "transactions"

    txid    = Column(String, primary_key=True)
    block   = Column(String)
    raw     = Column(Binary)
    
class Block(Base):
    __tablename__ = "blocks"

    hash     = Column(String, primary_key=True)
    prevhash = Column(String)
    
class Bulletin(Base):
    __tablename__ = "bulletins"

    txid    = Column(String, ForeignKey("transactions.txid"), primary_key=True)
    data    = Column(Binary)

class Reference(Base):
    __tablename__ = "references"

    txid    = Column(ForeignKey("bulletins.txid"), primary_key=True)
    addr    = Column(ForeignKey("addresses.addr"), primary_key=True)

class Address(Base):
    __tablename__ = "addresses"

    addr    = Column(String, primary_key=True)


if __name__ == "__main__":
    Base.metadata.create_all(engine)
