# bytebot-irc

## Running bytebot with docker-compose
First, let's configure the bot with a .env file, supposing that an IRC server runs at irc.example.com/6697, fill your .env as follows:

```
BYTEBOT_NICK="your-bots-nick"
BYTEBOT_SERVER="irc.example.com:6697"
BYTEBOT_CHANNELS="#test,#bytebot"
```

Then, we can run the main gateway with `docker-compose up`, the redis server will be accessible on localhost's port 6379, and in the #test and #bytebot channels of the server.

## Writing an application
You can find a commented reaction application in the examples, named `reaction`.
