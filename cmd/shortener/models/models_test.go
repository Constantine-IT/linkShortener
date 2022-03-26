package models

import (
	"testing"
)

func TestInsert(t *testing.T) {
	type want struct {
		shortURL string
		longURL  string
		flag     bool
	}
	tests := []struct {
		name  string
		inURL string
		want  want
	}{
		{name: "Simple insert test #1",
			inURL: "http://tudzqakmoorcb.net/bflsgr36aqo4x6/mmktfboj8",
			want: want{
				shortURL: "ISQJSGDNXNSG",
				longURL:  "http://tudzqakmoorcb.net/bflsgr36aqo4x6/mmktfboj8",
				flag:     true,
			},
		},
		{name: "Simple insert test #2",
			inURL: "http://commemns.edu/374hwsjhsdhh/mmktfboj8/sdejhjwh",
			want: want{
				shortURL: "LKJSDJJSNNDJDFDJC",
				longURL:  "http://abirvalg.shop/joowjubwjejjw/tfboj8/sdejhjwh1111111",
				flag:     false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Insert(tt.want.shortURL, tt.inURL)
			gotLongURL, gotFlag := Get(tt.want.shortURL)

			if (gotLongURL != tt.want.longURL) && (gotFlag == tt.want.flag) {
				t.Errorf("Get() longURL = %v, want %v", gotLongURL, tt.want.longURL)
			}
			if (gotLongURL == tt.want.longURL) && (gotFlag != tt.want.flag) {
				t.Errorf("Get() longURL the same, but gotFlag = %v, want %v", gotFlag, tt.want.flag)
			}
		})
	}
}
