# line_bot_test_app_v2

This is a bot that I created to learn the LINE Messaging API. It is written in Go.

## Environment Variables

The following environment variables must be set for the bot to run properly:

`BOT_HOST`: This should be set to the bot's hostname

`BETA_LINE_CHANNEL_SECRET`: This should be set to your channel's CHANNEL SECRET for the Beta environment that is found in the Channel Console.

`REAL_LINE_CHANNEL_SECRET`: Same as above but for the real environment. This is only used if `USE_REAL_ENVIRONMENT` is set to `TRUE`

`BETA_LINE_CHANNEL_ACCESS_TOKEN`: This should be set to the Channel Access Token for your bot's channel in the Beta environment. 

`REAL_LINE_CHANNEL_ACCESS_TOKEN`: Same as above, but for the real environment. This is only used if `USE_REAL_ENVIRONMENT` is set to `TRUE`

`USE_REAL_ENVIRONMENT`: If this is set to `TRUE`, the bot will use the endpoints for the Real environment. If this is variable is not set or not set to `TRUE`, the bot will default to using the Beta environment.

## Dependency Management

This project uses godep to manage its external dependencies.

https://github.com/tools/godep

Running `godep save` in the src directory will register the currently used in the project to the godeps f

## Bot Functionality

This bot has the following *Trivial* functions:

### Buttons Template Message ("find zombie")
* If a user says "find zombie", the bot will send a buttons template message with a picture of a zombie and three options: "Run!", "Scream!" "EXPLODE".
** "Run!" Sends a *postback* event to the bot. The bot will handle the event and randomly determine if the user escaped or exploded. It will send a message telling the user the result.
** "Scream!" causes the user to say "AHHH!" in the chat with the bot.
** "Explode!" opens a url to a picture of an exploding kitten.

### Carousel Template Message ("multizombie")
* If a user says "multizombie", the bot will send a carousel that contains multiple instances of the Button Template Message

### Confirm Template Message ("explode")
* If the user says "explode", the bot sends a confirm dialog asking of the user wants to explode.
** If the user chooses "YES!", the user is sent to a picture of an exploding kitten.
** If the user chooses "NO!", a postback event is sent to the bot and it responds calling the user a coward.

### Image Map ("imagemap")
* If the user says "imagemap", the bot sends an imagemap with two choices.
** If the user clicks on "CATS", they will be sent to the official Exploding Kittens website.
** If the user clicks on "ZOMBIES", the user will say "ZOMBIES!" in the chat with the bot.

### Push Message ("push")
* If the user says "push", the bot will send a push message directly to the user without using a Reply Token.



