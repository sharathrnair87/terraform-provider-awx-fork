// Intialize Provider with username password
provider "awx" {
  // hostname can be set in the variable AWX_HOSTNAME
  hostname = "https://awx.example.com"
  // username can be set in the variable AWX_USERNAME
  username = "awxuser"
  // password can be set in the variable AWX_PASSWORD
  password = "awxpassword"
}

// initialize using Token
provider "awx" {
  hostname = "https://awx.example.com"
  // token can be set in the variable AWX_TOKEN
  token = "awxoauth2token"
}
