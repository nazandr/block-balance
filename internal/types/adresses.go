package types

import (
	"math/big"
	"sync"
)

type Addresses struct {
	mx        *sync.RWMutex
	addresses map[string]*big.Int
}

func NewAddresses() *Addresses {
	return &Addresses{
		mx:        &sync.RWMutex{},
		addresses: make(map[string]*big.Int),
	}
}

func (a *Addresses) Store(from, to string, value *big.Int) {
	a.mx.Lock()
	defer a.mx.Unlock()
	if v, ok := a.addresses[from]; ok {
		a.addresses[from] = big.NewInt(0).Set(value).Sub(v, value)
	} else {
		a.addresses[from] = big.NewInt(0).Neg(value)
	}

	if v, ok := a.addresses[to]; ok {
		a.addresses[to] = big.NewInt(0).Set(value).Add(v, value)
	} else {
		a.addresses[to] = big.NewInt(0).Set(value)
	}
}

func (a *Addresses) FindLargest() (string, *big.Int) {
	a.mx.RLock()
	defer a.mx.RUnlock()

	type account struct {
		address string
		changes *big.Int
	}
	res := account{
		address: "",
		changes: big.NewInt(0),
	}
	for add, val := range a.addresses {
		if i := val.CmpAbs(res.changes); i == +1 {
			res = account{
				address: add,
				changes: val.Abs(val),
			}
		}
	}

	return res.address, res.changes
}
