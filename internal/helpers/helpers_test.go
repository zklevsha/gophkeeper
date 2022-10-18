package helpers

import "testing"

func TestGetRandomStr(t *testing.T) {
	wantSize := 32
	testStr := GetRandomSrt(wantSize)
	haveSize := len(testStr)
	if wantSize != haveSize {
		t.Errorf("string lenght mismatch: want: %d, have: %d",
			wantSize, haveSize)
	}
}
