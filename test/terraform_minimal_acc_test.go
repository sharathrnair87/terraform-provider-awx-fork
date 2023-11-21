package test

import (
	"strconv"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestTerraformMinimalAccExample(t *testing.T) {
	terraformOptions := &terraform.Options{
		TerraformDir: "../terraform/",
	}
	terraformDataOptions := &terraform.Options{
		TerraformDir: "../terraform/",
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	output := terraform.Output(t, terraformOptions, "inventory_id")
	nr, err := strconv.Atoi(output)
	if err != nil {
		t.Logf("Inventory id is not a number")
		t.Fail()
	}
	assert.Greater(t, nr, 0)

	terraform.InitAndApply(t, terraformDataOptions)

	defer terraform.Destroy(t, terraformDataOptions)

	output = terraform.Output(t, terraformDataOptions, "job")

}
