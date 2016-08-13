package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/Zemanta/gracefulshutdown"
	"github.com/Zemanta/gracefulshutdown/shutdownmanagers/awsmanager"
	"github.com/Zemanta/gracefulshutdown/shutdownmanagers/posixsignal"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

// Lifecycled ...
type Lifecycled struct {
	queueName         string
	lifecycleHookName string
	port              uint16
	pingTime          time.Duration
	forwardRetries    int
	cmdWithArgs       []string
	execTimeout       time.Duration
}

const (
	backoff float64 = 10000.0
	version string  = "0.0.1"
)

var execCommand = exec.Command

var (
	queue   = kingpin.Flag("queue", "sqs queue name to use. Required").Required().String()
	hook    = kingpin.Flag("hook", "lifecycle hook name to use. Required").Required().String()
	port    = kingpin.Flag("port", "port to listen on. defaults to 7999").Default("7999").Uint()
	ping    = kingpin.Flag("heartbeat", "heartbeat rate for lifecycle hook. defaults to 10s").Default("10s").Duration()
	timeout = kingpin.Flag("exec-timeout", "exec timeout for script. defaults to 5s").Default("5s").Duration()
	script  = kingpin.Arg("script", "shell script to run on shutdown. Required").Required().Strings()
)

func main() {
	kingpin.Version(version)
	kingpin.Parse()

	ld := Lifecycled{
		queueName:         *queue,
		lifecycleHookName: *hook,
		port:              uint16(*port),
		pingTime:          *ping,
		forwardRetries:    3,
		cmdWithArgs:       *script,
		execTimeout:       *timeout,
	}

	if err := ld.StartMonitor(); err != nil {
		fmt.Printf("FATAL: %+v\n", err)
		os.Exit(1)
	}

	fmt.Println("Done.")
	os.Exit(0)
}

// StartMonitor ...
func (ld *Lifecycled) StartMonitor() error {
	// initialize gracefulshutdown with ping time
	gs := gracefulshutdown.New()

	// set error handler
	gs.SetErrorHandler(gracefulshutdown.ErrorFunc(func(err error) {
		fmt.Println("Error:", err)
	}))

	// add posix shutdown manager and handle SIGTERM
	gs.AddShutdownManager(posixsignal.NewPosixSignalManager(syscall.SIGTERM))

	// add aws shutdown manager and use a 10 second backoff
	gs.AddShutdownManager(awsmanager.NewAwsManager(&awsmanager.AwsManagerConfig{
		SqsQueueName:      ld.queueName,
		LifecycleHookName: ld.lifecycleHookName,
		NumForwardRetries: ld.forwardRetries,
		Port:              ld.port,
		BackOff:           backoff,
		PingTime:          ld.pingTime,
	}))

	// add your tasks that implement ShutdownCallback
	gs.AddShutdownCallback(gracefulshutdown.ShutdownFunc(func(shutdownManager string) error {
		fmt.Println("Shutdown: callback start")
		if shutdownManager == awsmanager.Name {
			fmt.Println("We are being told by the ASG lifecycle hook to shutdown")
			fmt.Println("This instance will be shut down after the shutdown callback finishes")
		}
		if shutdownManager == posixsignal.Name {
			fmt.Println("We are being signaled by a Posix signal to shutdown")
			fmt.Println("This process will exit after the shutdown callback finishes")
		}
		fmt.Println("Executing script")
		if err := RunCommand(ld.cmdWithArgs, ld.execTimeout); err != nil {
			return fmt.Errorf("error running cmd = %v", err)
		}
		fmt.Println("Shutdown: callback finished")
		return nil
	}))

	// start shutdown managers
	if err := gs.Start(); err != nil {
		return err
	}

	// wait until we are sent the SIGQUIT signal to end
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGQUIT)
	<-sigc

	return nil
}

// RunCommand ...
func RunCommand(cmdWithArgs []string, execTimeout time.Duration) error {
	var run *exec.Cmd
	if len(cmdWithArgs) == 0 {
		return fmt.Errorf("no cmd provided to execute")
	} else if len(cmdWithArgs) == 1 {
		run = execCommand(cmdWithArgs[0])
	} else {
		run = execCommand(cmdWithArgs[0], cmdWithArgs[1:]...)
	}

	run.Stdout = os.Stdout
	run.Stderr = os.Stderr

	err := Watch(run, execTimeout)
	if err != nil {
		return fmt.Errorf("error running cmd = %v", err)
	}

	return nil
}

// Watch ...
func Watch(run *exec.Cmd, execTimeout time.Duration) error {
	if err := run.Start(); err != nil {
		return fmt.Errorf("error starting cmd = %v", err)
	}

	// http://stackoverflow.com/a/11886829
	done := make(chan error)
	go func() { done <- run.Wait() }()

	select {
	case err := <-done:
		// The process is done
		if err != nil {
			return fmt.Errorf("process done with error = %v", err)
		}
	case <-time.After(execTimeout):
		// The timout is reached
		if err := run.Process.Kill(); err != nil {
			return fmt.Errorf("failed to kill process = %v", err)
		}
		fmt.Println("process killed as timeout reached")
	}

	return nil
}
