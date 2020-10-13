[![Build Status](https://travis-ci.org/bostontrader/okprobe.svg?branch=master)](https://travis-ci.org/bostontrader/okprobe)
[![MIT license](http://img.shields.io/badge/license-MIT-brightgreen.svg)](http://opensource.org/licenses/MIT)

# Welcome to OKProbe
OKProbe is a command line tool that will enable the user to easily invoke the OKEx API.  The user can specify API endpoints and parameters that he expects should succeed, or he can request that OKProbe submit a test suite of erroneous requests, in order to observe the reaction of the OKEX API under these circumstances. The user can also specify a variety of variations of the output format.

## Usage
```
./okprobe
```
When executed without any command line params OKProbe will print a help screen.  You can dissect this at your leisure, but the general form of any command is ./okprobe [command] --flags.

Some of the flags are global and apply to all sub commands and other flags are specific to a specific sub command.

For example:

```
./okprobe accountWallet --credentialsFile /home/myhome/okex-read.json --forReal
```

This command will invoke the OKProbe on the default url https://www.okex.com and request the /api/account/v3/wallet endpoint.  This endpoint requires no input parameters so there are none sent.  It will not intentionally generate any error conditions.  It will print a single string of text that should be parsable into JSON.
  
Please realize that the OKProbe knows the actual URL path to OKEx in order to this, so we're not specifying that anywhere here. Notice the use of a file, specified by a global flag, to hold the API credentials.

The --forReal global flag directs the OKProbe to _really_ make a real request to the server.  At first glance this flag may seem extraneous.    However, it's worth will become evident soon.

## Digging Deeper...

The OKProbe is useful because it can submit a valid request to the server in order to provoke a sensible result.  However, in addition to this ordinary operation we might also want to submit erroneous requests in order to observe the server's behavior.  Upon further reflection we have identified three classes of erroneous requests that are potentially interesting, so the OKProbe provide the means to individually specify these test suites by a global flag.

Some of the endpoints of the OKEx server are idempotent, and some produce changes on the server.  Because we might want to test only erroneous requests and not actually try to get a real result, said real result will only be produced if the user specifies the --forReal global flag.  Although this is perhaps not real important for the idempotent endpoints, it becomes much more important for endpoints that produce real change on the OKEx server.

That said, the global flag available are:

### --baseURL

By default all requests go to https://www.okex.com.  If you have some reason to use another URL, such as that of an [OKCatbox](https://github.com/bostontrader/okcatbox), use this global flag.

Note that this is the _base_ URL.  The rest of the path to specify a particular endpoint and query string gets appended later.

### --credentialsFile

This global flag points to a file that contains the desired credentials in this example format:

{"api_key":"f28f7845-1b85-4463-4378-8f6cc442efdf","api_secret_key":"44F60E04780B71427C8A67E4AF8D6633","passphrase":"valid passphrase"}

(We invite hackers to try to steal something using these credentials since they have been generated randomly to look real and have nothing to do with any real OKEx account.)


### --forReal

This flag will direct the OKProbe to actually attempt a real call to the OKEx server.


### --makeErrorsCredentials

The API calls require four headers in the request dealing with the credentials.  There are a variety of errors that can be made doing this, and the OKEx server will respond with expected http status codes and error messages.  If you use this global flag then the OKProbe will submit several erroneous requests to the given endpoint and test that the OKEx server responds as expected, in addition to any other operation specified by other flags.


### --makeErrorsParams

Some of the API calls expect/allow parameters.  If you use this global flag then the OKProbe will submit several requests to the given endpoint containing a variety of errors and test that the OKEx server responds as expected, in addition to any other operation specified by other flags.

If an endpoint does not have any parameters then the associated probe silently ignores this option.


### -- makeErrorsWrongCredentialsType

OKEx credentials come in three flavors: read only, read & trade, read & withdraw.  If you want to test that some credentials are not of the correct type for a particular endpoint and will thus fail if used, use this global flag.  

### --postBody

Any POST endpoint that can accept parameters expects to receive them via the request body.  The user should make a suitable string and pass it to the OKProbe via this flag.  For example:

--postBody `{"from":"6", "to":"1", "amount":"0.1", "currency":"btc"}`

If an endpoint does not have any parameters then the associated probe silently ignores this option.
 
### --queryString

Any GET endpoint that can accept parameters expects to receive them via a query string.  The user should make a query string and pass it to the OKProbe via this flag.  For example:

--queryString "?param1=A&param2=B"

If an endpoint does not have any parameters then the associated probe silently ignores this option.
 

### Fun with headers.

The responses from the OKEx server have a variety of headers that contain undocumented and mostly mysterious and unknowable things.  In actual practice we have yet to find any practical value in understanding this more deeply.  In the beginning we attempted to observe and test these headers.  Unfortunately, these headers change so frequently that it's the Devil's work to keep on top of them.  The bottom line is that this area is too much trouble to figure out, relative to the benefit in doing so.  There's still some legacy code related to observing these headers and checking for unexpected or omitted headers, but we won't develop this further at this time.
