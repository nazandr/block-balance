package utils

import (
	"math/big"
	"reflect"
	"testing"
)

func TestHexToInt(t *testing.T) {
	bn, ok := big.NewInt(0).SetString("25599999999999999999999999999999999999766", 10)
	if !ok {
		t.Error("wrong bigint")
	}
	type args struct {
		hex string
	}
	tests := []struct {
		name    string
		args    args
		want    *big.Int
		wantErr bool
	}{
		{
			name: "big",
			args: args{
				hex: "0x4B3B4CA85A86C47A098A223FFFFFFFFF16",
			},
			want:    bn,
			wantErr: false,
		},
		{
			name: "zero",
			args: args{
				hex: "0x0",
			},
			want:    big.NewInt(0),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HexToBigInt(tt.args.hex)
			if (err != nil) != tt.wantErr {
				t.Errorf("HexToInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HexToInt() = %v, want %v", got, tt.want)
			}
		})
	}
}
