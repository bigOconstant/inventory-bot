# GPU stock bot

Tired of not being able to get a 6800xt, this bot will defeat all other slow python selenium bots.

# How to use
Add items to urls section in settings.json. Three items are given as an example. Currently it just parses the page for a "add to cart" string.

**item**: can be what ever name you want. It's just a name.

**url**: is the url to check for the add to cart string. Works with bestbuy

**delayseconds**: delay in seconds between checks. 

Each url is checked in a seperate thread so they can be ran concurrently.

**useragent**: String to tell server what web browser you are using. Spoofing browser