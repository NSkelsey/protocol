package ahimsa

import (
	"bytes"
	"fmt"

	"code.google.com/p/goprotobuf/proto"
	"github.com/conformal/btcnet"
	"github.com/conformal/btcscript"
	"github.com/conformal/btcutil"
	"github.com/conformal/btcwire"
)

var (
	ProtocolVersion uint32 = 0x1
	Magic                  = [8]byte{
		0x42, 0x52, 0x45, 0x54, 0x48, 0x52, 0x45, 0x4e, /* | BRETHREN | */
	}
)

type Author string

type Bulletin struct {
	Txid    *btcwire.ShaHash
	Block   *btcwire.ShaHash
	Author  string
	Topic   string
	Message string
}

func extractData(txOuts []*btcwire.TxOut) ([]byte, error) {
	// Munges the pushed data of TxOuts into a single universal slice that we can
	// use as whole message.

	alldata := make([]byte, 0)

	first := true
	for _, txout := range txOuts {

		pushMatrix, err := btcscript.PushedData(txout.PkScript)
		if err != nil {
			return alldata, err
		}
		for _, pushedD := range pushMatrix {
			if len(pushedD) != 20 {
				return alldata, fmt.Errorf("Pushed Data is not the right length")
			}

			alldata = append(alldata, pushedD...)
			if first {
				// Check to see if magic bytes match and slice accordingly
				first = false
				lenM := len(Magic)
				if !bytes.Equal(alldata[:lenM], Magic[:]) {
					return alldata, fmt.Errorf("Magic bytes don't match, Saw: [% x]", alldata[:lenM])
				}
				alldata = alldata[lenM:]
			}

		}

	}
	// trim trailing zeros
	for j := len(alldata) - 1; j > 0; j-- {
		b := alldata[j]
		if b != 0x00 {
			alldata = alldata[:j+1]
			break
		}
	}
	return alldata, nil
}

func NewBulletin(tx *btcwire.MsgTx, author string, blkhash *btcwire.ShaHash) (*Bulletin, error) {
	// Creates a new bulletin from the containing Tx, supplied author and optional blockhash

	// unpack txOuts that are considered data, We are going to ignore extra junk at the end of data
	wireBltn := &WireBulletin{}

	// Bootleg solution, but if unmarshal fails slice txout and try again until we can no more or it fails
	var err error
	for j := len(tx.TxOut); j > 1; j-- {
		rel_txouts := tx.TxOut[:j] // slice off change txouts
		err = *new(error)
		bytes, err := extractData(rel_txouts)
		if err != nil {
			continue
		}

		err = proto.Unmarshal(bytes, wireBltn)
		if err != nil {
			continue
		} else {
			// No errors, we found a good decode
			break
		}
	}
	if err != nil {
		return nil, err
	}

	hash, _ := tx.TxSha()
	bltn := &Bulletin{
		Txid:    &hash,
		Block:   blkhash,
		Author:  author,
		Topic:   wireBltn.GetTopic(),
		Message: wireBltn.GetMessage(),
	}

	return bltn, nil
}

func NewBulletinFromStr(author string, topic string, msg string) (*Bulletin, error) {
	if len(topic) > 30 {
		return nil, fmt.Errorf("Topic too long! Length is: %d", len(topic))
	}

	if len(msg) > 500 {
		return nil, fmt.Errorf("Message too long! Length is: %d", len(msg))
	}

	bulletin := Bulletin{
		Author:  author,
		Topic:   topic,
		Message: msg,
	}
	return &bulletin, nil
}

func (bltn *Bulletin) TxOuts(toBurn int64, net *btcnet.Params) ([]*btcwire.TxOut, error) {
	// Converts a bulletin into public key scripts for encoding

	rawbytes, err := bltn.Bytes()
	if err != nil {
		return []*btcwire.TxOut{}, err
	}

	numcuts, _ := bltn.NumOuts()

	cuts := make([][]byte, numcuts, numcuts)
	for i := 0; i < numcuts; i++ {
		sliceb := make([]byte, 20, 20)
		copy(sliceb, rawbytes)
		cuts[i] = sliceb
		if len(rawbytes) >= 20 {
			rawbytes = rawbytes[20:]
		}
	}

	// Convert raw data into txouts
	txouts := make([]*btcwire.TxOut, 0)
	for _, cut := range cuts {

		fakeaddr, err := btcutil.NewAddressPubKeyHash(cut, net)
		if err != nil {
			return []*btcwire.TxOut{}, err
		}
		pkscript, err := btcscript.PayToAddrScript(fakeaddr)
		if err != nil {
			return []*btcwire.TxOut{}, err
		}
		txout := &btcwire.TxOut{
			PkScript: pkscript,
			Value:    toBurn,
		}

		txouts = append(txouts, txout)
	}
	return txouts, nil
}

func GetAuthor(authorTx *btcwire.MsgTx, i uint32, params *btcnet.Params) (string, error) {
	// Returns the "Author" who signed a message from the outpoint at index i.
	relScript := authorTx.TxOut[i].PkScript
	// This pubkeyscript defines the author of the post

	scriptClass, addrs, _, err := btcscript.ExtractPkScriptAddrs(relScript, params)
	if err != nil {
		return "", err
	}
	if scriptClass != btcscript.PubKeyHashTy {
		return "", fmt.Errorf("Author script is not p2pkh")
	}
	// We know that the returned value is a P2PKH; therefore it must have
	// one address which is the author of the attached bulletin
	author := addrs[0].String()

	return author, nil

}

func (bltn *Bulletin) Bytes() ([]byte, error) {
	// Takes a bulletin and converts into a byte array. A bulletin has two
	// components. The leading 8 magic bytes and then the serialized protocol
	// buffer that contains the real message 'payload'.
	payload := make([]byte, 0)

	wireb := &WireBulletin{
		Version: proto.Uint32(ProtocolVersion),
		Topic:   proto.String(bltn.Topic),
		Message: proto.String(bltn.Message),
	}

	pbytes, err := proto.Marshal(wireb)
	if err != nil {
		return payload, err
	}

	payload = append(payload, Magic[:]...)
	payload = append(payload, pbytes...)
	return payload, nil
}

func (bltn *Bulletin) NumOuts() (int, error) {
	// Returns the number of txouts needed for this bulletin

	rawbytes, err := bltn.Bytes()
	if err != nil {
		return 0, err
	}

	numouts := len(rawbytes) / 20
	if len(rawbytes)%20 != 0 {
		numouts += 1
	}

	return numouts, nil
}
