# go-aoj: local judge command for AOJ

The go-aoj command judges your solution on your local machine with downloaded
test cases.
Currently it supports the Introduction set of problems (ITP1, ALDS1, etc.)

##  Install

```
$ go get github.com/ktateish/go-aoj/aoj
```

## Usage

### Check all cases

```
$ aoj check ALDS1_1_A mybinary
```

### Check the specified case

```
$ aoj check -i 0 ALDS1_1_A mybinary
```

### Get the number of cases

```
$ aoj testcase -l ALDS1_1_A
5
```

### Get the number of cases

```
$ aoj testcase -l ALDS1_1_A
5
```

## Author

Katsuyuki Tateishi <kt@wheel.jp>
