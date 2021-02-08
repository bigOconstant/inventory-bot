# Inventory stock bot

Tired of not being able to get a 6800xt, this bot will defeat all other slow python selenium bots.
compiles to a single executable for easy deployment. All html and images are compiled within.  

## How to use
Add items in the Add Item Tab. Currently it just parses the page for a "add to cart" string.

**item**: can be what ever name you want. It's just a name.

**url**: is the url to check for the add to cart string. Tested with bestbuy and newegg

## Settings page

**delayseconds**: delay in seconds between checks. Longer delays could lessen chances of a ip ban.

Each url is checked in a seperate thread.

**useragent**: String to tell server what web browser you are using. Spoofing browser

**Discord Webhook**: Optional field. If given a discord webhook a message will be sent to the webhook channel informing of in stock items.


## How to run manually (requires golang,gcc and make)

`make`

`./goinventory <Port>`

Open up the UI on `localhost:Port`

UI refreshes every 1 second and updates the table with all the items. If a discord webhook was given a discord message will be sent.

## Run in docker-compose (Suggested method) (one line)

`docker-compose up -d`

Database file is stored in a file called inventory.db in the container. You can back it if you need to delete or recreate the container. 

## forward ip 3000 to 80 

I forward 3000 to port 80 for deployment on my raspberrypi. so I don't need to specify a port in the browser

`sudo iptables -A INPUT -i eth0 -p tcp --dport 80 -j ACCEPT`

`sudo iptables -A INPUT -i eth0 -p tcp --dport 3000 -j ACCEPT`

`sudo iptables -A PREROUTING -t nat -i eth0 -p tcp --dport 80 -j REDIRECT --to-port 3000`


## To do

Work on a few important one off parsers for pages like https://www.amd.com/en/direct-buy/us where stock pops in for multiple items at a time. 
Will likely need https://pkg.go.dev/golang.org/x/net/html


## Screenshots

![screenshot1](screenshots/screenshot1.png?raw=true "screenshot1")

![screenshot2](screenshots/screenshot2.png?raw=true "screenshot2")

![screenshot3](screenshots/screenshot3.png?raw=true "screenshot3")

![screenshot4](screenshots/screenshot4.png?raw=true "screenshot4")


