# World of Warcraft Market Trends
A web application plus a collection of tools for importing World of Warcraft auction house data and tracking price trends. See a live demo at [https://wow.open-mailbox.com/](https://wow.open-mailbox.com/).

## Dependencies
* [Go](https://golang.org/) (tested on 1.9)
* [PostgreSQL](https://www.postgresql.org/) (tested on 9.5)

## Development
1. `git clone`
2. `cd wow_market_trends`
3. `make`
4. `cd cmd/wowexchange`
5. `./wowexchange` to fetch the latest auction house snapshot
6. `./wowexchange serve` to start a webserver locally on port 8081.

In order to fetch the latest auction house snapshot data, a Blizzard API key is required. You can obtain one by registering a free account at the [Blizzard Devloper Portal](https://dev.battle.net/). Set the `BLIZZARD_API_KEY` in your environment using the obtained key.

Change the constants in `server.go` for a different port or logging directory. Change the constant in `auctions.go` for a different WoW realm.

## Contributing
1. Fork it ([http://github.com/openmailbox/wow_market_trends/form](http://github.com/openmailbox/wow_market_trends/fork))
2. Create your feature branch (git checkout -b my-new-feature)
3. Commit your changes (git commit -am 'Add some feature')
4. Push to the branch (git push origin my-new-feature)
5. Create new Pull Request
