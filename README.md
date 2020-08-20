[![Build Status](https://travis-ci.org/bostontrader/okprobe.svg?branch=master)](https://travis-ci.org/bostontrader/okprobe)
[![MIT license](http://img.shields.io/badge/license-MIT-brightgreen.svg)](http://opensource.org/licenses/MIT)

# Welcome to OKProbe
OKProbe is a command line tool that will enable the user to easily invoke the OKEx API.  The user can specify API endpoints and parameters that he expects should succeed, or he can request that OKProbe submit a test suite of erroneous requests, in order to observe the reaction of the OKEX API under these circumstances. The user can also specify a variety of variations of the output format.

## Usage
```
./okprobe
```
When executed without any command line params OKProbe will print a help screen.  You can disect this at your leisure but the general form of any command is ./okprobe [command] --flags.

Some of the flags are global and apply to all sub commands and other flags are specific to a specific sub command.

For example:

```
./okprobe accountWallet --credentialsFile /home/myhome/okex-read.json
```

This command will invoke the OKProbe on the default url https://www.okex.com and request the /api/account/v3/wallet endpnt.  This endpoint requires no input parameters so there are none sent.  It will not intentionally generate any error conditions nor will it worry about the contents of the headers.  It will print a single string of text that should be parsable into JSON.
  
Please realize that the OKProbe knows the actual URL path in order to this, so we're not specifying that anywhere here. Notice the use of a file, specified by a global flag, to hold the API credentials.

## Digging Deeper...

The ordinary operation of the OKProbe is to simply submit a request to the OKEx server and get and print the response. Said response should be parsable into JSON and can thus easily be included in other testing.  If there's an error then print an error message and exit with status code 1.
  
This ordinary pattern can be tweaked via a variety of flags.

### --baseURL

By default all requests go to https://www.okex.com.  If you have some reason to use another URL, such as that of an [OKCatbox](https://github.com/bostontrader/okcatbox), use this global flag.

Note that this is the _base_ URL.  The rest of the path to specify a particular endpoint and query string gets appended later.

### --credentialsFile

This global flag points to a file that contains the desired credentials in this example format:

{"api_key":"f28f7845-1b85-4463-4378-8f6cc442efdf","api_secret_key":"44F60E04780B71427C8A67E4AF8D6633","passphrase":"valid passphrase"}

(We invite hackers to try to steal something using these credentials since they have been generated randomly to look real and have nothing to do with any real OKEx account.)


### --makeErrorsCredentials

The API calls require four headers in the request dealing with the credentials.  There are a variety of errors that can be made doing this, and the OKEx server will respond with expected http status codes and error messages.  If you use this global flag then the OKProbe will submit several requests to the given endpoint and test that the OKEx server responds as expected, in addition to any other error testing and the final ordinary operation.

If the tests pass then the OKProbe will exit with a status code 0, else the OKProbe will exit with status code 1.

### --makeErrorsParams

Some of the API calls expect/allow parameters.  If you use this global flag then the OKProbe will submit several requests to the given endpoint containing a variety of errors and test that the OKEx server responds as expected, in addition to any other error testing and the final ordinary operation.

If the tests pass then the OKProbe will exit with a status code 0, else the OKProbe will exit with status code 1.

If an endpoint does not have any parameters then the associated probe silently ignores this option.

### -- makeErrorsWrongCredentialsType

OKEx credentials come in three flavors: read only, read & trade, read & withdraw.  If you want to test that some credentials are not of the correct type for a particular endpoint and will thus fail if used, use this global flag.  

In this case the ordinary operation of the OKProbe will be attempted using whatever credentials specified via --credentialsFile.  The user expects that the API invocation will fail because of a credential type error. If so this command will exit with status code 0, else status code 1.

### --queryString

Any GET endpoint that can accept parameters expects to receive them via a query string.  The user should make a query string and pass it to OKProbe via this flag.  For example:

--queryString "?param1=A&param2=B"

If an endpoint does not have any parameters then the associated probe silently ignores this option.
 
### Fun with headers.

The responses from the OKEx server have a variety of headers that contain undocumented and mostly mysterious and unknowable things.  In actual practice we have yet to find any practical value in understanding this more deeply.  In the beginning we attempted to observe and test these headers.  Unfortunately, these headers change so frequently that it's the Devil's work to keep on top of them.  The bottom line is that this area is too much trouble to figure out.  There's still some legacy code related to observing these headers and checking for unexpected or omitted headers, but we won't develop this further at this time.
  