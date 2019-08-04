# friendbug

Friendbug is an app that provides you WhatsApp messages so you can contact friends.

It uses Golang as a main programming language, Redis as a primary database and Twilio for the integration with WhatsApp.

I was used primarily to help me keep in touch with friends, and to try out Golang and Redis

## Setup

To get things started, you need the following things:

* [A Twilio account](https://www.twilio.com/login)
* Your WhatsApp linked to the Twilio account
* A heroku account
* A persistent instance (e.g. [Redis to go](https://elements.heroku.com/addons/redistogo) )
* A scheduler to trigger the code (e.g. [Advanced scheduler](https://elements.heroku.com/addons/advanced-scheduler) )
* A list of friends to contact in the redis instance in the following format
    * Redis list name: `friends`
    * Friends structure in a JSON object
        * name (name of the friend)
        * dateContacted (the last day that you contacted a friend)
        * contactFrequency (integer, plus minus how many days you want to contact this friend)
* An .env file with the following variables
    * TWILIO_ACCOUNT_SID
    * TWILIO_AUTH_TOKEN
    * TWILIO_SOURCE_NUMBER (the Twilio number)
    * TARGET_NUMBER (your number)
    * REDISTOGO_URL
    
    
You can also host this on another platform than Heroku, but then the deployment scripts are not provided.

## Installation

Prerequisites:

* A Heroku account
* A Redis account
* A Twilio account

Actual installation:

1. Clone the repository
2. Create the .env file with the correct variables
3. Create a `friends` list in Redis, and populate it with friends
4. Deploy to Heroku
5. Run a one-off dyno with `friendbug`

To run it locally, go `go install` and then execute `friendbug`.
    
## Notes

**Important**: Because the WhatsApp of Twilio is still in Beta, it cannot send a number a message unless that number has sent the Twilio number a message in the last 24 hours. I suggest you sent `Thank you` back every time it reminds you to contact friends :).
