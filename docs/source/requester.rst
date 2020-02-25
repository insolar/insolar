.. _building_requester:

===================
Build API Requester
===================

This tutorial walks you through the process of creating a requester to Insolar API. You will learn how to **form and sign requests** that create a member capable of transferring funds to other members.

.. _what_you_will_build:

What You Will Build
-------------------

You will build a program (requester) that creates a member in the Insolar network and, as a member, transfers funds from its account.

The requester forms and sends the following requests to the Insolar’s JSON RPC 2.0 API endpoint and receives the corresponding responses:

#. Gets a **seed** that enables you to call a contract method.

#. Forms, signs, and sends a **member creation** request with the seed and receives your member’s reference in response.

#. Forms, signs, and sends a **transfer** request with another seed that sends an amount of Insolar coins (XNS):

   * from your member’s account given your reference received in the previous response;
   * to another’s account given a reference to the existing recipient member.

.. _what_you_will_need:

What You Will Need
------------------

* About an hour.
* Your favorite IDE for Golang.
* `Insolar API specification <https://apidocs.insolar.io/platform/latest>`_ as a reference.
* :ref:`Local Insolar deployment <deploying_devnet>` or testing environment provided by Insolar.

.. _how_to_complete:

How to Complete This Tutorial
-----------------------------

* To start from scratch, go through the :ref:`step-by-instructions <build_requester>` listed below and pay attention to comments in code examples.
* To skip the basics, read (and copy-paste) the working :ref:`requester code <requester_example>` provided at the end.

.. _build_requester:

Building the Requester
----------------------

To build the requester, go through the following steps:

#. **Prepare**: install the necessary tools, import the necessary packages, and initialize an HTTP client.

#. **Declare** request **structures** in accordance with the Insolar’s API specification.

#. **Create a seed getter** function. The getter will be reused to put a new seed into every signed request.

#. **Create a sender** function that signs and sends requests given a private key.

#. **Generate a key pair**, export the private key into a file, and store it in some secure place.

#. **Form a member creation request** and call the sender function to send it.

#. **Form a transfer request** and, again, call the sender function.

#. **Test** the requester against the testing environment.

All the above steps are detailed in sections below.

.. _prepare:

Step 1: Prepare
~~~~~~~~~~~~~~~

To build the requester, install, import, and set up the following:

#. Set up your development environment if you do not have one. Install the `Go programming tools <https://golang.org/doc/install>`_.
#. With the Golang programming tools, you do not need to “reinvent the wheel”: create a ``Main.go`` file and, inside, import the packages your requester will use. For example:

   .. code-block:: Go
      :linenos:

      package main

      import (
        // You will need:
        // - Some basic Golang functionality.
        "os"
        "bytes"
        "io/ioutil"
        "fmt"
        "log"
        "strconv"
        // - HTTP client and a cookiejar.
        "net/http"
        "net/http/cookiejar"
        "golang.org/x/net/publicsuffix"
        // - Big numbers to store signatures.
        "math/big"
        // - Basic cryptography.
        "crypto/x509"
        "crypto/elliptic"
        "crypto/ecdsa"
        "crypto/rand"
        "crypto/sha256"
        // - Basic encoding capabilities.
        "encoding/pem"
        "encoding/json"
        "encoding/base64"
        "encoding/asn1"
      )

#. To prepare the requester, do the following:

   #. Insolar supports ECDSA-signed requests. Since an ECDSA signature in Golang consists of two big integers, declare a single structure to contain it.

      .. _set_url:

   #. Set the API endpoint URL for the testing environment, either the public one provided by Insolar or :ref:`locally deployed <deploying_devnet>`.
   #. Create and initialize an HTTP client for connection re-use and store a ``cookiejar`` inside.
   #. Create a variable for the JSON RPC 2.0 request identifier. The identifier is to be incremented for every request and each corresponding response will contain it.

   .. _cookie:

   For example:

   .. code-block:: Go
      :linenos:

      // Declare a structure to contain the ECDSA signature:
      type ecdsaSignature struct {
        R, S *big.Int
      }

      // Set the endpoint URL for the testing environment:
      const (
        TestNetURL = "https://wallet-api.test.insolar.io/api/rpc"
      )

      // Create and initialize an HTTP client for connection re-use and put a cookiejar into it:
      var client *http.Client
      var jar cookiejar.Jar
      func init() {
        // All users of cookiejar should import "golang.org/x/net/publicsuffix"
        jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
        if err != nil {
          log.Fatal(err)
        }
        client = &http.Client{
          Jar: jar,
        }
      }

      // Create a variable for the JSON RPC 2.0 request identifier:
      var id int = 1
      // The identifier is to be incremented for every request and each corresponding response will contain it.

