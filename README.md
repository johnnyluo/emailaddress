# emailaddress : validate email address based on RFC-5322

## Overview [![GoDoc](https://godoc.org/github.com/johnnyluo/emailaddress?status.svg)](https://godoc.org/github.com/johnnyluo/emailaddress) [![Build Status](https://travis-ci.com/johnnyluo/emailaddress.svg?branch=master)](https://travis-ci.com/johnnyluo/emailaddress)

Email address is defined in [RFC-5322](https://tools.ietf.org/html/rfc5322#section-3.4.1) , a more easy to read version can be find on [wikipedia](https://en.wikipedia.org/wiki/Email_address). This package provide the function to validate an email address against the standard.

## Install

```bash
go get github.com/johnnyluo/emailaddress
```

## Example

### How to validate an email address

```go
b, err := emailaddress.Validate(input)
if nil != err {
    panic(err)
}
if b {
    fmt.Println("%s is a valid email address",input)
} else {
    fmt.Println("%s is an invalid email address",input)
}

```

### Check whether two mailbox is equal

johnny+1@test.net and johnny+2@test.net are both legitimate emamil address, but they might all end up to johnny@test.net mailbox.  This library provide a method to check whether two email address are semantically equal

```go
if emailaddress.Equals("johnny+1@test.net","johnny+2@test.net") {
    fmt.Println("They are equal")
} else {
    fmt.Println("They are different")
}

```

## License

Apache 2.0.
