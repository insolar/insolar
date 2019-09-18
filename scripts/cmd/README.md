## track
Can be used to track a sequence of events. Object tracking example:

    `ag 1tJDBK2rT2hJvuGGCmzCkV8QPDKA4DQ278noZ86wX1 .artifacts/launchnet/logs | bin/track`
    
As a result, time-sorted logs will be printed. Logs will be grouped by file for read simplicity. 

## logstat

Parses logs and calculates aggregates by traceid. Use to analyse request messages, performance and stats.

### Usage
Set log format to `json` and log level to at least `info`.

Usage example:
```
go run ./scripts/cmd/logstat/logstat.go -d .artifacts/launchnet/logs/
Analysed 36848 logs in 216 traces.
[Average count per trace]
    TypeSetResult: 16.5
    TypeCallMethod: 10.3
    TypeSetIncomingRequest: 10.3
    TypeGetObject: 9.7
    TypeGetRequestInfo: 9.3
    TypeSetOutgoingRequest: 9.3
    TypeReturnResults: 8.8
    TypePassState: 4.8
    TypeUpdate: 2.5
    TypeSagaCallAcceptNotification: 1.5
    TypeActivate: 0.6
    TypeHasPendings: 0.2
    TypeGetIndex: 0.1
    TypeReplication: 0.1
    TypeGetFilament: 0.1
    TypePendingFinished: 0.0
    TypePass: 0.0
    TypeUpdateJet: 0.0
    TypeGetJet: 0.0

[Average reply times per trace, ms]
    TypeCallMethod: 282
    TypeSetIncomingRequest: 191
    TypeSetResult: 180
    TypeGetObject: 126
    TypeSetOutgoingRequest: 93
    TypeGetRequestInfo: 82
    TypeReturnResults: 69
    TypeUpdate: 24
    TypePendingFinished: 1
    TypeActivate: 1
    TypeSagaCallAcceptNotification: 0
    TypeUpdateJet: 0
    TypePass: 0
    TypeGetJet: 0
    TypeHasPendings: 0
    TypePassState: 0
    TypeGetIndex: 0
    TypeGetFilament: 0
    TypeReplication: 0

[Call return percentiles]
< 1s | 205
   c9f7053e-a31e-4484-b395-b3029bd4daac | 535.263ms
   d9c07bac-458d-43b4-9d19-f68ac5804d26 | 518.862ms
   66881671-b9e0-4cb2-b512-92b59b870b8c | 449.375ms
< 10s | 0
< 20s | 0
< 40s | 0
< 1m0s | 0
> 1m0s | 0

[Total time percentiles]
< 1s | 58
   4840efa3-1739-41b4-8449-c2e9fce5f30c | 64.906ms
   8bd84109-e13e-43ca-af9b-d612352dc5da | 68.226ms
   c095f564-979c-4274-8c5f-4b9a5727fa36 | 95.508ms
< 10s | 158
   c9f7053e-a31e-4484-b395-b3029bd4daac | 4.346744s
   d9c07bac-458d-43b4-9d19-f68ac5804d26 | 4.159466s
   66881671-b9e0-4cb2-b512-92b59b870b8c | 3.338537s
< 20s | 0
< 40s | 0
< 1m0s | 0
> 1m0s | 0
```

Use `logstat -h` for more info.
