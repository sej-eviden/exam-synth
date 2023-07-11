# Exam Topics Scraper - reimagined

## Purpose

test tests

> Scrape Exam Topics exams, but better.

Pseudo-scrape html page to get json data of exams & questions.
Ease of use and portability are a must. The main focus is to make life easier,
any step away from this, is a step bakwards.

### Main points to improve:

- Portability
- Size
- Ease of use
- Scope a.k.a amount of manual work for the end product

## Roadmap

- [x] Port all *features* from original Python code
    - [x] Create necessary folders (results + specific)
    - [x] Scrape all questions from html
    - [x] Produce json
- [x] Create `img` folder too
- [ ] Paste imgs from original to the destination folder
[Possible solution](https://stackoverflow.com/questions/50740902/move-a-file-to-a-different-drive-with-go) 
- [ ] Create `master` file with general information about the exam (img answers,
empty answers, amount of questions, amount of topics, alternative names (code), etc.)
- [ ] Give some useful output in terminal and describe next steps (terminal and, 
optionally, a small file)

## Use

1. Download and add to path
2. Create `main` folder (name does't matter) with all exams inside `exams` folder
(different name can be specified)
3. Run from the terminal inside the `main` folder

```sh
e-synth ./path/to/exams -dest results
```

Where `./path/to/exams` is the folder where the exams are located, and `results`
is the folder where the output will be placed.

