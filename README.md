# Delco Door
This is a Go server that uses Twilio to let authorized users into my apartment, ranther than having to manually answer my phone and let them in. Strategy inspired by Francesc Campoy https://www.youtube.com/watch?v=mTd3hHUy9OU.

# To Use Yourself
1. Get Twilio account, twilio number
2. Get GCE account, install gcloud locally
3. Get API key from Google Speech API
4. Clone repo
5. Make an app.yaml from the template, add your ```SPEECH_API_KEY```
6. Make an app engine instance in a new project folder (might need to switch to that new folder locally ```gcloud config configurations set```...)
7. ```gcloud deploy app``` from this folder
8. Set the public address (ex: https://delco-door.appspot.com) as twilio hook

* to test locally use the below, and change the twilio endpoint to ngrok's 
```
dev_appserver.py app.yaml
ngrok http 8080
```
* the dev server is a gcloud thing

# Notes:
+ This is a go library. If you want to deploy this manually, just change it back to a regular package (init -> main), ignore the app.yaml stuff, and make a compute engine or AWS instance. I have instructions for dockerizing go applciations in a couple steps here if containerizing would be easier https://github.com/max-hoffman/go-starter

# Password
It's in the code. I'll make it a private env variable if that becomes a problem.

# Changing Password (in progress)
The password for the door is public and can be changed with a "king of the hill" type Ethereum smart contract. In non-technical terms, people can outbid each other to change the password, until the end of the semester when the change price is set back to zero.
