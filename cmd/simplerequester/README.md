SIMPLE requester
===============
   Makes a single signed request to the API URL and shows the result.
   Supported P-256 and P-256K curves.

Usage
----------
#### Build

    make simplerequester
   
#### Start insolard

    ./scripts/insolard/launchnet.sh -g
   
#### Start simplerequester

     ./bin/simplerequester -u=http://localhost:19101/api -f=./bin/request.json -k=scripts/insolard/configs/root_member_keys.json
     ./bin/simplerequester -u=http://localhost:19101/api -m=CreateMember -p={\"name\":\"John\"} -i=1tJDL5m9pKyq2mbanYfgwQ5rSQdrpsXbzc1Dk7a53d.1tJDJLGWcX3TCXZMzZodTYWZyJGVdsajgGqyq8Vidw -k=scripts/insolard/configs/root_member_keys.json
 
### Options

        -k memberkeys
                Path to file with Member keys.

        -u url
                API url for requests (default - http://localhost:19101/api).

        -f paramsFile
                JSON request parameters file
  
           or

        -m method 
                SmartContract method name (for example 'CreateMember')
                
        -p params  
                Request parameters in JSON (for example '{"name":"John"}'

        -i address 
                Smart contract address in the Insolar platform 
                   (for example "1tJDL5m9pKyq2mbanYfgwQ5rSQdrpsXbzc1Dk7a53d.1tJDJLGWcX3TCXZMzZodTYWZyJGVdsajgGqyq8Vidw",
                   default using RootMember address)


### Params file structure

        {
          	    "reference":"1tJDL5m9pKyq2mbanYfgwQ5rSQdrpsXbzc1Dk7a53d.1tJDJLGWcX3TCXZMzZodTYWZyJGVdsajgGqyq8Vidw",
               	"method":"CreateMember",
               	"params": {
           	        "name": "John"
               	}
        }


## Example 1:

    ./bin/simplerequester -u=http://localhost:19101/api -f=./bin/request.json -k=scripts/insolard/configs/root_member_keys.json   

    {
      "jsonrpc": "2.0",
      "id": 1,
      "method": "api.call",
      "params": {
        "callSite": "contract.createMember",
        "callParams": {
        }
      }
    }

    Reference can be null, in this case using RootMemberRef automatically

    Execute result:  &{2.0 map[] map[traceID:f9b592c5-7530-4028-bdfc-a1979ccd285c payload:6qf26syFY27GfdmQ3HHf2PHkxCvXyscmk6X8uvKtuR4.11111111111111111111111111111111]}

## Example 2:
 
    ./bin/simplerequester -u=http://localhost:19101/api -f=./bin/addBurnAddrs.json -k=.artifacts/launchnet/configs/migration_admin_member_keys.json
    {
      "jsonrpc": "2.0",
      "id": 1,
      "method": "api.call",
      "params": {
        "callSite": "wallet.addBurnAddresses",
        "reference": "1tJCnAptaXWX2hopa2hxWN3AJ3f84YZCYkVbqJeJSV.11111111111111111111111111111111",
        "callParams": {
          "burnAddresses": [
            "0xbd8647335d76d5c7b053e769a45239e0c2f6f19a",
            "0x7e8028a397387975fc3ccd84b4253245a7a28da6",
            "0x513337ff17019eaf75117d7d574c6f78ac102b44",
            "0x0ec963ff8965366ce0db256227f7666be386dfd0",
            "0x1a3272eeeac7759adef8bc527804e9a4584fced7",
            "0xf4a471022ff88a8cf0326087a025872d204b8fe2",
            "0xc849e4210083975226385015f9fee7d749f01654",
            "0xd7b9a9b2f665849c4071ad5af77d8c76aa30fb32",
            "0x56ec87150c44dce6d9a73281fe4c5d1f9671805a",
            "0xf4a471022ff88a8cf0326087a025872d204b8fe2",
            "0x32c7a87c60c0417c0a5371348bbdc42e19216574",
            "0x25d48e5f628c73f049cd0299d613e1f9c0a35256",
            "0x77441159d634758249ce852bd5322d2d4138a8f6",
            "0x66312c6ae9f2a90f540f662f3622c2f705852139",
            "0x27e1d264dcce43fbe9e1b48ba1d774d9b8047353",
            "0x27e1d264dcce43fbe9e1b48ba1d774d9b8047353",
            "0xba262de2538f534652a98784a4586955c38434cb",
            "0xb1298e180afd7e3450e7f1c9c46d729b070e2c65",
            "0xd483354dc64d470d930441220e30a849a02fb7b5",
            "0x4388d818f3973ea13b3cd2ef98197a9cfac99ef3"
          ]
        }
      }
    }

    Reference you can get from «network.GetInfo» request, field «MigrationAdminMember» 

    Execute result:  &{2.0 map[] map[payload:<nil> traceID:873e1141-d1e2-4410-b3f0-5761a55a931b]}