With that, everything your requester needs is set up.

.. _declare_structs_or_classes:

Step 2: Declare Request Structures
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Next, declare request structures in accordance with the Insolar’s API specification.

To transfer funds, you need structures or classes for:

#. Information request: ``node.getSeed``.
#. Contract requests: ``member.create`` and ``member.transfer``.

Both information and contract requests have the same base structure in accordance with the `JSON RPC 2.0 specification <https://www.jsonrpc.org/specification>`_.
Therefore, define the base structure once and expand it for all requests with their specific fields.

For example:

.. code-block:: Go
   :linenos:

   // Continue in the Main.go file...

   // Declare a nested structure to form requests to Insolar API in accordance with the specification.
   // The Platform uses the basic JSON RPC 2.0 request structure:
   type requestBody struct {
     JSONRPC        string         `json:"jsonrpc"`
     ID             int            `json:"id"`
     Method         string         `json:"method"`
   }

   type requestBodyWithParams struct {
     JSONRPC        string         `json:"jsonrpc"`
     ID             int            `json:"id"`
     Method         string         `json:"method"`
     // Params is a structure that depends on a particular method:
     Params         interface{}    `json:"params"`
   }

   // The Platform defines params of the signed request as follows:
   type params struct {
     Seed            string       `json:"seed"`
     CallSite        string       `json:"callSite"`
     // CallParams is a structure that depends on a particular method.
     CallParams      interface{}  `json:"callParams"`
     PublicKey       string       `json:"publicKey"`
   }

   type paramsWithReference struct {
     params
     Reference       string  `json:"reference"`
   }

   // The member.create request has no parameters, so it's an empty structure:
   type memberCreateCallParams struct {}

   // The transfer request sends an amount of funds to member identified by a reference:
   type transferCallParams struct {
     Amount            string    `json:"amount"`
     ToMemberReference string    `json:"toMemberReference"`
   }

Now that the requester knows which information and contract requests it is supposed to send, create the following functions:

#. Seed getter for the information request.
#. Sender for contract requests.

.. _create_seed_getter:

Step 3: Create a Seed Getter
~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Each signed request to Insolar API has to contain a seed in its body. Seed is a unique piece of information generated by a node that:

* has a short lifespan;
* expires upon first use;
* protects from duplicate requests.

.. tip:: Due to these qualities, a new seed is required to form each signed contract request.

.. caution:: Since the seed is generated by a node, each subsequent contract request containing the seed must be sent to the node in question. Otherwise, a node will reject the seed generated by a different one. To ensure that the contract request is routed to the correct node, make sure to retrieve all the cookies from the node and store them in the HTTP client intended for re-use as described in the :ref:`preparation step <cookie>`.

To be able to send signed requests, create a seed getter function to re-use upon forming each such request.

The seed getter:

#. Forms a ``node.getSeed`` request body in JSON format.
#. Creates an *unsigned* HTTP request with the body and a Content-Type (``application/json``) HTTP header.
#. Sends the request and receives a response.
#. Retrieves the acquired seed from the response and returns it.

For example:

