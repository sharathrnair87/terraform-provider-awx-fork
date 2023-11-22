
data "awx_credential_type" "my_cust_cred_type" {
  id = var.my_cust_cred_type_id
}

output "my_cust_cred_type_inputs" {
  value = data.awx_credential_type.my_cust_cred_type.inputs
}

