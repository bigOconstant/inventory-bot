# Inventory stock bot

Tired of not being able to get a 6800xt, this bot will defeat all other slow python selenium bots.
compiles to a single executable and needs only a single json file in it's directory to run run. 

# How to use
Add items to urls section in settings.json. Three items are given as an example. Currently it just parses the page for a "add to cart" string.

**item**: can be what ever name you want. It's just a name.

**url**: is the url to check for the add to cart string. Tested with bestbuy and newegg

**delayseconds**: delay in seconds between checks. Longer delays could lessen chances of a ip ban.

**host**: host name to be used in the UI.

**port**: port to be used in the UI.

Each url is checked in a seperate thread.

**useragent**: String to tell server what web browser you are using. Spoofing browser

**discord**: Optional field. If given a discord webhook a message will be sent to the webhook channel informing of in stock items.

## How to run

`go run main.go <Port>`

Open up the UI on `localhost:Port`

UI refreshes every 1 second and updates the table with all the items. If an item suddendly becomes in stock a single notification pops up. If a discord webhook was given a discord message will be sent.


## Run in Docker


**build**

`docker build -t goinventory`

**Run**

```bash

docker run -d \
  -it \
  -p 3000:3000 \
  -v "$(pwd)"/settings.json:/app/settings.json:ro \
  goinventory
  
```

## Run in docker-compose (one line)

`docker-compose up -d`


## forward ip 3000 to 80 

I forward 3000 to port 80 for deployment on my raspberrypi. so I don't need to specify a port in the browser

`sudo iptables -A INPUT -i eth0 -p tcp --dport 80 -j ACCEPT`

`sudo iptables -A INPUT -i eth0 -p tcp --dport 3000 -j ACCEPT`

`sudo iptables -A PREROUTING -t nat -i eth0 -p tcp --dport 80 -j REDIRECT --to-port 3000`


## To do

Work on a few important one off parsers for pages like https://www.amd.com/en/direct-buy/us where stock pops in for multiple items at a time. 
Will likely need https://pkg.go.dev/golang.org/x/net/html
