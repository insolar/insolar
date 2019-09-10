# inslogrotator

Tool for controlling logs reopening.
Used for log rotation which untied from system environment (syslog/journald/etc).

## Use case

1. pass log streams through worker instances to files
2. rotate files by removing/moving files, then send `SIGUSR2` signal to all `inslogrotator` processes.

## How it works

- passes all incoming stdin to file, provided in commandline.
- if got SIGUSR2 signal, just reopen file and continues to pass incoming stream
