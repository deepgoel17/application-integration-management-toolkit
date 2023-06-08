## integrationcli certificates delete

Delete a certificate from a region

### Synopsis

Delete a certificate from a region

```
integrationcli certificates delete [flags]
```

### Options

```
  -h, --help          help for delete
  -n, --name string   Integration flow name
```

### Options inherited from parent commands

```
  -a, --account string   Path Service Account private key in JSON
      --api api          Sets the control plane API. Must be one of prod, staging or autopush; default is prod
      --disable-check    Disable check for newer versions
      --no-output        Disable printing all statements to stdout
      --print-output     Control printing of info log statements (default true)
  -p, --proj string      Integration GCP Project name
  -r, --reg string       Integration region name
  -t, --token string     Google OAuth Token
      --verbose          Enable verbose output from integrationcli
```

### SEE ALSO

* [integrationcli certificates](integrationcli_certificates.md)	 - Manage certificates used by Integration

###### Auto generated by spf13/cobra on 29-Apr-2023