.. code-block:: Go
   :linenos:

   // Continue in the Main.go file...

   // Create a function to get a new seed for each signed request:
   func getNewSeed() (string) {
     // Form a request body for getSeed:
     getSeedReq := requestBody{
       JSONRPC: "2.0",
       Method:  "node.getSeed",
       ID:      id,
     }
     // Increment the id for future requests:
     id++

     // Marshal the payload into JSON:
     jsonSeedReq, err := json.Marshal(getSeedReq)
     if err != nil {
       log.Fatalln(err)
     }

     // Create a new HTTP request and send it:
     seedReq, err := http.NewRequest("POST", TestNetURL, bytes.NewBuffer(jsonSeedReq))
     if err != nil {
       log.Fatalln(err)
     }
     seedReq.Header.Set("Content-Type", "application/json")

     // Perform the request:
     seedResponse, err := client.Do(seedReq)
     if err != nil {
       log.Fatalln(err)
     }
     defer seedReq.Body.Close()

     // Receive the response body:
     seedRespBody, err := ioutil.ReadAll(seedResponse.Body)
     if err != nil {
       log.Fatalln(err)
     }

     // Unmarshal the response:
     var newSeed map[string]interface{}
     err = json.Unmarshal(seedRespBody, &newSeed)
     if err != nil {
       log.Fatalln(err)
     }

     // (Optional) Print the request and its response:
     print := "POST to " + TestNetURL +
       "\nPayload: " + string(jsonSeedReq) +
       "\nResponse status code: " +  strconv.Itoa(seedResponse.StatusCode) +
       "\nResponse: " + string(seedRespBody) + "\n"
     fmt.Println(print)

     // Retrieve and return the current seed:
     return newSeed["result"].(map[string]interface{})["seed"].(string)
   }

Now, every ``getNewSeed()`` call will return a living seed that can be put into the contract request body.

The next step is to create a sender function that signs and sends contract requests.

.. _create_sender:

Step 4: Create a Sender Function
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

The sender function:

#. Takes some request body (payload) and the ECDSA private key.
#. Forms an HTTP request with the payload and relevant HTTP headers:

   #. *Content-Type* — ``application/json``.
   #. *Digest* that contains (1) a SHA-256 hash of the payload's bytes (2) represented as a Base64 string.
   #. *Signature* that contains (1) the ECDSA signature of the hash's bytes (2) in the ASN.1 DER format (3) represented as a Base64 string.

#. Sends the request.
#. Returns the response JSON object.

For example:

.. tip:: In Golang, the ECDSA signature consists of two big integers. To convert the signature into the ASN.1 DER format, put it into the ``ecdsaSignature`` structure.

.. code-block:: Go
   :linenos:

   // Continue in the Main.go file...

   // Create a function to send signed requests:
   func sendSignedRequest(payload requestBodyWithParams, privateKey *ecdsa.PrivateKey) map[string]interface{} {
     // Marshal the payload into JSON:
     jsonPayload, err := json.Marshal(payload)
     if err != nil {
       log.Fatalln(err)
     }

     // Take a SHA-256 hash of the payload's bytes:
     hash := sha256.Sum256(jsonPayload)

     // Sign the hash with the private key:
     r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
     if err != nil {
       log.Fatalln(err)
     }

     // Convert the signature into ASN.1 DER format:
     sig := ecdsaSignature{
       R: r,
       S: s,
     }
     signature, err := asn1.Marshal(sig)
     if err != nil {
       log.Fatalln(err)
     }

     // Convert both hash and signature into a Base64 string:
     hash64 := base64.StdEncoding.EncodeToString(hash[:])
     signature64 := base64.StdEncoding.EncodeToString(signature)

     // Create a new request and set its headers:
     request, err := http.NewRequest("POST", TestNetURL, bytes.NewBuffer(jsonPayload))
     if err != nil {
       log.Fatalln(err)
     }
     request.Header.Set("Content-Type", "application/json")

     // Put the hash string into the HTTP Digest header:
     request.Header.Set("Digest", "SHA-256="+hash64)

     // Put the signature string into the HTTP Signature header:
     request.Header.Set("Signature", "keyId=\"public-key\", algorithm=\"ecdsa\", headers=\"digest\", signature="+signature64)

     // Send the signed request:
     response, err := client.Do(request)
     if err != nil {
       log.Fatalln(err)
     }
     defer response.Body.Close()

     // Receive the response body:
     responseBody, err := ioutil.ReadAll(response.Body)
     if err != nil {
       log.Fatalln(err)
     }

     // Unmarshal it into a JSON object:
     var JSONObject map[string]interface{}
     err = json.Unmarshal(responseBody, &JSONObject)
     if err != nil {
       log.Fatalln(err)
     }

     // (Optional) Print the request and its response:
     print := "POST to " + TestNetURL +
       "\nPayload: " + string(jsonPayload) +
       "\nResponse status code: " + strconv.Itoa(response.StatusCode) +
       "\nResponse: " + string(responseBody) + "\n"
     fmt.Println(print)

     // Return the response:
     return JSONObject
   }

