package plagiarism

import (
	"bufio"
	"fmt"
	"strings"
)

const (
	// N default n-gram size
	N = 8
	// LANG default language
	LANG = "en"
)

// Set of n-grams and scores
type Set struct {
	NGram []string
	Score int
}

// Detector struct
type Detector struct {
	N               int
	Lang            string
	StopWords       []string
	SourceText      string
	TargetText      string
	SourceStopWords []string
	TargetStopWords []string
	SourceNGrams    [][]string
	TargetNGrams    [][]string
	Score           float64
	Similar         int
	Total           int
}

// NewDetector implements the detector interface. Will return a new detector or an error
// if any of the optional arguments fails.
func NewDetector(options ...Option) (*Detector, error) {
	// implement a new detector interface with defaults
	detector := &Detector{N: N, Lang: LANG, StopWords: StopWords[LANG].([]string)}

	// iterrate over options, apply or return an error on failure
	for _, opt := range options {
		if err := opt(detector); err != nil {
			return nil, err
		}
	}
	// retrun the detecor
	return detector, nil
}

// Option applies detector options and returns an error on failure.
type Option func(*Detector) (err error)

// SetN option sets the detector's n-gram size and must be a positive integer larger than 0,
// otherwise an error will be returned. The default n-gram size is 8.
func SetN(n int) Option {
	return func(p *Detector) (err error) {
		// check if n-gram size is larger than 0
		if n > 0 {
			p.N = n
			return
		}
		// otherwise return an error
		return fmt.Errorf("illegal n-gram size %d, must be a positive integer larger than 0 (tip consider using values within range 7-16)", n)
	}
}

// SetLang option sets the detector's language as well as the stopwords for the specified language.
// Use ISO 639-1 formatted language codes (https://en.wikipedia.org/wiki/List_of_ISO_639-1_codes).
// Refer to stopwords.go for all supported languages. If the specified language doesn't exists or
// has no stopwords (empty []string), will return an error. If you want to use a custom
// language or a custom stopwords list use SetStopWords option instead.
func SetLang(lang string) Option {
	return func(p *Detector) (err error) {
		// check if language exists and has stopwords
		if val, ok := StopWords[lang]; ok && val != nil {
			p.Lang = lang
			p.StopWords = val.([]string)
			return
		}
		// otherwise return an error
		return fmt.Errorf("language %s not found or not supported yet (tip consider using a custom stopwords list with SetStopWords option", lang)
	}
}

// SetStopWords option will set a custom language and stopword list.
func SetStopWords(stopWords []string) Option {
	return func(p *Detector) (err error) {
		// check if stopwords list is not empty, otherwise return an error
		if len(stopWords) < 1 {
			return fmt.Errorf("stopwords list cannot be empty")
		}
		p.Lang = "custom"
		p.StopWords = stopWords
		return
	}
}

// Tokenize method will split the input string using bufio.Scanner into word tokens
// in order to filter out the unnecessary words. You can always use your own
// tokenizer and provide only the stopwords by using the SetStopWords option instead.
func (p *Detector) Tokenize(input string) []string {
	var output []string
	scanner := bufio.NewScanner(strings.NewReader(strings.ToLower(input)))
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		output = append(output, scanner.Text())
	}

	return output
}

// GetStopWords returns the stopwords list for a given token list.
func (p *Detector) GetStopWords(input []string) []string {
	var output []string
	for _, token := range input {
		if p.IsStopWord(token) {
			output = append(output, token)
		}
	}
	return output
}

// IsStopWord will check if a given token is in the stopwords list.
func (p *Detector) IsStopWord(token string) bool {
	for _, stopWord := range p.StopWords {
		if stopWord == token {
			return true
		}
	}
	return false
}

// GetNGrams returns the 2D array representation of the stopword list.
func (p *Detector) GetNGrams(tokens []string) [][]string {
	// implement ngram 2D list
	grams := make([][]string, 0)

	// calculate offset and max for N, length
	offset := int(p.N / 2)
	max := len(tokens)

	// loop through tokens and append to ngram list
	for i := range tokens {
		if i < offset || i+p.N-offset > max {
			continue
		}
		grams = append(grams, tokens[i-offset:i+p.N-offset])
	}

	// return the n-gram list
	return grams
}

// DeepEquaility something like Jaccard coefficient but with strict position.
// Instead of intersection / union we use position / union == 1.0
func (p *Detector) DeepEquaility(source, target *[][]string) [][]Set {
	// Copy Temp Slices, I > J
	tempI := *source
	tempJ := *target

	// initilize source sets and set scores to 0
	setI := make([]Set, len(tempI))
	for i := range tempI {
		setI[i] = Set{NGram: tempI[i], Score: 0}
	}

	// initilize target sets and set scores to 0
	setJ := make([]Set, len(tempJ))
	for j := range tempJ {
		setJ[j] = Set{NGram: tempJ[j], Score: 0}
	}

	// find equals for I/J and set score to 1
	for i := range setI {
		for j := range setJ {
			if p.Equal(setI[i].NGram, setJ[j].NGram) {
				setI[i].Score = 1
				setJ[j].Score = 1
			}
		}
	}

	// return the sets
	return [][]Set{setI, setJ}
}

// Equal will test if Slices are Equal (NxN).
func (p *Detector) Equal(source, target []string) bool {
	for i := range source {
		if source[i] != target[i] {
			return false
		}
	}
	return true
}

// DetectWithStrings returns an error on failure, otherwise will invoke
// DetectWithStopWords method.
func (p *Detector) DetectWithStrings(source, target string) error {
	// check if any of source or target text is an empty string and return an error
	if source == "" || target == "" {
		return fmt.Errorf("both, source and target text cannot be empty")
	}

	// assign detector values
	p.SourceText = source
	p.TargetText = target

	// tokenize sting, filter stopwords and return DetectWithStopWords method
	return p.DetectWithStopWords(
		p.GetStopWords(p.Tokenize(p.SourceText)),
		p.GetStopWords(p.Tokenize(p.TargetText)),
	)
}

// DetectWithStopWords returns an error on failure, otherwise will set Score, Similar and Total
// values to the detector interface.
func (p *Detector) DetectWithStopWords(source, target []string) error {
	// check if any of source or target stopwords list is an empty string array and return an error
	if len(source) < 1 || len(target) < 1 {
		return fmt.Errorf("both, source and target stopwords list cannot be empty")
	}

	// assign detector values
	p.SourceStopWords = source
	p.TargetStopWords = target

	// get the n-grams and assign detector values
	p.SourceNGrams = p.GetNGrams(p.SourceStopWords)
	p.TargetNGrams = p.GetNGrams(p.TargetStopWords)

	// test n-grams equality
	equality := p.DeepEquaility(&p.SourceNGrams, &p.TargetNGrams)

	// sum source similarity score
	for i := range equality[0] {
		p.Similar += equality[0][i].Score
	}

	// sum target similarity score
	for j := range equality[1] {
		p.Similar += equality[1][j].Score
	}

	// sum totals
	p.Total = len(equality[0]) + len(equality[1])

	// calculate probability
	p.Score = float64(p.Similar) / float64(p.Total)

	return nil
}
