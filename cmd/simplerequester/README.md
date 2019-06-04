SIMPLE requester
===============
   Makes a single signed request to the API URL and shows the result.

Usage
----------
#### Build

    make simplerequester
   
#### Start insolard

    ./scripts/insolard/launchnet.sh -g
   
#### Start simplerequester

    ./bin/simplerequester -k=scripts/insolard/configs/root_member_keys.json

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
                   default use RootMember address)


### Params file structure

        {
    	    "reference":"1tJDL5m9pKyq2mbanYfgwQ5rSQdrpsXbzc1Dk7a53d.1tJDJLGWcX3TCXZMzZodTYWZyJGVdsajgGqyq8Vidw",
        	"method":"CreateMember",
        	"params": {
    	        "name": "John"
        	}
        }
