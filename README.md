# SSH TUNNEL

Test project used to connect to private database servers using a stepping-stone server

## Install the command

```
go install github.com/mehix/sshtunnel/cmd/sshtunnel
```

## Run

```
sshtunnel <local-addr> <stepping-stone> <pk-file> <remote-private-addr>
```

**local-addr**: can be IP:PORT but also IP:0 and the port will be random and displayed at runtime
**stepping-stone**: user@server
**pk-file**: ssh private key used to connect to the stepping-stone server. The password can be provided in a file called `.secret`
**remote-private-addr**: IP:PORT of the server we actually want to connect to
