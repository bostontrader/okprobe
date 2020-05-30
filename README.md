[![Build Status](https://travis-ci.org/bostontrader/okprobe.svg?branch=master)](https://travis-ci.org/bostontrader/okprobe)
[![MIT license](http://img.shields.io/badge/license-MIT-brightgreen.svg)](http://opensource.org/licenses/MIT)

# Welcome to OKProbe
The purpose of OKProbe is to probe the OKEx API using a variety of correct and incorrect methods in order to observe its behavior.

In order to use said API it's generally helpful to understand how exactly it behaves.  RTFM is a good first start, but there's always extra devils in the details.  OKProbe gives you a means of carving your hard-fought field-knowledge of the actual API behavior into blocks of solid software.

## Usage
```
./okprobe
```
When executed without any command line params it will print a help screen.

A functioning example:

```
./okprobe -url https://www.okex.com -endpnt wallet -errors -keyfile /home/myhome/okex-read.json

```
This command will invoke the Probe on the given url and request the wallet endpnt.  Consider that the Probe knows the actual URL path in order to this, we're not specifying that anywhere here. The command also specifies that a variety of errors be intentionally created in the calls.  Finally, notice the use of a file to hold the API credentials. Speaking of...

## Credentials

One important difficulty with this general task is that the real OKEx API frequently requires credentials, and you may be hesitant to use your real OKEx credentials when taunting the real OKEx API for practice.  OKProbe addresses this situation thus:

### Store the Credentials in a File

The problem is not so much that OKEx will retaliate against your real credentials for using the Probe, although they might conceivably do so given excessive probing.  The real problem is how to get your sensitive credentials into the running program without hard-wiring them into code or command-line history or any other similar faux pas.  This is a deeper question then anything we can address here, but our answer-of-expediency is to use a key file in a location of your choice, which thou shalt guard carefully.


### OKCatbox

Another solution is to use a simulated OKEx server for your Probing needs.  We happen to know of one at https://github.com/bostontrader/OKCatbox.  With the OKCatbox you can create your own credentials that you can freely use with the OKProbe.  This will give you a reasonable first sense of correct operation of your code, before the time comes to really unleash it upon the real OKEx API.  Please realize that the credentials created by the OKCatbox are in no way connected to your real credentials and cannot be used on the real OKEx API.

## Testing

Automated testing works by installing [OKCatbox](https://github.com/bostontrader/okcatbox) on the CI server and executing the probes against that.  Doing this has the pleasant side effect of testing OKCatbox as well.

### Limits of Testing

An important limitation of these tests is that we don't carefully inspect the 200 type results.  There are issues with doing so that are too tedious to even describe here, much less implement in the software.  Fortunately there is another layer of testing that starts with a blank db, submits API calls, and examines the specific correct responses at that time.  This is the suitable place for such testing.