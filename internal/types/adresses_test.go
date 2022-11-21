package types

import (
	"math/big"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddresses_Store(t *testing.T) {
	mx := &sync.RWMutex{}
	type fields struct {
		mx        *sync.RWMutex
		addresses map[string]*big.Int
	}
	type args struct {
		from  string
		to    string
		value *big.Int
	}
	tests := []struct {
		name   string
		fields fields
		args   []args
		want   map[string]*big.Int
	}{
		{
			name:   "2 address",
			fields: fields{mx: mx, addresses: make(map[string]*big.Int)},
			args: []args{
				{from: "a", to: "b", value: big.NewInt(100)},
				{from: "a", to: "b", value: big.NewInt(100)},
				{from: "a", to: "b", value: big.NewInt(100)}},
			want: map[string]*big.Int{"a": big.NewInt(-300), "b": big.NewInt(300)},
		},
		{
			name:   "3 address",
			fields: fields{mx: mx, addresses: make(map[string]*big.Int)},
			args: []args{
				{from: "a", to: "b", value: big.NewInt(100)},
				{from: "a", to: "b", value: big.NewInt(100)},
				{from: "b", to: "d", value: big.NewInt(100)}},
			want: map[string]*big.Int{"a": big.NewInt(-200), "b": big.NewInt(100), "d": big.NewInt(100)},
		},
		{
			name:   "4 address",
			fields: fields{mx: mx, addresses: make(map[string]*big.Int)},
			args: []args{
				{from: "a", to: "b", value: big.NewInt(100)},
				{from: "a", to: "b", value: big.NewInt(100)},
				{from: "b", to: "d", value: big.NewInt(100)},
				{from: "d", to: "c", value: big.NewInt(150)}},
			want: map[string]*big.Int{
				"a": big.NewInt(-200),
				"b": big.NewInt(100),
				"d": big.NewInt(-50),
				"c": big.NewInt(150)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Addresses{
				mx:        tt.fields.mx,
				addresses: tt.fields.addresses,
			}
			for _, v := range tt.args {
				a.Store(v.from, v.to, v.value)
			}

			assert.Equal(t, tt.want, a.addresses)
		})
	}
}

func TestAddresses_FindLargest(t *testing.T) {
	mx := &sync.RWMutex{}
	type fields struct {
		mx        *sync.RWMutex
		addresses map[string]*big.Int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
		want1  *big.Int
	}{
		{
			name: "max positive",
			fields: fields{
				mx:        mx,
				addresses: map[string]*big.Int{"a": big.NewInt(100), "b": big.NewInt(-200), "c": big.NewInt(300)},
			},
			want:  "c",
			want1: big.NewInt(300),
		},
		{
			name: "max negative",
			fields: fields{
				mx:        mx,
				addresses: map[string]*big.Int{"a": big.NewInt(100), "b": big.NewInt(-200), "c": big.NewInt(-400)},
			},
			want:  "c",
			want1: big.NewInt(400),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Addresses{
				mx:        tt.fields.mx,
				addresses: tt.fields.addresses,
			}
			got, got1 := a.FindLargest()
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.want1, got1)
		})
	}
}
