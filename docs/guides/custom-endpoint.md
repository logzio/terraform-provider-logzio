## Using custom_api_url

The `custom_api_url` provider argument allows you to override the default Logz.io API endpoint. This is useful if you need to route API requests through an internal gateway, proxy, or custom endpoint for security or compliance reasons.

### Example: Provider block with custom_api_url

```hcl
provider "logzio" {
  api_token      = var.api_token
  custom_api_url = "https://my-internal-gateway.company.com/logzio"
}
```

- If `custom_api_url` is set, it takes precedence over the `region` argument. All API requests will use this URL.
- You can also set it via the `LOGZIO_CUSTOM_API_URL` environment variable:

```bash
export LOGZIO_CUSTOM_API_URL="https://my-internal-gateway.company.com/logzio"
```

Then your provider block can omit the argument:

```hcl
provider "logzio" {}
```

**Tip:**
- Make sure your custom endpoint is reachable from where Terraform runs, and that it properly proxies or implements the Logz.io API.

For more details, see the [main readme](../../readme.md#configuring-the-provider). 