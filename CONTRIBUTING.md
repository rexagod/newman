#### Development
- Replace the token placeholder in `internal/private.json` with a valid one (you can create one [here](https://discord.com/developers/applications/984883201682067517/bot)) and run `make run` to start the bot.
	- Please note that **this is the recommended way**.
- The bot uses `mssql` for database management. To trigger an instance run,
```bash
docker run \
	--cap-add SYS_PTRACE \
	--name mssql \
	-e 'ACCEPT_EULA=Y' \
	-e 'MSSQL_SA_PASSWORD=Qwertyuiop1#' \
	-p 1433:1433 \
	-h mssql \
	-d mcr.microsoft.com/azure-sql-edge
```
and then a `go run .` to start the bot.
