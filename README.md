# Newman `(งツ)ว`

Newman is a discord bot that aims to have functionalities that ultimately revolve around the hit-comedy show, Seinfeld, and the hilarious Wayne Knight (AKA "Newman") in particular. This is fun side-project that I hope to maintain and have folks collaborate on in the future.

_Feel free to chime in!_

![The gang](./assets/seinfeld.jpg)

***
#### Development
- The bot uses `mssql` for database management. To trigger an instance run,
```bash
docker run \
  -e "ACCEPT_EULA=Y" \
  -e "SA_PASSWORD=Qwertyuiop1#" \
  -p 1433:1433 \
  --name mssql \
  -h mssql \
  -d mcr.microsoft.com/mssql/server:2019-latest
```
- Replace the token placeholder in `internal/private.json` with a valid one and run `make run` to start the bot.
