package ahimsa

import (
	"bytes"
	"errors"
	"fmt"
	"time"

	"code.google.com/p/goprotobuf/proto"
	"github.com/conformal/btcnet"
	"github.com/conformal/btcscript"
	"github.com/conformal/btcutil"
	"github.com/conformal/btcwire"
)

const (
	MaxBoardLen int = 30
)

var (
	Magic = [8]byte{
		0x42, 0x52, 0x45, 0x54, 0x48, 0x52, 0x45, 0x4e, /* | BRETHREN | */
	}
)

type Author string

type Bulletin struct {
	Txid      *btcwire.ShaHash
	Block     *btcwire.ShaHash
	Author    string
	Board     string
	Message   string
	Timestamp time.Time
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

// Creates a new bulletin from the containing Tx, supplied author and optional blockhash
// by unpacking txOuts that are considered data. It ignores extra junk behind the protobuffer.
// NewBulletin also asserts aspects of valid bulletins by throwing errors when msg len
// is zero or board len is greater than MaxBoardLen.
func NewBulletin(tx *btcwire.MsgTx, author string, blkhash *btcwire.ShaHash) (*Bulletin, error) {
	wireBltn := &WireBulletin{}

	// Bootleg solution, but if unmarshal fails slice txout and try again until we can try no more or it fails
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

	board := wireBltn.GetBoard()
	// assert that the length of the board is within its max size!
	if len(board) > MaxBoardLen {
		return nil, errors.New("Board length is too large.")
	}

	msg := wireBltn.GetMessage()
	// assert that the bulletin has a non zero message length.
	if len(msg) < 1 {
		return nil, errors.New("Message has no content.")
	}

	hash, _ := tx.TxSha()

	bltn := &Bulletin{
		Txid:      &hash,
		Block:     blkhash,
		Author:    author,
		Board:     board,
		Message:   msg,
		Timestamp: time.Unix(wireBltn.GetTimestamp(), 0),
	}

	return bltn, nil
}

func NewBulletinFromStr(author string, board string, msg string) (*Bulletin, error) {
	if len(board) > 30 {
		return nil, fmt.Errorf("Board too long! Length is: %d", len(board))
	}

	if len(msg) > 500 {
		return nil, fmt.Errorf("Message too long! Length is: %d", len(msg))
	}

	bulletin := Bulletin{
		Author:    author,
		Board:     board,
		Message:   msg,
		Timestamp: time.Now(),
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

func GetAuthor(tx *btcwire.MsgTx, net *btcnet.Params) (string, error) {
	// Returns the "Author" who signed the first txin of the transaction
	sigScript := tx.TxIn[0].SignatureScript

	dummyTx := btcwire.NewMsgTx()

	// Setup a script executer to parse the raw bytes of the signature script.
	script, err := btcscript.NewScript(sigScript, make([]byte, 0), 0, dummyTx, 0)
	if err != nil {
		return "", err
	}
	// Step twice due to <sig> <pubkey> format of pay 2pubkeyhash
	script.Step()
	script.Step()
	// Pull off the <pubkey>
	pkBytes := script.GetStack()[1]

	addrPubKey, err := btcutil.NewAddressPubKey(pkBytes, net)
	if err != nil {
		return "", err
	}

	return addrPubKey.EncodeAddress(), nil
}

func (bltn *Bulletin) Bytes() ([]byte, error) {
	// Takes a bulletin and converts into a byte array. A bulletin has two
	// components. The leading 8 magic bytes and then the serialized protocol
	// buffer that contains the real message 'payload'.
	payload := make([]byte, 0)

	wireb := &WireBulletin{
		Board:     proto.String(bltn.Board),
		Message:   proto.String(bltn.Message),
		Timestamp: proto.Int64(bltn.Timestamp.Unix()),
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
