package typeddata

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	signercore "github.com/ethereum/go-ethereum/signer/core/apitypes"
)

// x19 to avoid collision with rlp encode. x01 version byte defined in EIP-191.
var messagePadding = []byte{0x19, 0x01}

func encodeData(typed TypedData) (rst common.Hash, err error) {
	domainSeparator, err := hashStruct(eip712Domain, typed.Domain, typed.Types)
	if err != nil {
		return rst, err
	}

	primary, err := hashStruct(typed.PrimaryType, typed.Message, typed.Types)
	if err != nil {
		return rst, err
	}

	return crypto.Keccak256Hash(messagePadding, domainSeparator[:], primary[:]), nil
}

func encodeDataV4(typedData signercore.TypedData) ([]byte, error) {
	domainSeparator, err := typedData.HashStruct("EIP712Domain", typedData.Domain.Map())
	if err != nil {
		return nil, err
	}

	typedDataHash, err := typedData.HashStruct(typedData.PrimaryType, typedData.Message)
	if err != nil {
		return nil, err
	}

	rawData := []byte(fmt.Sprintf("\x19\x01%s%s", string(domainSeparator), string(typedDataHash)))
	sighash := crypto.Keccak256(rawData)

	return sighash, nil
}

func HashTypedDataV4(typedData signercore.TypedData) (common.Hash, error) {
	hashBytes, err := encodeDataV4(typedData)
	if err != nil {
		return common.Hash{}, err
	}

	return common.BytesToHash(hashBytes), nil
}

func SignHashTypedDataV4(sighash common.Hash, prv *ecdsa.PrivateKey) (hexutil.Bytes, error) {
	sig, err := crypto.Sign(sighash[:], prv)
	if err != nil {
		return nil, err
	}

	sig[64] += 27

	return sig, nil
}

func SignTypedDataV4(typedData signercore.TypedData, prv *ecdsa.PrivateKey) (hexutil.Bytes, error) {
	sighash, err := HashTypedDataV4(typedData)
	if err != nil {
		return nil, err
	}

	return SignHashTypedDataV4(sighash, prv)
}

// Sign TypedData with a given private key. Verify that chainId in the typed data matches currently selected chain.
func Sign(typed TypedData, prv *ecdsa.PrivateKey, chain *big.Int) ([]byte, error) {
	hash, err := ValidateAndHash(typed, chain)
	if err != nil {
		return nil, err
	}

	sig, err := crypto.Sign(hash[:], prv)
	if err != nil {
		return nil, err
	}

	sig[64] += 27

	return sig, nil
}

func Verify(hash common.Hash, sig []byte, sender common.Address) bool {
	pubKey, err := crypto.SigToPub(hash.Bytes(), sig)
	if err != nil {
		return false
	}

	return crypto.PubkeyToAddress(*pubKey) == sender
}
