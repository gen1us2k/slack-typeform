# slack-typegorm
[![Go Report Card](http://goreportcard.com/badge/gen1us2k/slack-typeform)](http://goreportcard.com/report/gen1us2k/slack-typeform)
Automate user addition into slack via typeform

Installation
===
Prerequisites
=
1. Golang installed 
```
  git clone https://github.com/gen1us2k/slack-typeform.git
  cd slack-typeform
  go get ./...
  go build
```

Configuring
===

```
  cp config.json.example config.json
```
Configure all params.
If you don't know what fields name of type form open
https://api.typeform.com/v0/form/UID?key=key&completed=true

License
===
Licensed under MIT

Feel free for pull requests, issues or any feedback
