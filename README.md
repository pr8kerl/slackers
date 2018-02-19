# slackers

compare SSO slack users to an LDAP source

## build

```
docker-compose run make
```

# configure it

* copy config.example to config.ini
* edit the values in config.ini for your ldap/AD server
* set a couple of secrets as environment variables

  ```
  SLACK_API_TOKEN="<my slack api token>"
  LDAP_PASSWD=SuperSecretSquirrel1

  export SLACK_API_TOKEN LDAP_PASSWD
```

## run

  ```
  ./slackers --help
  Usage: slackers [--version] [--help] <command> [<args>]

  Available commands are:
      disabled    scan for all disabled AD users that are live slackers
      report      list all live slackers
  ```

