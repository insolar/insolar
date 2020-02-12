# heavy-badger tool

scans badger database and extracts some valuable info from it (addition for badger tool)

## Usage examples

print all known scopes:

    ./bin/heavy-badger scopes

fix (open/close) db (without it badger binary could fail on unclosed db)

    ./bin/heavy-badger --dir ./data fix

show keys only stat for all scopes:

    ./bin/heavy-badger --dir ./data scan scopes-stat --fast


show stat for all scopes:

    ./bin/heavy-badger --dir ./data scan scopes-stat

show stat on records scope

    ./bin/heavy-badger --dir ./data/ scan scope-pulses -s ScopeRecord --limit=1800

show stat on records scope with extra total stat by record type

    ./bin/heavy-badger --dir ./data/ scan scope-pulses -s ScopeRecord --record-types-stat

show stat on records scope with extra total stat by record type

    ./bin/heavy-badger --dir ./data/ scan scope-pulses -s ScopeRecord --per-pulse

show stat on records scope with extra pulse sizes graph report:

    ./bin/heavy-badger --dir ./data/ scan scope-pulses -s ScopeRecord --graph=console

    ./bin/heavy-badger --dir ./data/ scan scope-pulses -s ScopeRecord --per-pulse --graph=web

print all pulses line by line:

    ./bin/heavy-badger --dir ./data/ scan pulses --all

only print short pulses info:

    ./bin/heavy-badger --dir ./data/ scan pulses

### How to find and print values that exceed some threshold

search values in records scope in a first 3000 pulses, which exceeds 600kb threshold:

    ./bin/heavy-badger --dir ./data/ scan scope-pulses -s ScopeRecord --limit=3000 --print-value-gt-size=600000

it prints all found keys in hex and some extra info (pulse number, record id, value type if available) in format:

    key=02016805ec45b214bb07805950ecea9a948ff4f61ef81fdb9cd25692a807de7d2f - FOUND VALUE: size=944 kB type=*record.Amend id=18q4S3nzBAWpZCfH5Z4eNp1LQwWQc5MacU7W3zEBBxKW.record pulse=33646597

now we can show binary dump of value by key:

    ./bin/heavy-badger --dir ./data/ dump bin 020168074028f5debd3359226e348e9f09b7f2c26aa762095a0bc4a1e681613961 | hexdump -C | vim -

or try to deserialize value to record by key:

    ./bin/heavy-badger --dir ./data/ dump record 020168074028f5debd3359226e348e9f09b7f2c26aa762095a0bc4a1e681613961


you can get:

    ------------------------------------------------------------------------------
    Material Record:
                      ID: 16VQqDEE1tw75K9ib85MRHsqCbsfUApnFHZGt9x3H6XH.record
                   JetID: [JET 5 11001]
                ObjectID: 11tJDm7rr8RVx1PA9Sakpwf8YmyS5fpKHMEVuc9E5P2.record

    Virtual Record:
                    Type: *record.Amend
                 request: 16VQqCzpd5Smbdi89b8N7CTK6a8sNWFk4ZNFLWGBSCbt.record
                  memory: 935 kB
                   image: 0111A5x8N1VJTm7BKYgzSe6TWHcFi98QZgw3AnkYiKML
             isPrototype: false
               prevState: 16VQpHpBmwLis59DYUEoFEWbHtKJ8tDNWDKtRhxZWRDx.record
    ------------------------------------------------------------------------------