Now, every ``sendSignedRequest(payload, privateKey)`` call will return the result of a contract method.

With the seed getter and sender functions, you can get the seed and send signed contract requests. The next step is to generate a key pair.

.. _generate_key_pair:

Step 5: Generate a Key Pair
~~~~~~~~~~~~~~~~~~~~~~~~~~~

The body of each request that calls a contract method must be hashed by a ``SHA256`` algorithm. Each hash must be signed by a private key generated by a ``P256`` elliptic curve.

To be able to sign requests, do the following:

#. Generate a key pair using the said curve and convert it into PEM format.

   .. warning:: You will not be able to access your member object without the private key and, as such, transfer funds.

#. Export the private key into a file.
#. Save the file to some secure place.

For example:

.. tip:: In Golang, to encode the key into the PEM format, first, convert it into ASN.1 DER using the ``x509`` library.

.. code-block:: Go
   :linenos:

   // Continue in the Main.go file...

   // Create the main function to form and send signed requests:
   func main() {
     // Generate a key pair:
     privateKey := new(ecdsa.PrivateKey)
     privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
     var publicKey ecdsa.PublicKey
     publicKey = privateKey.PublicKey

     // Convert both private and public keys into PEM format:
     x509PublicKey, err := x509.MarshalPKIXPublicKey(&publicKey)
     if err != nil {
       log.Fatalln(err)
     }
     pemPublicKey := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: x509PublicKey})

     x509PrivateKey, err := x509.MarshalECPrivateKey(privateKey)
     if err != nil {
       log.Fatalln(err)
     }
     pemPrivateKey := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: x509PrivateKey})

     // The private key is required to sign requests.
     // Make sure to put into a file to save it in some secure place later:
     file, err := os.Create("private.pem")
     if err != nil {
       fmt.Println(err)
       return
     }
     file.WriteString(string(pemPrivateKey))
     file.Close()

      // The main function is to be continued...
    }

Now that the key pair is generated and saved, you can form contract requests.

.. _form_member_create:

Step 6: Form and Send a Member Creation Request
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

The member creation request is a signed request to a contract method that does the following in the blockchain:

* Creates a new member and corresponding account objects.
* Returns the new member reference — address in the Insolar network.
* Binds a given public key to the member. Insolar uses this public key to identify a member and check the signature generated by the paired private key.

To create a member:

#. Call the ``getNewSeed()`` function and store the new seed into a variable.
#. Form the ``member.create`` request payload with the seed and the public key generated in the :ref:`previous step <generate_key_pair>`.
#. Call the ``sendSignedRequest()`` function and pass it the payload and the private key.
#. Put the returned member reference into a variable. The subsequent transfer request requires it.

For example:

.. code-block:: Go
   :linenos:

   // Continue in the main() function...

   // Get a seed to form the request:
   seed := getNewSeed()
   // Form a request body for member.create:
   createMemberReq := requestBodyWithParams{
     JSONRPC: "2.0",
     Method:  "contract.call",
     ID:      id,
     Params:params {
       Seed: seed,
       CallSite: "member.create",
       CallParams:memberCreateCallParams {},
       PublicKey: string(pemPublicKey)},
   }
   // Increment the JSON RPC 2.0 request identifier for future requests:
   id++

   // Send the signed member.create request:
   newMember := sendSignedRequest(createMemberReq, privateKey)

   // Put the reference to your new member into a variable to send transfer requests:
   memberReference := newMember["result"].(map[string]interface{})["callResult"].(map[string]interface{})["reference"].(string)
   fmt.Println("Member reference is " + memberReference)

   // The main function is to be continued...

