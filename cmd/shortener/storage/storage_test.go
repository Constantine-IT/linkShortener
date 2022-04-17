package storage

import (
	"testing"
)

func TestInsert(t *testing.T) {
	type want struct {
		shortURL string
		longURL  string
	}

	tests := []struct {
		name     string
		shortURL string
		longURL  string
		want     want
	}{
		{name: "Data methods test #1: POST longURL then GET with right shortURL",
			shortURL: "SDFGHJK",
			longURL:  "http://test.test/test1",
			want: want{
				shortURL: "SDFGHJK",
				longURL:  "http://test.test/test1",
			},
		},
		{name: "Data methods test #2: POST longURL then GET with wrong shortURL",
			shortURL: "QWERTYU",
			longURL:  "http://test.test/test1",
			want: want{
				shortURL: "UYTREWQ",
				longURL:  "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := NewStorage()
			err := storage.Insert(tt.shortURL, tt.longURL, "testUSER", "")
			if err != nil {
				t.Errorf("Error in INSERT method: %s", err.Error())
			}
			gotLongURL, _, gotFlag := storage.Get(tt.want.shortURL)

			if (gotFlag == true) && (gotLongURL != tt.want.longURL) {
				t.Errorf("GET return longURL = %v, but want %v", gotLongURL, tt.want.longURL)
			}
			if (gotFlag == false) && (gotLongURL != tt.want.longURL) {
				t.Errorf("GET with not existing shortURL, then longURL is empty, but want %v", tt.want.longURL)
			}
		})
	}
}