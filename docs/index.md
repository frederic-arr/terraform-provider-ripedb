---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "ripedb Provider"
subcategory: ""
description: |-
  RIPE DB.
---

# ripedb Provider

RIPE DB.

## Example Usage

```terraform
provider "ripe" {}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `certificate` (String, Sensitive) the certificate of the MNTNER object
- `database` (String) the source of the RIPE DB
- `endpoint` (String) the endpoint of the RIPE DB
- `key` (String, Sensitive) the private key of the MNTNER object
- `password` (String, Sensitive) the password of the MNTNER object
- `user` (String) the name of the MNTNER object
