A simple library that uses https://github.com/Zemanta/gracefulshutdown to provide functionality similar to https://github.com/lox/lifecycled.

More info on how to use auto scaling lifecycle hooks can be found here: http://docs.aws.amazon.com/autoscaling/latest/userguide/lifecycle-hooks.html.

Travis.

[![Build Status](https://travis-ci.org/jason-riddle/lifecycled.svg?branch=master)](https://travis-ci.org/jason-riddle/lifecycled)

Coverage.

[![codecov](https://codecov.io/gh/jason-riddle/lifecycled/branch/master/graph/badge.svg)](https://codecov.io/gh/jason-riddle/lifecycled)

Install.

```
$ go get github.com/Chansey-Org/lifecycled
```

Help.

```
$ lifecycled --help
usage: lifecycled --queue=QUEUE --hook=HOOK [flags>] script>...

Flags:
  --help             Show context-sensitive help (also try --help-long and --help-man).
  --queue=QUEUE      sqs queue name to use. Required
  --hook=HOOK        lifecycle hook name to use. Required
  --port=7999        port to listen on. defaults to 7999
  --heartbeat=10s    heartbeat rate for lifecycle hook. defaults to 10s
  --exec-timeout=5s  exec timeout for script. defaults to 5s
  --version          Show application version.

Args:
  script>  shell script to run on shutdown. Required

exit status 1
```

Usage.

```
$ lifecycled --queue MyQueue --hook my-lifecycle-hook -- sleep 10
```

Ensure you have created your lifecycle hook and SQS queue before running. See the [GUIDE](GUIDE.md) for more info.