Now that you have your member reference, you can transfer funds to other members.

.. _form_transfer:

Step 7: Form and Send a Transfer Request
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

The transfer request is a signed request to a contract method that transfers some amount of funds to another member.

To transfer funds:

#. Acquire the recipient reference — the reference to an existing member to whom you want to transfer the funds.
#. Call the ``getNewSeed()`` function and store the new seed into a variable.
#. Form a ``member.transfer`` request payload with:

   * a new seed,
   * an amount of funds to transfer,
   * the recipient reference,
   * your reference (for identification),
   * and your public key (to check the signature).

#. Call the ``sendSignedRequest()`` function and pass it the payload and the paired private key.

The transfer request will return the factual fee value in its response.

For example:

.. attention:: In the highlighted line, replace the ``<recipient_member_reference>`` placeholder value with the reference to the existing recipient member.

.. code-block:: Go
   :linenos:
   :emphasize-lines: 15

   // Continue in the main() function...

   // Get a new seed to form a transfer request:
   seed = getNewSeed()
   // Form a request body for transfer:
   transferReq := requestBodyWithParams{
     JSONRPC: "2.0",
     Method:  "contract.call",
     ID:      id,
     Params:paramsWithReference{ params:params{
       Seed: seed,
       CallSite: "member.transfer",
       CallParams:transferCallParams {
         Amount: "100",
         ToMemberReference: "<recipient_member_reference>",
         },
       PublicKey: string(pemPublicKey),
       },
       Reference: string(memberReference),
     },
   }
   // Increment the id for future requests:
   id++

   // Send the signed transfer request:
   newTransfer := sendSignedRequest(transferReq, privateKey)
   fee := newTransfer["result"].(map[string]interface{})["callResult"].(map[string]interface{})["fee"].(string)

   // (Optional) Print out the fee.
   fmt.Println("Fee is " + fee)

   // Remember to close the main function.
   }

With that, the requester, as a member, can send funds to other members of the Insolar network.

.. _test_requester:

Step 8: Test the Requester
~~~~~~~~~~~~~~~~~~~~~~~~~~

To test the requester, do the following:

#. Make sure the :ref:`endpoint URL <set_url>` is set to that of the testing environment.
#. Run the requester:

   .. code-block:: console

      $ go run Main.go

.. _Summary:

Summary
-------

Congratulations! You have just developed a requester capable of forming signed requests to interact with the Insolar API.

Build upon it:

#. Create structures for other requests in accordance with the Insolar API specification.
#. Export the getter and sender functions to use them in other packages.

.. _requester_example:

Complete Requester Code Example
-------------------------------

Below are the complete requester code example in Golang.

.. attention:: To be able to send transfer requests, in the highlighted line, replace the ``<recipient_member_reference>`` placeholder value with the reference to the existing recipient member.

