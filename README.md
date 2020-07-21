voluum-apis
This utility can be used to fetch reports from your Voluum account to Google Sheets. 
This script can be used with a cron setup and which pulls out the daily report and pushes to set Google Sheet.

To run this script you just need Docker and config.json which will have your Access Keys and couple of flags file in root directory of project.

```
git clone https://github.com/bhambri94/voluum-apis.git

cd voluum-apis/

cp config.json .

docker build -t voluum-apis:v1.0 .

docker images ls

docker run -it -v $PWD/src:/go/src/voluum-apis voluum-apis:v1.0

```

On the first run for this utility, it needs to generate a token.json file which is actually access_token and refresh tokens to write to the mentioned Google Sheet. 
Once you run the build first time, a link will be displayed in command line, which we need to open in a browser, it will ask for Allow message from your gmail account and will share a 
Access id, which we need to paste in command line and token.json file will be generated, which can be used going forward.
