// test/security_groups_test.go

package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestSecurityGroupModule(t *testing.T) {
	t.Parallel()

	awsRegion := "us-east-2"

	terraformOptions := &terraform.Options{
		TerraformDir: "../modules/Security_groups",
		Vars: map[string]interface{}{
			"vpc_id":                    "vpc-0b740d622a5170532",
			"cluster_security_group_id": "sg-0c81a9e03de46c1ae",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	albSGID := terraform.Output(t, terraformOptions, "alb_sg_id")

	assert.NotEmpty(t, albSGID)

	securityGroup := aws.GetSecurityGroupById(t, albSGID, awsRegion)
	assert.Equal(t, "alb-sg", securityGroup.Name)

	ingressRules := aws.GetSecurityGroupIngressRules(t, albSGID, awsRegion)
	assert.Len(t, ingressRules, 1)
	assert.Equal(t, int64(80), *ingressRules[0].FromPort)
	assert.Equal(t, int64(80), *ingressRules[0].ToPort)
}
