
[![Language](https://img.shields.io/badge/Language-Go-blue.svg)](https://golang.org/)
[![Build Status](https://github.com/cvcio/go-plagiarism/workflows/Go/badge.svg)](https://github.com/cvcio/go-plagiarism/actions)
[![GoDoc](https://pkg.go.dev/badge/github.com/cvcio/go-plagiarism)](https://pkg.go.dev/github.com/cvcio/go-plagiarism)
[![Go Report Card](https://goreportcard.com/badge/github.com/cvcio/go-plagiarism)](https://goreportcard.com/report/github.com/cvcio/go-plagiarism)

# Plagiarism detection using stopwords *n*-grams

`go-plagiarism` is the main algorithm that utilizes [MediaWatch](https://mediawatch.io) and is inspired by [Efstathios Stamatatos](https://www3.icsd.aegean.gr/lecturers/stamatatos/) paper [Plagiarism detection using stopwords *n*-grams](http://dx.doi.org/10.1002/asi.21630).

We only rely on a small list of stopwords, for each [language](#supported-languages), to calculate the plagiarism probability between two texts, in combination with *n*-grams that let us find, not only plagiarism but also paraphrase and patchwork plagiarism. Take a look at the images bellow to help you better understand the process.

During the 1st step we tokenize the strings and keep only the stopwords (red tokens) for each document, as **SourceStopWords** and **TargetStopWords**.

![Plagiarism Detection Algorithm - Function Words - Tokens](https://github.com/cvcio/go-plagiarism/raw/main/assets/Plagiarism%20Detection%20Algorithm%20-%20Function%20Words%20-%20Tokens.png)

Later we transform the stopwords for each document into *n*-grams, with default **N = 8**, and calculate the score for each set of *n*-grams.

![N-Grams](https://github.com/cvcio/go-plagiarism/raw/main/assets/N-Grams.png)

In our case (cc [MediaWatch](https://mediawatch.io)) we use this algorithm to create relationships between similar articles and map the process, or **the chain of misinformation**. As our scope is to track propaganda networks in the news ecosystem, this algorithm only tested in such context.

![The Chain of Misinformation](https://github.com/cvcio/go-plagiarism/raw/main/assets/The%20Chain%20of%20Misinformation.png)

<p align="center">The Chain of Misinformation</p>

![Similarity Network](https://github.com/cvcio/go-plagiarism/raw/main/assets/Similarity%20Network.png)

<p align="center">Similarity Network</p>

## Usage

```bash
go get github.com/cvcio/go-plagiarism
```

To use the detector you must provide either source/target texts, when using with `DetectWithStrings`, or a list of stopwords for each text, when usign with `DetectWithStopWords`. You can pass [options](#options) to the detector to set your [language](#supported-languages), *n*-gram size or a custom stopword list. After executing one of the available detection methods, the detector will write in its interface the final score (float64), the similar *n*-grams (int) and the total *n*-grams (int). Though it seems highily experimental you can see the algorithm in action, in real-time, to [app.mediawatch.io](https//app.mediawatch.io), where we continuosly monitor Greek news outlets. Read the complete documentation at [go-plagiarism](https://pkg.go.dev/github.com/cvcio/go-plagiarism).

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

## Supported Languages
You can find all supported languages in the [stopwords.go](/stopwords.go) file. All supported languages use the ISO639-1 code format as a key (string) and the corresponding stopwods list ([]string) as a value.

| ISO 639-1 	| Language       	| Tested                  	| Tests     |
|-----------	|----------------	|-------------------------	|-------    |
| bg        	| Bulgarian      	| Partially Tested        	| 1         |
| de        	| German         	| Tested (>10K Articles)  	| 1         |
| el        	| Greek          	| Tested (>10M Articles)  	| 5         |
| en        	| English        	| Tested (>1M Articles)   	| 2         |
| fi        	| Finnish        	| Partially Tested        	| 1         |
| fr        	| French         	| Partially Tested        	| 1         |
| hr        	| Croatian       	| Partially Tested        	| 1         |
| hu        	| Hungarian      	| Partially Tested        	| 1         |
| it        	| Italian        	| Tested (>10K Articles)  	| 1         |
| nl        	| Dutch, Flemish 	| Partially Tested        	| 1         |
| no        	| Norwegian      	| Partially Tested        	| 1         |
| pl        	| Polish         	| Partially Tested        	| 1         |
| pt        	| Portuguese     	| Partially Tested        	| 1         |
| ro        	| Romanian       	| Partially Tested        	| 1         |
| ru        	| Russian        	| Tested (>10K Articles)  	| 1         |
| tr        	| Turkish        	| Tested (>100K Articles) 	| 1         |
| uk        	| Ukrainian      	| Partially Tested        	| 1         |

### TODO List

- [ ] Include additional test cases for each language
- [ ] Introduce a `GetSimilar` method to retrieve similar passages
  

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