

to run this server, you need the following file structure:



data
  |-  latest
  |     |- v0.3.4   (empty file with filename = latest version of the platform binary files)
  |
  |-v0.3.1
  |    |- insgocc
  |    |- insgorund
  |    |- insolar
  |    |- insolard
  |    |- insupdater
  |    |- pulsard
  |
  |-v0.3.4
       |- insgocc
       |- insgorund
       |- insolar
       |- insolard
       |- insupdater
       |- pulsard
  
  
  
  HTTP GET REQUEST returns latest version: 
  localhost:2345/latest
  
  RESPONSE:
  {
      "latest": "v0.3.1",
      "major": 0,
      "minor": 3,
      "revision": 1
  }
  
  
  HTTP GET REQUEST returns file
  http://localhost:2345/{VERSION}/{UTILITY_NAME} 
  For example:   http://localhost:2345/v0.3.1/pulsard
  
  Default port: 2345
  For install random port you can start updateserv with parameter "port"