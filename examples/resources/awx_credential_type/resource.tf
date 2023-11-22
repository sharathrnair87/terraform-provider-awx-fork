
data "awx_organization" "default" {
  name = "Default"
}

resource "awx_credential_type" "custom_cred_type" {
  name = "customcreds"
  inputs = jsonencode(
    {
      fields = [
        {
          id     = "url",
          label  = "URL",
          secret = false,
          type   = "string",
        },
        {
          id     = "url_token",
          label  = "URL TOKEN",
          secret = true,
          type   = "string",
        }
      ]
    }
  )
  injectors = jsonencode(
    {
      "env" = {
        url       = "{{url}}",
        url_token = "{{url_token}}",
      }
    }
  )
}

