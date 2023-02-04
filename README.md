# clsr

`clsr` stands for **C**ommand **L**ine **S**paced **R**epetition.
It is a learning tool that uses [spaced repetition](https://en.wikipedia.org/wiki/Spaced_repetition)
to help you learn and retain information efficiently.
It is similar to other spaced repetition software, but takes a
minimalist, text-based, and version-controllable approach.

`clsr` was inspired by [`ledger`](https://github.com/ledger/ledger), and
more generally by the [plain text accounting](https://plaintextaccounting.org/)
ecosystem.


## How does `clsr` work?

First, some terms:

A **card** is a virtual flash card. It has a question or prompt for something
you want to learn, and the answer to that question.  
A **deck** is a group of related cards. For example, you might make a deck
for French words you want to learn, or for parts of syntax of a programming
language you want to learn.

In `clsr`, decks take the form of JSON files. The idea is that you keep
your deck files in a single directory. Then when you run `clsr` commands
from inside that directory, `clsr` can work with those files.
Having these files in one directory also lends itself to the use of
version control.


## Should I use `clsr`?

`clsr` will work well for you if:

- You want to easily understand how your data is stored
- You want to store your card data in version control
- You want to write scripts that use your card data
- You work over SSH, or want the option to do so

You should not use `clsr` if:

- You are not comfortable with the command line
- You need to include anything other than text (i.e. pictures, sounds)
  in your cards
- You need fancier features such as cloze and reversed cards


## Installation

Have go 1.17 or later installed. Then:

```
go install github.com/adamkpickering/clsr@latest
```


## Credits

Thanks to SUSE for holding [Hack Week](https://hackweek.opensuse.org/) 22,
which helped to get `clsr` to the point where it is usable!
