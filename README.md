# Gophkeeper Password Manager

## Getting Started 

On Mac OS with Apple silicon processor.

```
brew install postgresql@16
brew services start postgresql@16
psql postgres

create database gophkeeper;
create user gophkeeper with encrypted password 'gophkeeper';
grant all privileges on database gophkeeper to gophkeeper;
alter database gophkeeper owner to gophkeeper;
```

Run server from project3/cmd/server
```
go run main.go -d postgresql://gophkeeper:gophkeeper@localhost:5432/gophkeeper -a localhost:8080
```

Run command line client from project3/cmd/client/bin
```
./gophkeeper-darwin-arm64
```

Client sqlite file at project3/cmd/client/bin/gophkeeper.db

## Features

- Data encrypted on client end: credentials, text, files up to 1 MB, bank cards.
- Storage on server and secure sync over multiple clients.
- Auto sign out and brute-force prevention.
- Offline access from signed client.
- Single sqlite file with user data.
- Program for Mac, Linux and Windows.
