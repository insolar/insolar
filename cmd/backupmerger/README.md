Insolar — Backup merger
================
Utility for merging incremental backups into one db.
If given db does not exist, it creates new one.

Usage
----------
#### Build

    make backupmerger
   
#### Run

    bin/backupmerger
    
#### Options

    bin/backupmerger -h
    Usage of ./bin/backupmerger:
      -n, --bkp-name string    file name if incremental backup (required)
      -h, --help               show this help
      -t, --target-db string   directory where backup will be roll to (required)
      -w, --workers-num int    number of workers to read backup file (default 1)

