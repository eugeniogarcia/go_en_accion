package text

import "testing"

func TestCountCharacters(t *testing.T) {
	total, err := CountCharacters("testdata/sample1.txt")
	if err != nil {
		t.Error("Unexpected error:", err) // no deberia fallar
	}
	if total != 35 {
		t.Error("Expected 35, got", total) // deberia contar 35 caracteres
	}
	_, err = CountCharacters("testdata/no_file.txt")
	if err == nil {
		t.Error("Expected an error") // deberia fallar
	}
}
