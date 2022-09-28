#### Development
- Replace the token placeholder in `internal/private.json` with a valid one (you can create one [here](https://discord.com/developers/applications/984883201682067517/bot)) and run `make run` to start the bot.
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
or, just simply use the `make run` target for a no-hassle setup.
