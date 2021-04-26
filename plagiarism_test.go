package plagiarism

import (
	"io/ioutil"
	"reflect"
	"testing"
)

var sourceString = `Plagiarism detection using stopwords n-grams. go-plagiarism is the main algorithm 
that utilizes MediaWatch and is inspired by Efstathios Stamatatos paper. 
We only rely on a list of stopwords to calculate 
the plagiarism probability between two texts, in combination with n-gram 
loops that let us find, not only plagiarism but also 
paraphrase and patchwork plagiarism. In our case (cc MediaWatch) we 
use this algorithm to create relationships between similar articles and 
map the process, or the chain of misinformation. As our 
scope is to track propaganda networks in the news ecosystem, 
this algorithm only tested in such context.`

var sourceStopWords = []string{
	"using", "is", "the", "that", "and", "is", "by", "we", "only", "on", "a", "of", "to", "the", "between", "two", "in", "with", "that", "let", "us", "not", "only", "but", "also", "and", "in", "our", "case", "we", "use", "this", "to", "between", "similar", "and", "the", "or", "the", "of", "as", "our", "is", "to", "in", "the", "this", "only", "in", "such",
}

var targetString = `We only rely on a list of stopwords to calculate 
the plagiarism probability between two texts, in combination with n-gram 
loops that let us find, not only plagiarism but also 
paraphrase and patchwork plagiarism. In our case (cc MediaWatch) we 
use this algorithm to create relationships between similar articles and 
map the process, or the chain of misinformation. As our 
scope is to track propaganda networks in the news ecosystem, 
this algorithm only tested in such context.`

var targetStopWords = []string{
	"we", "only", "on", "a", "of", "to", "the", "between", "two", "in", "with", "that", "let", "us", "not", "only", "but", "also", "and", "in", "our", "case", "we", "use", "this", "to", "between", "similar", "and", "the", "or", "the", "of", "as", "our", "is", "to", "in", "the", "this", "only", "in", "such",
}

