# GPU stock bot

Tired of not being able to get a 6800xt, this bot will defeat all other slow python selenium bots.

# How to use
Add items to urls section in settings.json. Three items are given as an example. Currently it just parses the page for a "add to cart" string.

**item**: can be what ever name you want. It's just a name.

**url**: is the url to check for the add to cart string. Tested with bestbuy and newegg

**delayseconds**: delay in seconds between checks. Longer delays could lessen chances of a ip ban.

Each url is checked in a seperate thread.

**useragent**: String to tell server what web browser you are using. Spoofing browser

## How to run

`go run main.go`

Open up the UI on `localhost:3000`

UI refreshes every 1 second and updates the table with all the items. If an item suddendly becomes in stock a single notification pops up. Only one notification appears but in stock column gets updated. 