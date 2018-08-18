# World of Warcraft Market Trends
A web application plus a collection of tools for importing WoW auction house data and tracking price trends.
## Development
1. `git clone`
2. `cd wow_market_trends/internal && go build`
3. `cd wow_market_trends/cmd/wowexchange && go build`

In order to fetch the latest auction house snapshot data, a Blizzard API key is required. You can obtain one by registering a free account at the [Blizzard Devloper Portal](https://dev.battle.net/). Set the `BLIZZARD_API_KEY` in your environment using the obtained key.

Run `wowexchange` to import the latest snapshot of auction house data and create the current price data for all unique items. Run `wowexchange serve` to start the web server on localhost port 8081. 

Change the constants in `server.go` for a different port or logging directory. Change the constant in `auctions.go` for a different WoW realm.

## Contributing
1. Fork it ([http://github.com/openmailbox/wow_market_trends/form](http://github.com/openmailbox/wow_market_trends/form))
2. Create your feature branch (git checkout -b my-new-feature)
3. Commit your changes (git commit -am 'Add some feature')
4. Push to the branch (git push origin my-new-feature)
5. Create new Pull Request