.. code-block:: Go
   :linenos:
   :emphasize-lines: 294

   package main

   import (
     // You will need:
     // - Some basic Golang functionality.
     "os"
     "bytes"
     "io/ioutil"
     "fmt"
     "log"
     "strconv"
     // - HTTP client and a cookiejar.
     "net/http"
     "net/http/cookiejar"
     "golang.org/x/net/publicsuffix"
     // - Big numbers to store signatures.
     "math/big"
     // - Basic cryptography.
     "crypto/x509"
     "crypto/elliptic"
     "crypto/ecdsa"
     "crypto/rand"
     "crypto/sha256"
     // - Basic encoding capabilities.
     "encoding/pem"
     "encoding/json"
     "encoding/base64"
     "encoding/asn1"
   )

   // Declare a structure to contain the ECDSA signature:
   type ecdsaSignature struct {
     R, S *big.Int
   }

   // Set the endpoint URL for the testing environment:
   const (
     TestNetURL = "https://wallet-api.test.insolar.io/api/rpc"
   )

   // Create and initialize an HTTP client for connection re-use and put a cookiejar into it:
   var client *http.Client
   var jar cookiejar.Jar
   func init() {
     // All users of cookiejar should import "golang.org/x/net/publicsuffix"
     jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
     if err != nil {
       log.Fatal(err)
     }
     client = &http.Client{
       Jar: jar,
     }
   }

   // Create a variable for the JSON RPC 2.0 request identifier:
   var id int = 1
   // The identifier is to be incremented for every request and each corresponding response will contain it.

   // Declare a nested structure to form requests to Insolar API in accordance with the specification.
   // The Platform uses the basic JSON RPC 2.0 request structure:
   type requestBody struct {
     JSONRPC        string         `json:"jsonrpc"`
     ID             int            `json:"id"`
     Method         string         `json:"method"`
   }

   type requestBodyWithParams struct {
     JSONRPC        string         `json:"jsonrpc"`
     ID             int            `json:"id"`
     Method         string         `json:"method"`
     // Params is a structure that depends on a particular method:
     Params         interface{}    `json:"params"`
   }

   // The Platform defines params of the signed request as follows:
   type params struct {
     Seed            string       `json:"seed"`
     CallSite        string       `json:"callSite"`
     // CallParams is a structure that depends on a particular method.
     CallParams      interface{}  `json:"callParams"`
     PublicKey       string       `json:"publicKey"`
   }

   type paramsWithReference struct {
     params
     Reference       string  `json:"reference"`
   }

   // The member.create request has no parameters, so it's an empty structure:
   type memberCreateCallParams struct {}

   // The transfer request sends an amount of funds to member identified by a reference:
   type transferCallParams struct {
     Amount            string    `json:"amount"`
     ToMemberReference string    `json:"toMemberReference"`
   }

   // Create a function to get a new seed for each signed request:
   func getNewSeed() (string) {
     // Form a request body for getSeed:
     getSeedReq := requestBody{
       JSONRPC: "2.0",
       Method:  "node.getSeed",
       ID:      id,
     }
     // Increment the id for future requests:
     id++

     // Marshal the payload into JSON:
     jsonSeedReq, err := json.Marshal(getSeedReq)
     if err != nil {
       log.Fatalln(err)
     }

     // Create a new HTTP request and send it:
     seedReq, err := http.NewRequest("POST", TestNetURL, bytes.NewBuffer(jsonSeedReq))
     if err != nil {
       log.Fatalln(err)
     }
     seedReq.Header.Set("Content-Type", "application/json")

     // Perform the request:
     seedResponse, err := client.Do(seedReq)
     if err != nil {
       log.Fatalln(err)
     }
     defer seedReq.Body.Close()

     // Receive the response body:
     seedRespBody, err := ioutil.ReadAll(seedResponse.Body)
     if err != nil {
       log.Fatalln(err)
     }

     // Unmarshal the response:
     var newSeed map[string]interface{}
     err = json.Unmarshal(seedRespBody, &newSeed)
     if err != nil {
       log.Fatalln(err)
     }

     // (Optional) Print the request and its response:
     print := "POST to " + TestNetURL +
       "\nPayload: " + string(jsonSeedReq) +
       "\nResponse status code: " +  strconv.Itoa(seedResponse.StatusCode) +
       "\nResponse: " + string(seedRespBody) + "\n"
     fmt.Println(print)

     // Retrieve and return the current seed:
     return newSeed["result"].(map[string]interface{})["seed"].(string)
   }

   // Create a function to send signed requests:
   func sendSignedRequest(payload requestBodyWithParams, privateKey *ecdsa.PrivateKey) map[string]interface{} {
     // Marshal the payload into JSON:
     jsonPayload, err := json.Marshal(payload)
     if err != nil {
       log.Fatalln(err)
     }

     // Take a SHA-256 hash of the payload's bytes:
     hash := sha256.Sum256(jsonPayload)

     // Sign the hash with the private key:
     r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
     if err != nil {
       log.Fatalln(err)
     }

     // Convert the signature into ASN.1 DER format:
     sig := ecdsaSignature{
       R: r,
       S: s,
     }
     signature, err := asn1.Marshal(sig)
     if err != nil {
       log.Fatalln(err)
     }

     // Convert both hash and signature into a Base64 string:
     hash64 := base64.StdEncoding.EncodeToString(hash[:])
     signature64 := base64.StdEncoding.EncodeToString(signature)

     // Create a new request and set its headers:
     request, err := http.NewRequest("POST", TestNetURL, bytes.NewBuffer(jsonPayload))
     if err != nil {
       log.Fatalln(err)
     }
     request.Header.Set("Content-Type", "application/json")

     // Put the hash string into the HTTP Digest header:
     request.Header.Set("Digest", "SHA-256="+hash64)

     // Put the signature string into the HTTP Signature header:
     request.Header.Set("Signature", "keyId=\"public-key\", algorithm=\"ecdsa\", headers=\"digest\", signature="+signature64)

     // Send the signed request:
     response, err := client.Do(request)
     if err != nil {
       log.Fatalln(err)
     }
     defer response.Body.Close()

     // Receive the response body:
     responseBody, err := ioutil.ReadAll(response.Body)
     if err != nil {
       log.Fatalln(err)
     }

     // Unmarshal it into a JSON object:
     var JSONObject map[string]interface{}
     err = json.Unmarshal(responseBody, &JSONObject)
     if err != nil {
       log.Fatalln(err)
     }

     // (Optional) Print the request and its response:
     print := "POST to " + TestNetURL +
       "\nPayload: " + string(jsonPayload) +
       "\nResponse status code: " + strconv.Itoa(response.StatusCode) +
       "\nResponse: " + string(responseBody) + "\n"
     fmt.Println(print)

     // Return the response:
     return JSONObject
   }

   // Create the main function to form and send signed requests:
   func main() {
     // Generate a key pair:
     privateKey := new(ecdsa.PrivateKey)
     privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
     var publicKey ecdsa.PublicKey
     publicKey = privateKey.PublicKey

     // Convert both private and public keys into PEM format:
     x509PublicKey, err := x509.MarshalPKIXPublicKey(&publicKey)
     if err != nil {
       log.Fatalln(err)
     }
     pemPublicKey := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: x509PublicKey})

     x509PrivateKey, err := x509.MarshalECPrivateKey(privateKey)
     if err != nil {
       log.Fatalln(err)
     }
     pemPrivateKey := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: x509PrivateKey})

     // The private key is required to sign requests.
     // Make sure to put into a file to save it in some secure place later:
     file, err := os.Create("private.pem")
     if err != nil {
       fmt.Println(err)
       return
     }
     file.WriteString(string(pemPrivateKey))
     file.Close()

     // Get a seed to form the request:
     seed := getNewSeed()
     // Form a request body for member.create:
     createMemberReq := requestBodyWithParams{
       JSONRPC: "2.0",
       Method:  "contract.call",
       ID:      id,
       Params:params {
         Seed: seed,
         CallSite: "member.create",
         CallParams:memberCreateCallParams {},
         PublicKey: string(pemPublicKey)},
     }
     // Increment the JSON RPC 2.0 request identifier for future requests:
     id++

     // Send the signed member.create request:
     newMember := sendSignedRequest(createMemberReq, privateKey)

     // Put the reference to your new member into a variable to send transfer requests:
     memberReference := newMember["result"].(map[string]interface{})["callResult"].(map[string]interface{})["reference"].(string)
     fmt.Println("Member reference is " + memberReference)

     // Get a new seed to form a transfer request:
     seed = getNewSeed()
     // Form a request body for transfer:
     transferReq := requestBodyWithParams{
       JSONRPC: "2.0",
       Method:  "contract.call",
       ID:      id,
       Params:paramsWithReference{ params:params{
         Seed: seed,
         CallSite: "member.transfer",
         CallParams:transferCallParams {
           Amount: "100",
           ToMemberReference: "<recipient_member_reference>",
           },
         PublicKey: string(pemPublicKey),
         },
         Reference: string(memberReference),
       },
     }
     // Increment the id for future requests:
     id++

     // Send the signed transfer request:
     newTransfer := sendSignedRequest(transferReq, privateKey)
     fee := newTransfer["result"].(map[string]interface{})["callResult"].(map[string]interface{})["fee"].(string)

     // (Optional) Print out the fee.
     fmt.Println("Fee is " + fee)
   }