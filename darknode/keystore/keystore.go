package keystore

import (
	"encoding/json"
	"fmt"
	"math/big"
)

type Keystore struct {
	Ecdsa `json:"ecdsa"`
	Rsa   `json:"rsa"`
}

func RandomKeystore() (Keystore, error) {
	ecdsaPrivkey, err := RandomEcdsaPrivKey()
	if err != nil {
		return Keystore{}, fmt.Errorf("error generating ecdsa privkey: %v", err)
	}

	rsaPrivkey, err := RandomRsaPrivKey()
	if err != nil {
		return Keystore{}, fmt.Errorf("error generating rsa privkey: %v", err)
	}

	return Keystore{ecdsaPrivkey, rsaPrivkey}, nil
}

func unmarshalStringFromMap(m map[string]json.RawMessage, k string) (string, error) {
	if val, ok := m[k]; ok {
		str := ""
		if err := json.Unmarshal(val, &str); err != nil {
			return "", err
		}
		return str, nil
	}
	return "", fmt.Errorf("%s is nil", k)
}

func unmarshalIntFromMap(m map[string]json.RawMessage, k string) (int, error) {
	if val, ok := m[k]; ok {
		i := 0
		if err := json.Unmarshal(val, &i); err != nil {
			return 0, err
		}
		return i, nil
	}
	return 0, fmt.Errorf("%s is nil", k)
}

func unmarshalBigIntFromMap(m map[string]json.RawMessage, k string) (*big.Int, error) {
	if val, ok := m[k]; ok {
		bytes := []byte{}
		if err := json.Unmarshal(val, &bytes); err != nil {
			return nil, err
		}
		return big.NewInt(0).SetBytes(bytes), nil
	}
	return nil, fmt.Errorf("%s is nil", k)
}

func unmarshalBigIntsFromMap(m map[string]json.RawMessage, k string) ([]*big.Int, error) {
	bigInts := []*big.Int{}
	if val, ok := m[k]; ok {
		vals := []json.RawMessage{}
		if err := json.Unmarshal(val, &vals); err != nil {
			return bigInts, err
		}
		for _, val := range vals {
			bytes := []byte{}
			if err := json.Unmarshal(val, &bytes); err != nil {
				return bigInts, err
			}
			bigInts = append(bigInts, big.NewInt(0).SetBytes(bytes))
		}
		return bigInts, nil
	}
	return bigInts, fmt.Errorf("%s is nil", k)
}
