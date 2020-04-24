# Welcome to OKProbe
The purpose of OKProbe is to probe the OKEx API using a variety of correct and incorrect methods in order to observe its behavior.

In order to use said API it's generally helpful to understand how exactly it behaves.  RTFM is a good first start, but there's always extra devils in the details.  OKProbe gives you a means of carving your hard-fought field-knowledge of the actual API behavior into blocks of solid software.

One important difficulty with this general task is that the real OKEx API frequently requires credentials, and you may be hesitant to use your real OKEx credentials when taunting the real OKEx API for practice.  OKProbe addresses this situation thus:

## Store the Credentials in a File

The problem is not so much that OKEx will retaliate against your real credentials for using the Probe, although they might conceivably do so given excessive probing.  The real problem is how to get your sensitive credentials into the running program without hard-wiring them into code or command-line history or any other similar faux pas.  This is a deeper question then anything we can address here, but our answer-of-expediency is to hardware a file.json name, in your home directory, into the Probe's code and to let you configure and guard said file.json as you wish.

## OKCatbox

Another solution is to use a simulated OKEx server for your Probing needs.  We happen to know of one at https://github.com/bostontrader/OKCatbox.  With the OKCatbox you can create your own credentials that you can freely use with the OKProbe.  This will give you a reasonable first sense of correct operation of your code, before the time comes to really unleash it upon the real OKEx API.  Please realize that the credentials created by the OKCatbox are in no way connected to your real credentials and cannot be used on the real OKEx API.