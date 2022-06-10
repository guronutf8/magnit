package Url

import (
	"github.com/stretchr/testify/require"
	"net/http"
	"strconv"
	"testing"
)

func TestFormatToIndex(t *testing.T) {
	tests := []struct {
		input uint64
	}{
		{input: 0},
		{input: 1},
		{input: 50},
		{input: 61},
		{input: 62},
		{input: 80},
		{input: 122},
		{input: 125},
		{input: 180},
		{input: 400},
		{input: 500},
		{input: 3500},
		{input: 7000},
		{input: 7500},
		{input: 7680},
		{input: 7687},
		{input: 7688},
		{input: 7700},
		{input: 8000},
		{input: 50000},
		{input: 90000},
		{input: 100000},
		{input: 150000},
		{input: 200000},
		//{input: 2000000},
	}

	for _, tc := range tests {
		t.Run("test "+strconv.FormatUint(tc.input, 10), func(t *testing.T) {
			result := FormatToShort(tc.input)
			index := FormatToIndex(result)
			require.Equal(t, int(index), int(tc.input), "no equal((")
		})
	}

}

func TestUrl_Save(t *testing.T) {
	type fields struct {
		Long string
	}
	type args struct {
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{name: "default", fields: fields{Long: "http://foo.bar/baz"}, want: 200},
		{name: "default https", fields: fields{Long: "https://foo.bar/baz"}, want: 200},
		{name: "default ftp", fields: fields{Long: "ftp://foo.bar/baz"}, want: http.StatusBadRequest},
		{name: "default no path", fields: fields{Long: "http://foo.bar/"}, want: http.StatusBadRequest},
		{name: "default no slash", fields: fields{Long: "http://foo.bar"}, want: http.StatusBadRequest},
		{name: "no domain", fields: fields{Long: "http:///"}, want: http.StatusBadRequest},
		{name: "clear url", fields: fields{Long: "http://"}, want: http.StatusBadRequest},
		{name: "clear full", fields: fields{Long: ""}, want: http.StatusBadRequest},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := Url{
				Long: tt.fields.Long,
			}
			if got := u.Save(); got != tt.want {
				t.Errorf("Save() = %v, want %v", got, tt.want)
			}
		})
	}
}
