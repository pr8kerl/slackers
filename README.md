# slackers

Find SSO slack users that should be deactivated according to an LDAP source

## build

```
docker-compose run make <darwin/linux/windows>
```

## configure

* copy **config.example** to config.ini
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

                  -purge ...let the app to kill the targeted users from Slack
                  -o json ...outputs result in JSON that is machine friendly
                  -o stdout ...outputs result in stdout (default)

      report      list all live slackers
  ```

