
[![Language](https://img.shields.io/badge/Language-Go-blue.svg)](https://golang.org/)
[![Build Status](https://github.com/cvcio/go-plagiarism/workflows/Go/badge.svg)](https://github.com/cvcio/go-plagiarism/actions)
[![GoDoc](https://pkg.go.dev/badge/github.com/cvcio/go-plagiarism)](https://pkg.go.dev/github.com/cvcio/go-plagiarism)
[![Go Report Card](https://goreportcard.com/badge/github.com/cvcio/go-plagiarism)](https://goreportcard.com/report/github.com/cvcio/go-plagiarism)

# Plagiarism detection using stopwords *n*-grams

`go-plagiarism` is the main algorithm that utilizes [MediaWatch](https://mediawatch.io) and is inspired by [Efstathios Stamatatos](https://www3.icsd.aegean.gr/lecturers/stamatatos/) paper [Plagiarism detection using stopwords *n*-grams](http://dx.doi.org/10.1002/asi.21630).

![Plagiarism Detection Algorithm - Function Words - Tokens](https://github.com/cvcio/go-plagiarism/raw/main/assets/Plagiarism%20Detection%20Algorithm%20-%20Function%20Words%20-%20Tokens.png)
<p align="center">filter stopwords</p>

![N-Grams](https://github.com/cvcio/go-plagiarism/raw/main/assets/N-Grams.png)
<p align="center">loop through n-grams</p>

We only rely on a small list of stopwords to calculate the plagiarism probability between two texts, in combination with *n*-gram loops that let us find, not only plagiarism but also paraphrase and patchwork plagiarism. In our case (cc [MediaWatch](https://mediawatch.io)) we use this algorithm to create relationships between similar articles and map the process, or **the chain of misinformation**. As our scope is to track propaganda networks in the news ecosystem, this algorithm only tested in such context.

![The Chain of Misinformation](https://github.com/cvcio/go-plagiarism/raw/main/assets/The%20Chain%20of%20Misinformation.png)
<p align="center">The Chain of Misinformation</p>

![Similarity Network](https://github.com/cvcio/go-plagiarism/raw/main/assets/Similarity%20Network.png)
<p align="center">Similarity Network</p>

## Usage

```bash
go get github.com/cvcio/go-plagiarism
```

```go
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

// > Probability: 0.91, Similar n-grams 72, Total n-grams 79
```
## Options

Detector can be initialized with options, `SetN` to set the *n*-gram size, `SetLang` to set the detector's language model and assign the approrpiate stopwords and `SetStopWords` to assign a custom list of stopwords. Do not use `SetLang` alongside with `SetStopWords` as it will override one another.
```go
plagiarism.SetN(n int) Option // will set the desired n-gram size
plagiarism.SetLang(lang string) Option // will set the detector's language and assign the default stopwords
plagiarism.SetStopWords(stopWords []string) Option // will set a custom list of stopwords as the default
```

To use the detector with options, simple pass the options during initialization.
```go
// create a detector with 12 N n-gram size and set the language to Greek
detector, err := plagiarism.NewDetector(plagiarism.SetN(12), plagiarism.SetLang("el"))
```

```go
// create a detector with default n-gram size (8) and set a custom stopword list
detector, err := plagiarism.NewDetector(plagiarism.SetStopWords([]string{"ο", "του", "η", "της", "αλλά"}))
```

## Test Coverage
```bash
go test -v
```
## Contributing

If you're new to contributing to Open Source on Github, [this guide](https://opensource.guide/how-to-contribute/) can help you get started. Please check out the contribution guide for more details on how issues and pull requests work. Before contributing be sure to review the [code of conduct](/CODE_OF_CONDUCT.md).

<a href="https://github.com/cvcio/go-plagiarism/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=cvcio/go-plagiarism" />
</a>

## license

This library is distributed under the MIT license found in the [LICENSE](/LICENSE) file.