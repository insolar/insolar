Insolar — Backup manager
================
Utility for merging incremental backups into one db.
If given db does not exist, it creates new one.

Usage
----------
#### Build

    make backupmanager
   
#### Run

    bin/backupmanager
    
#### Options

       ./bin/backupmanager -h
       backupmanager is the command line client for managing backups
       
       Usage:
         backupmanager [command]
       
       Available Commands:
         create         create new empty badger
         help           Help about any command
         merge          merge incremental backup to existing db
         prepare_backup prepare backup for usage
       
       Flags:
         -h, --help   help for backupmanager
       
       Use "backupmanager [command] --help" for more information about a command.

