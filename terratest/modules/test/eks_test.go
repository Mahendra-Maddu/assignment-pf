// test/eks_test.go

package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestEKSModule(t *testing.T) {
	t.Parallel()

	awsRegion := "us-east-2"

	terraformOptions := &terraform.Options{
		TerraformDir: "../modules/eks",
		Vars: map[string]interface{}{
			"vpc_id":             "vpc-0b740d622a5170532",
			"private_subnet_ids": []string{"subnet-00ed61b50207d6724", "subnet-0bfe9e0d4d7af6bb0"},
			"cluster_name":       "test-eks-cluster",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	clusterID := terraform.Output(t, terraformOptions, "cluster_id")
	clusterEndpoint := terraform.Output(t, terraformOptions, "cluster_endpoint")
	clusterSGID := terraform.Output(t, terraformOptions, "cluster_security_group_id")

	assert.NotEmpty(t, clusterID)
	assert.NotEmpty(t, clusterEndpoint)
	assert.NotEmpty(t, clusterSGID)

	cluster, err := aws.GetEksClusterByNameE(t, awsRegion, "main-eks-cluster")
	assert.NoError(t, err)
	assert.Equal(t, "ACTIVE", *cluster.Status)
}
