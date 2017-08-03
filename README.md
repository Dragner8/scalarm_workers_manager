[![Build Status](https://travis-ci.org/Dragner8/scalarm_workers_manager.svg?branch=master)](https://travis-ci.org/Dragner8/scalarm_workers_manager)   [![Codacy Badge](https://api.codacy.com/project/badge/Grade/cdb30f2077894dd7953060f443a009ec)](https://www.codacy.com/app/Dragner8/scalarm_workers_manager?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=Dragner8/scalarm_workers_manager&amp;utm_campaign=Badge_Grade)

Scalarm Workers Manager
============ 

Installation guide: 
---------------------- 
Go 
-- 
To build and install Scalarm Workers Manager you need to install go programming language. 
You can install it from official binary distribution: 

https://golang.org/doc/install

or from source: 

https://golang.org/doc/install/source 

After that you have to specify your $GOPATH. Read more about it here: 

https://golang.org/doc/code.html#GOPATH 

Installation 
-------------- 
You can download Scalarm Workers Manager directly from GitHub. You have to download it into your $GOPATH/src folder 
``` 
go get github.com/scalarm/scalarm_workers_manager
``` 
Now you can install Scalarm Workers Manager: 
```` 
go install github.com/scalarm/scalarm_workers_manager
```` 
This command will install Scalarm Workers Manager in $GOPATH/bin. It's name will be scalarm_workers_manager.

Config 
-------- 
The config folder contains single file config.json that contains required informations for Scalarm Workers Manager:

* InformationServiceAddress - address of working Information Service
* Login, Password - Scalarm credentials
* Infrastructures - list of infrastructures monitor has to check for records
* ScalarmCertificatePath - path to custom certificate (optional, by default looking in standard certificate directory)
* ScalarmScheme - http or https (default)
* InsecureSSL - should invalid certificates (eg. self-signed) be accepted
* ProbeFrequencySecs (optional, default: 10) - number of delay in seconds between getting records and checking their state if last check was completed
* ExitTimeout (optional, default: 0) - number of seconds to wait before exiting if there are no more records available; if this value is greater than 0, the minimal time to exit is ProbeFrequencySecs; if this values is less than 0, process continues work indefinitely (use with caution!)
* VerboseMode - should more detailed logs be shown

Example config:

```
{
	"InformationServiceAddress": "149.156.10.32:31034",
	"Login": "login",
	"Password": "password",
	"Infrastructures": 
	[
		"qsub",
		"qcg"
	],
	"ScalarmCertificatePath" : "cert.pem",
	"ScalarmScheme" : "https",
	"InsecureSSL" : true,
	"ProbeFrequencySecs": 10,
	"ExitTimeout": 0,
	"VerboseMode": false
}
```
Example config can be found in config/config.json.

Run 
---- 
Before running program you have to copy contents of config folder to folder with executable file of Scalarm Workers Manager. By default it will be $GOPATH/bin 

