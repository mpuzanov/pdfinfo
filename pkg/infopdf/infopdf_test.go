package infopdf_test

import (
	pi "pdfinfo/pkg/infopdf"
	"testing"
)

func TestGetPageCountPdf(t *testing.T) {
	fileName := "./../../example/test.pdf"
	want := 2
	got, err := pi.GetPageCountPdf(fileName)
	if err != nil {
		t.Errorf("got=%v, expected=%v, error: %v", got, want, err)
	}
	if got != want {
		t.Errorf("got=%v, expected=%v", got, want)
	}
}
