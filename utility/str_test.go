package utility

import (
	"reflect"
	"testing"
)

func TestByteToString(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test1",
			args: args{[]byte("abcd")},
			want: "abcd",
		},
		{
			name: "test2",
			args: args{[]byte("")},
			want: "",
		},
		{
			name: "test2",
			args: args{nil},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ByteToString(tt.args.b); got != tt.want {
				t.Errorf("ByteToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringToBytes(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "test1",
			args: args{"test"},
			want: []byte("test"),
		},

		{
			name: "test2",
			args: args{""},
			want: nil,
		},

		{
			name: "test3",
			args: args{"----"},
			want: []byte("----"),
		},
		{
			name: "test3",
			args: args{"   "},
			want: []byte("   "),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringToBytes(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StringToBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}
