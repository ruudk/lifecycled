##### Usage

Ensure everything is working okay.

```
mkdir ~/.go
export GOPATH=~/.go
export GOBIN=$GOPATH/bin
export PATH=$GOBIN:$PATH
go get golang.org/x/tour/gotour
gotour
```

Install lifecycled.

```
go get -u github.com/Chansey-Org/lifecycled
```

Help.

```
lifecycled --help
```

`Output`

```
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

Run.

```
lifecycled --queue MyQueue --hook my-lifecycle-hook --exec-timeout 30s -- sleep 10
```

If the SIGTERM signal was sent, the posix signal manager will fire off the shell command.

```
Shutdown: callback start
shutdownManager: PosixSignalManager
Doing some work for posixsignal
```

Otherwise if a lifecycle hook was called, then the aws manager will fire off the shell command.

```
Shutdown: callback start
shutdownManager: AwsManager
Doing some work for awsmanager
Shutdown: callback finished
```

When work has been completed, the shutdown callback will complete and the program will exit

```
Shutdown: callback start
...Work
Shutdown: callback finished
```

Send SIGQUIT or SIGINT to exit prematurely.

```
pkill -SIGQUIT lifecycled
```