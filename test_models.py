from models import *

Session = sessionmaker(bind=engine)

def test():
    s = Session()

    addr = Address(addr="ThisiABTCADDR")
    
    tx = Tx(txid="deadbeef", raw=bytearray([0,0,0,0,0]))
    bulletin =  Bulletin(data="asdfafy msg")

    bulletin.txid = tx.txid
    

    print bulletin.txid, tx.txid

    ref = Reference()
    ref.addr = addr.addr
    ref.txid = bulletin.txid

    s.add(tx)
    s.add(addr)
    s.add(bulletin)
    #s.add(ref)

    s.commit()

if __name__ == '__main__':
    test()

    

