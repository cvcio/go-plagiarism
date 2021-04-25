package main

import (
	"fmt"

	"github.com/cvcio/go-plagiarism"
)

var source = `Plagiarism detection using stopwords n-grams. go-plagiarism is the main algorithm 
that utilizes MediaWatch and is inspired by Efstathios Stamatatos paper. 
We only rely on a list of stopwords to calculate 
the plagiarism probability between two texts, in combination with n-gram 
loops that let us find, not only plagiarism but also 
paraphrase and patchwork plagiarism. In our case (cc MediaWatch) we 
use this algorithm to create relationships between similar articles and 
map the process, or the chain of misinformation. As our 
scope is to track propaganda networks in the news ecosystem, 
this algorithm only tested in such context.`

var target = `We only rely on a list of stopwords to calculate 
the plagiarism probability between two texts, in combination with n-gram 
loops that let us find, not only plagiarism but also 
paraphrase and patchwork plagiarism. In our case (cc MediaWatch) we 
use this algorithm to create relationships between similar articles and 
map the process, or the chain of misinformation. As our 
scope is to track propaganda networks in the news ecosystem, 
this algorithm only tested in such context.`

func main() {
	detector, _ := plagiarism.NewDetector()
	err := detector.DetectWithStrings(source, target)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Probability: %.2f, Similar n-grams %d, Total n-grams %d\n", detector.Score, detector.Similar, detector.Total)
}
