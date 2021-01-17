# Inventory stock bot

Tired of not being able to get a 6800xt, this bot will defeat all other slow python selenium bots.
compiles to a single executable for easy deployment. All html and images are compiled within.  

# How to use
Add items to urls section in settings.json. Three items are given as an example. Currently it just parses the page for a "add to cart" string.

**item**: can be what ever name you want. It's just a name.

**url**: is the url to check for the add to cart string. Tested with bestbuy and newegg

**delayseconds**: delay in seconds between checks. Longer delays could lessen chances of a ip ban.

Each url is checked in a seperate thread.

**useragent**: String to tell server what web browser you are using. Spoofing browser

**discord**: Optional field. If given a discord webhook a message will be sent to the webhook channel informing of in stock items.

## Example settings.json

```json

{
    "delayseconds" : 30,
    "useragent" : "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100 Safari/537.36",
    "discord":"",
    "urls" : [
     {"item":"6800xt xfx","url":"https://www.bestbuy.com/site/xfx-amd-radeon-rx-6800xt-16gb-gddr6-pci-express-4-0-gaming-graphics-card-black/6441226.p?skuId=6441226"},
     {"item":"6800xt msi-radeon","url": "https://www.bestbuy.com/site/msi-radeon-rx-6800-xt-16g-16gb-gddr6-pci-express-4-0-graphics-card-black-black/6440913.p?skuId=6440913"},
     {"item":"new egg 6800xt other","url":"https://www.newegg.com/gigabyte-radeon-rx-6800-xt-gv-r68xt-16gc-b/p/N82E16814932373?"}
    ]
}

```

## How to run manually (requires golang and make)

`make`

`./goinventory <Port>`

Open up the UI on `localhost:Port`

UI refreshes every 1 second and updates the table with all the items. If an item suddendly becomes in stock a single notification pops up. If a discord webhook was given a discord message will be sent.


## Run in Docker


**build**

`docker build -t goinventory`

**Run**

```bash

docker build -t goinventory .;

docker run -d \
  -it \
  -p 3000:3000 \
  -v "$(pwd)"/settings.json:/app/settings.json:ro \
  goinventory
  
```

## Run in docker-compose (Suggested method) (one line)

`docker-compose up -d`


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


