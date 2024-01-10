// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sfn_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/YakDriver/regexache"

	"github.com/aws/aws-sdk-go/service/sfn"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"

	// TIP: You will often need to import the package that this test file lives
	// in. Since it is in the "test" context, it must import the package to use
	// any normal context constants, variables, or functions.

	"github.com/hashicorp/terraform-provider-aws/names"
)

// TIP: File Structure. The basic outline for all test files should be as
// follows. Improve this data source's maintainability by following this
// outline.
//
// 1. Package declaration (add "_test" since this is a test file)
// 2. Imports
// 3. Unit tests
// 4. Basic test
// 5. Disappears test
// 6. All the other tests
// 7. Helper functions (exists, destroy, check, etc.)
// 8. Functions that return Terraform configurations

// TIP: ==== ACCEPTANCE TESTS ====
// This is an example of a basic acceptance test. This should test as much of
// standard functionality of the data source as possible, and test importing, if
// applicable. We prefix its name with "TestAcc", the service, and the
// data source name.
//
// Acceptance test access AWS and cost money to run.
func TestAccSFNStateMachinesDataSource_basic(t *testing.T) {
	ctx := acctest.Context(t)
	// TIP: This is a long-running test guard for tests that run longer than
	// 300s (5 min) generally.
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	//var statemachines sfn.ListStateMachinesOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	dataSourceName := "data.aws_sfn_state_machines.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, sfn.EndpointsID)
			//testAccPreCheck(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.SFN),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckStateMachinesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccStateMachinesDataSourceConfig_basic(rName, 1),
				Check: resource.ComposeTestCheckFunc(
					// testAccCheckStateMachinesExists(ctx, dataSourceName, &statemachines),
					//resource.TestCheckResourceAttr(dataSourceName, "auto_minor_version_upgrade", "false"),
					//resource.TestCheckResourceAttr(dataSourceName, "names.0", rName),
					resource.TestCheckResourceAttr(dataSourceName, "arns.#", "1"),
					resource.TestCheckResourceAttrPair(dataSourceName, "arns.0", rName, "arn"),
					resource.TestCheckResourceAttr(dataSourceName, "names.#", "1"),
					resource.TestCheckResourceAttrPair(dataSourceName, "names.0", rName, "name"),
					//
					//resource.TestCheckResourceAttrSet(dataSourceName, "names"),
					//resource.TestCheckResourceAttrSet(dataSourceName, "arns"),
					/*
						resource.TestCheckTypeSetElemNestedAttrs(dataSourceName, "user.*", map[string]string{
							"console_access": "false",
							"groups.#":       "0",
							"username":       "Test",
							"password":       "TestTest1234",
						}),
					*/
					acctest.MatchResourceAttrRegionalARN(dataSourceName, "arn", "sfn", regexache.MustCompile(`statemachines:+.`)),
				),
			},
		},
	})
}

func testAccCheckStateMachinesExists(ctx context.Context, dataSourceName string, statemachines *sfn.ListStateMachinesOutput) error {
	return nil
}
func testAccCheckStateMachinesDestroy( /*ctx context.Context, */ s *terraform.State) error {
	return nil
}

/*
func testAccStateMachinesDataSourceConfig_basic(rName string) string {
	return `
data "aws_sfn_state_machines" "test" {
  #depends_on = [aws_sfn_state_machine.test]
}
`
}
*/

func testAccStateMachinesDataSourceConfig_basic(rName string, rMaxAttempts int) string {
	return fmt.Sprintf(`
resource "aws_iam_role_policy" "for_lambda" {
  name = "%[1]s-lambda"
  role = aws_iam_role.for_lambda.id

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [{
    "Effect": "Allow",
    "Action": [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents"
    ],
    "Resource": "arn:${data.aws_partition.current.partition}:logs:*:*:*"
  }]
}
EOF
}

resource "aws_iam_role" "for_lambda" {
  name = "%[1]s-lambda"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [{
    "Action": "sts:AssumeRole",
    "Principal": {
      "Service": "lambda.amazonaws.com"
    },
    "Effect": "Allow"
  }]
}
EOF
}

resource "aws_lambda_function" "test" {
  filename      = "test-fixtures/lambdatest.zip"
  function_name = %[1]q
  role          = aws_iam_role.for_lambda.arn
  handler       = "exports.example"
  runtime       = "nodejs16.x"
}

data "aws_region" "current" {}

data "aws_partition" "current" {}

resource "aws_iam_role_policy" "for_sfn" {
  name = "%[1]s-sfn"
  role = aws_iam_role.for_sfn.id

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [{
    "Effect": "Allow",
    "Action": [
      "lambda:InvokeFunction",
      "logs:CreateLogDelivery",
      "logs:GetLogDelivery",
      "logs:UpdateLogDelivery",
      "logs:DeleteLogDelivery",
      "logs:ListLogDeliveries",
      "logs:PutResourcePolicy",
      "logs:DescribeResourcePolicies",
      "logs:DescribeLogGroups",
      "xray:PutTraceSegments",
      "xray:PutTelemetryRecords",
      "xray:GetSamplingRules",
      "xray:GetSamplingTargets"
    ],
    "Resource": "*"
  }]
}
EOF
}

resource "aws_iam_role" "for_sfn" {
  name = "%[1]s-sfn"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [{
    "Effect": "Allow",
    "Principal": {
      "Service": "states.${data.aws_region.current.name}.amazonaws.com"
    },
    "Action": "sts:AssumeRole"
  }]
}
EOF
}

resource "aws_sfn_state_machine" "test" {
  name     = %[1]q
  role_arn = aws_iam_role.for_sfn.arn

  definition = <<EOF
{
  "Comment": "A Hello World example of the Amazon States Language using an AWS Lambda Function",
  "StartAt": "HelloWorld",
  "States": {
    "HelloWorld": {
      "Type": "Task",
      "Resource": "${aws_lambda_function.test.arn}",
      "Retry": [
        {
          "ErrorEquals": [
            "States.ALL"
          ],
          "IntervalSeconds": 5,
          "MaxAttempts": %[2]d,
          "BackoffRate": 8
        }
      ],
      "End": true
    }
  }
}
EOF
}

data "aws_sfn_state_machines" "test" {
  depends_on = [aws_sfn_state_machine.test]
}
`, rName, rMaxAttempts)
}
