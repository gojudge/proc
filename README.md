# proc [![Build Status](https://travis-ci.org/gojudge/proc.svg?branch=master)](https://travis-ci.org/gojudge/proc) [![Coverage Status](https://coveralls.io/repos/github/gojudge/proc/badge.svg?branch=master)](https://coveralls.io/github/gojudge/proc?branch=master)

- get process info
- get process chain
- kill process chain

# usage

```go
proc := GetProc(1)
proc.KillProcChainReverse()
```

# license

MIT License
