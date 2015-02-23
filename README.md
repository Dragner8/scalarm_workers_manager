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
git clone https://github.com/Scalarm/scalarm_workers_manager.git
``` 
Now you can install monitoring: 
```` 
go install scalarm_workers_manager 
```` 
This command will install monitoring in $GOPATH/bin. It's name will be scalarm_workers_manager.

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
	"ExitTimeout": 0
}
```
Example config can be found in config/config.json.

Run 
---- 
Before running program you have to copy contents of config folder to folder with executable file of Scalarm Workers Manager. By default it will be $GOPATH/bin 

