### voluum-apis

This utility can be used to fetch reports from your Voluum account to Google Sheets. 
This script can be used with a cron setup and which pulls out the daily report and pushes to set Google Sheet.

To run this script you just need Docker and config.json which will have your Access Keys and couple of flags file in root directory of project.

config.json should have below mandatory fields:
```
{
    "SpreadsheetId" : "*********************",
    "VoluumAccessId" : "***********************",
    "VoluumAccessKey" : "******************",
    "IncludeTrafficSources" :"ACTIVE",
    "TrafficSourcesShortlisted" :["Advertizer", "Facebook", "Zeropark"],
    "TrafficSourceFilteringEnabled" :false
}
```

```
git clone https://github.com/bhambri94/voluum-apis.git

cd voluum-apis/

vi config.json 
//save the configs shared above with spreadhseet id and voluum access key and id to the file

docker build -t voluum-apis:v1.0 .

docker images ls

docker run -it --name voluum-apis -v $PWD/src:/go/src/voluum-apis voluum-apis:v1.0

```

While we run this project for the first time, we would need Google Account Access Token and Refresh Token. We need to enter a code for the first time while running this project.

Once you run the `docker run` last command shared above, a link will be displayed in the command line, which we need to open in a browser, it will ask for `Allow message` to use your account and grant access to write in the Google Sheet. Once we click the Allow button, a code would be generated, that we need to paste in console, after successful verification a token.json file will be generated at the root directory of the project. 
Note: This file will need to be regenerated if we have created a new Docker build.

###Cron job

To setup a Daily Cron job, please follow following steps:
 
```
cd voluum-apis/

Vi bash.sh

```
```
#!/bin/bash
sudo /usr/bin/docker restart voluum-apis
```

Save the sheet script and run command 

```
chmod 777 bash.sh

Crontab -e

* 9 * * * /path_to_voluum-apis_repo/bash.sh > /path_to_voluum-apis_repo/voluum-apis.logs

```
This above command written in crontab will run the script daily at 9AM UTC time.
