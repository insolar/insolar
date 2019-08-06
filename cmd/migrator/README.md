Migrator
===============
   Emulate migration daemon utility for create deposits.
   Supported P-256 and P-256K curves.

Usage
----------
#### Build

    make simplerequester
   
#### Start insolard

    ./scripts/insolard/launchnet.sh -g
   
#### Start simplerequester

     ./bin/migrator -u=http://localhost:19101/api -f=./bin/migration.json -k=./.artifacts/launchnet/configs/ -i=0
   
### Options

        -k memberkeys
                Path to the dirrectory with Migration Member keys.

        -u url
                API url for requests (default - http://localhost:19101/api).

        -i index
                index of migration_daemon_{index}_member_keys.json file
     
        -i address 
                Smart contract address in the Insolar platform 
                   (for example "1tJDL5m9pKyq2mbanYfgwQ5rSQdrpsXbzc1Dk7a53d.1tJDJLGWcX3TCXZMzZodTYWZyJGVdsajgGqyq8Vidw",
                   default using RootMember address)


### Params file structure

     {
       "jsonrpc": "2.0",
       "id": 1,
       "method": "api.call",
       "params": {
         "callSite": "deposit.migration",
         "callParams": {
           "amount": "1000000000",
           "ethTxHash": "394578234932493486739856jfgd48756348563495846djf",  - this is unique value for 1 deposit
           "migrationAddress": "123"   - change here to valid migration address
         }
       }
     }



