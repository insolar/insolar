Insolar — Backup manager
================
Utility for merging incremental backups into one db.
If given db does not exist, it creates new one.

## Build

```
make backupmanager
```

## Typical usage

Create an empty backup:

```
./bin/backupmanager create -d ./heavy_backup
```

Configure Heavy to execute a given script when incremental backup is ready. Example:

```
ledger:
  ...skipped...
  backup:
    enabled: true
    tmpdirectory: "/tmp/heavy/tmp"
    targetdirectory: "/tmp/heavy/target"
    metainfofile: meta.json
    confirmfile: BACKUPED
    backupfile: incr.bkp
    dirnametemplate: pulse-%d
    backupwaitperiod: 60
    postprocessbackupcmd:
    - bash
    - -c
    - ./bin/backupmanager merge -n $INSOLAR_CURRENT_BACKUP_DIR/incr.bkp -t ./heavy_backup && touch $INSOLAR_CURRENT_BACKUP_DIR/BACKUPED
```

The script should run `backupmanager merge` for a given incremental backup (the path is provided in `$INSOLAR_CURRENT_BACKUP_DIR` environment variable) and the backup directory, created at the first step. When backup is merged the script should create a BACKUPED file in the `$INSOLAR_CURRENT_BACKUP_DIR` directory. When Heavy sees this file it knows that the backup was sucessfuly done and the pulse can be finilized. If necessary the `postprocessbackupcmd` can use `rsync`, `scp`, etc.

To restore the database from a backup run:

```
./bin/backupmanager prepare_backup -d ./heavy_backup/ -l last_backup_info.json
```

This command marks the last pulse in the backup as finalized. We have to do it because during the backup the last pulse is not finalized yet.

After executing the command replace `data` directory on Heavy with `heavy_backup` and start the network.

## Using a backup daemon

`backupmanager merge` is executed much faster when a backup daemon is used.

Start a backup daemon:

```
./bin/backupmanager daemon -t ./heavy_backup
```

Instead of `backupmanager merge` use:

```
./bin/backupmanager daemon-merge -a http://localhost:8099 -g -n $INSOLAR_CURRENT_BACKUP_DIR/incr.bkp
```

See `--help` output for more details.

To restore from a backup kill the backup daemon and use `prepare_backup` as usual.
