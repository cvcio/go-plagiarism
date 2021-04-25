package plagiarism

import (
	"io/ioutil"
	"testing"
)

func Test_NewDetector(t *testing.T) {
	_, err := NewDetector()
	if err != nil {
		t.Errorf("Error while creating detector: %s", err)
	}
}

func Test_NewDetectorSetN(t *testing.T) {
	var tests = []struct {
		n        int
		expected int
	}{
		{1, 1},
		{10, 10},
		{100, 100},
	}

	for _, test := range tests {
		detector, _ := NewDetector(SetN(test.n))
		if detector.N != test.expected {
			t.Errorf("Error while setting n-gram size %d, expected %d, got %d", test.n, test.expected, detector.N)
		}
	}
}

func Test_NewDetectorSetNError(t *testing.T) {
	var tests = []struct {
		n int
	}{
		{0},
		{-1},
		{},
	}

	for _, test := range tests {
		detector, err := NewDetector(SetN(test.n))
		if err == nil {
			t.Errorf("Error while setting n-gram size %d, expected Error, got %d", test.n, detector.N)
		}
	}
}

func Test_NewDetectorSetLang(t *testing.T) {
	var tests = []struct {
		lang     string
		expected string
	}{
		{"bg", "bg"},
		{"hr", "hr"},
		{"nl", "nl"},
		{"en", "en"},
		{"fi", "fi"},
		{"fr", "fr"},
		{"de", "de"},
		{"el", "el"},
		{"hu", "hu"},
		{"it", "it"},
		{"no", "no"},
		{"pl", "pl"},
		{"pt", "pt"},
		{"ro", "ro"},
		{"ru", "ru"},
		{"tr", "tr"},
		{"uk", "uk"},
	}

	for _, test := range tests {
		detector, _ := NewDetector(SetLang(test.lang))
		if detector.Lang != test.expected {
			t.Errorf("Error while setting lang %s, expected %s, got %s", test.lang, test.expected, detector.Lang)
		}
	}
}

func Test_NewDetectorSetLangError(t *testing.T) {
	var tests = []struct {
		lang string
	}{
		{"xx"},
		{""},
	}

	for _, test := range tests {
		detector, err := NewDetector(SetLang(test.lang))
		if err == nil {
			t.Errorf("Error while setting lang %s, expected Error, got %s", test.lang, detector.Lang)
		}
	}
}

func Test_NewDetectorWithStrings(t *testing.T) {
	var tests = []struct {
		lang    string
		source  string
		target  string
		Score   float64
		Similar int
		Total   int
	}{
		{"el", "examples/testfiles/el/t1-source.txt", "examples/testfiles/el/t1-source.txt", 1.0, 528, 528},
		{"el", "examples/testfiles/el/t1-source.txt", "examples/testfiles/el/t1-target.txt", 0.5544041450777202, 214, 386},
		{"el", "examples/testfiles/el/t2-source.txt", "examples/testfiles/el/t2-target.txt", 0.0, 0, 0},
		{"el", "examples/testfiles/el/t3-source.txt", "examples/testfiles/el/t3-target.txt", 0.34255129348795715, 384, 1121},
		{"el", "examples/testfiles/el/test-a.txt", "examples/testfiles/el/test-b.txt", 0.41025641025641024, 48, 117},
		{"el", "examples/testfiles/el/small-a.txt", "examples/testfiles/el/small-b.txt", 0.0, 0, 2},
	}

	for _, test := range tests {
		detector, _ := NewDetector(SetLang(test.lang))

		source, _ := readFile(test.source)
		target, _ := readFile(test.target)

		detector.DetectWithStrings(source, target)

		if detector.Score != test.Score {
			t.Errorf("Error in DetectWithStrings, expected %f, got %f", test.Score, detector.Score)
		}

		if detector.Similar != test.Similar {
			t.Errorf("Error in DetectWithStrings, expected %d, got %d", test.Similar, detector.Similar)
		}

		if detector.Total != test.Total {
			t.Errorf("Error in DetectWithStrings, expected %d, got %d", test.Total, detector.Total)
		}
	}
}

func readFile(path string) (string, error) {
	content, err := ioutil.ReadFile(path)

	if err != nil {
		return "", err
	}

	return string(content), nil
}