func Test_NewDetector(t *testing.T) {
	_, err := NewDetector()
	if err != nil {
		t.Fatalf("Error while creating detector: %s", err)
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

func Test_NewDetectorDetectError(t *testing.T) {
	detector := &Detector{}

	detector.N = 8
	detector.Lang = "en"
	detector.StopWords = StopWords[detector.Lang].([]string)

	err := detector.Detect()

	if err == nil {
		t.Fatalf("Error in Detect, expected Error, got nil")
	}
}
func Test_NewDetectorDetectStopWords(t *testing.T) {
	detector := &Detector{}

	detector.N = 8
	detector.Lang = "en"
	detector.StopWords = StopWords[detector.Lang].([]string)
	detector.SourceStopWords = sourceStopWords
	detector.TargetStopWords = targetStopWords

	err := detector.Detect()

	if err != nil {
		t.Fatalf("Error while creating detector: %s", err.Error())
	}
	if detector.Score != 0.9113924050632911 {
		t.Errorf("Error in DetectWithStrings, expected %f, got %f", 0.9113924050632911, detector.Score)
	}

	if detector.Similar != 72 {
		t.Errorf("Error in DetectWithStrings, expected %d, got %d", 72, detector.Similar)
	}

	if detector.Total != 79 {
		t.Errorf("Error in DetectWithStrings, expected %d, got %d", 79, detector.Total)
	}
}
func Test_NewDetectorDetectStrings(t *testing.T) {
	detector := &Detector{}

	detector.N = 8
	detector.Lang = "en"
	detector.StopWords = StopWords[detector.Lang].([]string)
	detector.SourceText = sourceString
	detector.TargetText = targetString

	err := detector.Detect()

	if err != nil {
		t.Fatalf("Error while creating detector: %s", err.Error())
	}
	if detector.Score != 0.9113924050632911 {
		t.Errorf("Error in DetectWithStrings, expected %f, got %f", 0.9113924050632911, detector.Score)
	}

	if detector.Similar != 72 {
		t.Errorf("Error in DetectWithStrings, expected %d, got %d", 72, detector.Similar)
	}

	if detector.Total != 79 {
		t.Errorf("Error in DetectWithStrings, expected %d, got %d", 79, detector.Total)
	}
}

func Test_NewDetectorWithStringsWithStruct(t *testing.T) {
	detector := &Detector{}

	detector.N = 8
	detector.Lang = "en"
	detector.StopWords = StopWords[detector.Lang].([]string)

	err := detector.DetectWithStrings(sourceString, targetString)

	if err != nil {
		t.Fatalf("Error while creating detector: %s", err.Error())
	}
	if detector.Score != 0.9113924050632911 {
		t.Errorf("Error in DetectWithStrings, expected %f, got %f", 0.9113924050632911, detector.Score)
	}

	if detector.Similar != 72 {
		t.Errorf("Error in DetectWithStrings, expected %d, got %d", 72, detector.Similar)
	}

	if detector.Total != 79 {
		t.Errorf("Error in DetectWithStrings, expected %d, got %d", 79, detector.Total)
	}
}
func Test_NewDetectorWithString(t *testing.T) {
	detector, _ := NewDetector()
	err := detector.DetectWithStrings(sourceString, targetString)

	if err != nil {
		t.Fatalf("Error while creating detector: %s", err.Error())
	}
	if detector.Score != 0.9113924050632911 {
		t.Errorf("Error in DetectWithStrings, expected %f, got %f", 0.9113924050632911, detector.Score)
	}

	if detector.Similar != 72 {
		t.Errorf("Error in DetectWithStrings, expected %d, got %d", 72, detector.Similar)
	}

	if detector.Total != 79 {
		t.Errorf("Error in DetectWithStrings, expected %d, got %d", 79, detector.Total)
	}
}

func Test_NewDetectorWithStopWords(t *testing.T) {
	detector, _ := NewDetector()
	err := detector.DetectWithStopWords(sourceStopWords, targetStopWords)

	if err != nil {
		t.Fatalf("Error while creating detector: %s", err.Error())
	}
	if detector.Score != 0.9113924050632911 {
		t.Errorf("Error in DetectWithStrings, expected %f, got %f", 0.9113924050632911, detector.Score)
	}

	if detector.Similar != 72 {
		t.Errorf("Error in DetectWithStrings, expected %d, got %d", 72, detector.Similar)
	}

	if detector.Total != 79 {
		t.Errorf("Error in DetectWithStrings, expected %d, got %d", 79, detector.Total)
	}
}

func Test_NewDetectorWithStringsMany(t *testing.T) {
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

func TestDeepEquaility(t *testing.T) {
	type args struct {
		a *[][]string
		b *[][]string
		n int
	}
	tests := []struct {
		name string
		args args
		want [][]Set
	}{
		// TODO: Add test cases.
	}

	detector, _ := NewDetector()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := detector.DeepEquaility(tt.args.a, tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeepEquaility() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEqual(t *testing.T) {
	type args struct {
		a []string
		b []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}

	detector, _ := NewDetector()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := detector.Equal(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNGrams(t *testing.T) {
	type args struct {
		tokens []string
		n      int
	}
	tests := []struct {
		name string
		args args
		want [][]string
	}{
		{
			name: "Small Array",
			args: args{
				tokens: []string{"Ολοι", "υποψιαζόμαστε", "και", "ξέρουμε"},
				n:      1,
			},
			want: [][]string{
				{"Ολοι"},
				{"υποψιαζόμαστε"},
				{"και"},
				{"ξέρουμε"},
			},
		},
		{
			name: "Small Array",
			args: args{
				tokens: []string{"Ολοι", "υποψιαζόμαστε", "και", "ξέρουμε"},
				n:      2,
			},
			want: [][]string{
				{"Ολοι", "υποψιαζόμαστε"},
				{"υποψιαζόμαστε", "και"},
				{"και", "ξέρουμε"},
			},
		},
		{
			name: "Small Array",
			args: args{
				tokens: []string{"Ολοι", "υποψιαζόμαστε", "και", "ξέρουμε"},
				n:      3,
			},
			want: [][]string{
				{"Ολοι", "υποψιαζόμαστε", "και"},
				{"υποψιαζόμαστε", "και", "ξέρουμε"},
			},
		},
		{
			name: "Small Array",
			args: args{
				tokens: []string{"Ολοι", "υποψιαζόμαστε", "και", "ξέρουμε"},
				n:      4,
			},
			want: [][]string{
				{"Ολοι", "υποψιαζόμαστε", "και", "ξέρουμε"},
			},
		},
		// TODO: Add test cases.
	}

	for _, tt := range tests {
		detector, _ := NewDetector(SetN(tt.args.n))
		t.Run(tt.name, func(t *testing.T) {
			if got := detector.GetNGrams(tt.args.tokens); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NGrams() = %v, want %v", got, tt.want)
			}
		})
	}
}

func readFile(path string) (string, error) {
	content, err := ioutil.ReadFile(path)

	if err != nil {
		return "", err
	}

	return string(content), nil
}
