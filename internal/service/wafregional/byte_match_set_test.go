package wafregional_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/waf"
	"github.com/aws/aws-sdk-go/service/wafregional"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfwafregional "github.com/hashicorp/terraform-provider-aws/internal/service/wafregional"
)

func TestAccWAFRegionalByteMatchSet_basic(t *testing.T) {
	ctx := acctest.Context(t)
	var v waf.ByteMatchSet
	byteMatchSet := fmt.Sprintf("byteMatchSet-%s", sdkacctest.RandString(5))
	resourceName := "aws_wafregional_byte_match_set.byte_set"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); acctest.PreCheckPartitionHasService(t, wafregional.EndpointsID) },
		ErrorCheck:               acctest.ErrorCheck(t, wafregional.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckByteMatchSetDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccByteMatchSetConfig_basic(byteMatchSet),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckByteMatchSetExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(
						resourceName, "name", byteMatchSet),
					resource.TestCheckResourceAttr(
						resourceName, "byte_match_tuples.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "byte_match_tuples.*", map[string]string{
						"field_to_match.#":      "1",
						"field_to_match.0.data": "referer",
						"field_to_match.0.type": "HEADER",
						"positional_constraint": "CONTAINS",
						"target_string":         "badrefer1",
						"text_transformation":   "NONE",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "byte_match_tuples.*", map[string]string{
						"field_to_match.#":      "1",
						"field_to_match.0.data": "referer",
						"field_to_match.0.type": "HEADER",
						"positional_constraint": "CONTAINS",
						"target_string":         "badrefer2",
						"text_transformation":   "NONE",
					}),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccWAFRegionalByteMatchSet_changeNameForceNew(t *testing.T) {
	ctx := acctest.Context(t)
	var before, after waf.ByteMatchSet
	byteMatchSet := fmt.Sprintf("byteMatchSet-%s", sdkacctest.RandString(5))
	byteMatchSetNewName := fmt.Sprintf("byteMatchSet-%s", sdkacctest.RandString(5))
	resourceName := "aws_wafregional_byte_match_set.byte_set"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); acctest.PreCheckPartitionHasService(t, wafregional.EndpointsID) },
		ErrorCheck:               acctest.ErrorCheck(t, wafregional.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckByteMatchSetDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccByteMatchSetConfig_basic(byteMatchSet),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckByteMatchSetExists(ctx, resourceName, &before),
					resource.TestCheckResourceAttr(
						resourceName, "name", byteMatchSet),
					resource.TestCheckResourceAttr(
						resourceName, "byte_match_tuples.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "byte_match_tuples.*", map[string]string{
						"field_to_match.#":      "1",
						"field_to_match.0.data": "referer",
						"field_to_match.0.type": "HEADER",
						"positional_constraint": "CONTAINS",
						"target_string":         "badrefer1",
						"text_transformation":   "NONE",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "byte_match_tuples.*", map[string]string{
						"field_to_match.#":      "1",
						"field_to_match.0.data": "referer",
						"field_to_match.0.type": "HEADER",
						"positional_constraint": "CONTAINS",
						"target_string":         "badrefer2",
						"text_transformation":   "NONE",
					}),
				),
			},
			{
				Config: testAccByteMatchSetConfig_changeName(byteMatchSetNewName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckByteMatchSetExists(ctx, resourceName, &after),
					resource.TestCheckResourceAttr(
						resourceName, "name", byteMatchSetNewName),
					resource.TestCheckResourceAttr(
						resourceName, "byte_match_tuples.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "byte_match_tuples.*", map[string]string{
						"field_to_match.#":      "1",
						"field_to_match.0.data": "referer",
						"field_to_match.0.type": "HEADER",
						"positional_constraint": "CONTAINS",
						"target_string":         "badrefer1",
						"text_transformation":   "NONE",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "byte_match_tuples.*", map[string]string{
						"field_to_match.#":      "1",
						"field_to_match.0.data": "referer",
						"field_to_match.0.type": "HEADER",
						"positional_constraint": "CONTAINS",
						"target_string":         "badrefer2",
						"text_transformation":   "NONE",
					}),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccWAFRegionalByteMatchSet_changeByteMatchTuples(t *testing.T) {
	ctx := acctest.Context(t)
	var before, after waf.ByteMatchSet
	byteMatchSetName := fmt.Sprintf("byte-batch-set-%s", sdkacctest.RandString(5))
	resourceName := "aws_wafregional_byte_match_set.byte_set"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); acctest.PreCheckPartitionHasService(t, wafregional.EndpointsID) },
		ErrorCheck:               acctest.ErrorCheck(t, wafregional.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckByteMatchSetDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccByteMatchSetConfig_basic(byteMatchSetName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckByteMatchSetExists(ctx, resourceName, &before),
					resource.TestCheckResourceAttr(
						resourceName, "name", byteMatchSetName),
					resource.TestCheckResourceAttr(
						resourceName, "byte_match_tuples.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "byte_match_tuples.*", map[string]string{
						"field_to_match.0.data": "referer",
						"field_to_match.0.type": "HEADER",
						"positional_constraint": "CONTAINS",
						"target_string":         "badrefer1",
						"text_transformation":   "NONE",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "byte_match_tuples.*", map[string]string{
						"field_to_match.0.data": "referer",
						"field_to_match.0.type": "HEADER",
						"positional_constraint": "CONTAINS",
						"target_string":         "badrefer2",
						"text_transformation":   "NONE",
					}),
				),
			},
			{
				Config: testAccByteMatchSetConfig_changeTuples(byteMatchSetName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckByteMatchSetExists(ctx, resourceName, &after),
					resource.TestCheckResourceAttr(
						resourceName, "name", byteMatchSetName),
					resource.TestCheckResourceAttr(
						resourceName, "byte_match_tuples.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "byte_match_tuples.*", map[string]string{
						"field_to_match.#":      "1",
						"field_to_match.0.data": "",
						"field_to_match.0.type": "METHOD",
						"positional_constraint": "EXACTLY",
						"target_string":         "GET",
						"text_transformation":   "NONE",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "byte_match_tuples.*", map[string]string{
						"field_to_match.#":      "1",
						"field_to_match.0.data": "",
						"field_to_match.0.type": "URI",
						"positional_constraint": "ENDS_WITH",
						"target_string":         "badrefer2+",
						"text_transformation":   "LOWERCASE",
					}),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccWAFRegionalByteMatchSet_noByteMatchTuples(t *testing.T) {
	ctx := acctest.Context(t)
	var byteMatchSet waf.ByteMatchSet
	byteMatchSetName := fmt.Sprintf("byte-batch-set-%s", sdkacctest.RandString(5))
	resourceName := "aws_wafregional_byte_match_set.byte_match_set"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); acctest.PreCheckPartitionHasService(t, wafregional.EndpointsID) },
		ErrorCheck:               acctest.ErrorCheck(t, wafregional.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckByteMatchSetDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccByteMatchSetConfig_noDescriptors(byteMatchSetName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckByteMatchSetExists(ctx, resourceName, &byteMatchSet),
					resource.TestCheckResourceAttr(resourceName, "name", byteMatchSetName),
					resource.TestCheckResourceAttr(resourceName, "byte_match_tuples.#", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccWAFRegionalByteMatchSet_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	var v waf.ByteMatchSet
	byteMatchSet := fmt.Sprintf("byteMatchSet-%s", sdkacctest.RandString(5))
	resourceName := "aws_wafregional_byte_match_set.byte_set"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); acctest.PreCheckPartitionHasService(t, wafregional.EndpointsID) },
		ErrorCheck:               acctest.ErrorCheck(t, wafregional.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckByteMatchSetDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccByteMatchSetConfig_basic(byteMatchSet),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckByteMatchSetExists(ctx, resourceName, &v),
					testAccCheckByteMatchSetDisappears(ctx, &v),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckByteMatchSetDisappears(ctx context.Context, v *waf.ByteMatchSet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).WAFRegionalConn()
		region := acctest.Provider.Meta().(*conns.AWSClient).Region

		wr := tfwafregional.NewRetryer(conn, region)
		_, err := wr.RetryWithToken(ctx, func(token *string) (interface{}, error) {
			req := &waf.UpdateByteMatchSetInput{
				ChangeToken:    token,
				ByteMatchSetId: v.ByteMatchSetId,
			}

			for _, ByteMatchTuple := range v.ByteMatchTuples {
				ByteMatchUpdate := &waf.ByteMatchSetUpdate{
					Action: aws.String("DELETE"),
					ByteMatchTuple: &waf.ByteMatchTuple{
						FieldToMatch:         ByteMatchTuple.FieldToMatch,
						PositionalConstraint: ByteMatchTuple.PositionalConstraint,
						TargetString:         ByteMatchTuple.TargetString,
						TextTransformation:   ByteMatchTuple.TextTransformation,
					},
				}
				req.Updates = append(req.Updates, ByteMatchUpdate)
			}

			return conn.UpdateByteMatchSetWithContext(ctx, req)
		})
		if err != nil {
			return fmt.Errorf("Error updating ByteMatchSet: %s", err)
		}

		_, err = wr.RetryWithToken(ctx, func(token *string) (interface{}, error) {
			opts := &waf.DeleteByteMatchSetInput{
				ChangeToken:    token,
				ByteMatchSetId: v.ByteMatchSetId,
			}
			return conn.DeleteByteMatchSetWithContext(ctx, opts)
		})
		if err != nil {
			return fmt.Errorf("Error deleting ByteMatchSet: %s", err)
		}

		return nil
	}
}

func testAccCheckByteMatchSetExists(ctx context.Context, n string, v *waf.ByteMatchSet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No WAF ByteMatchSet ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).WAFRegionalConn()
		resp, err := conn.GetByteMatchSetWithContext(ctx, &waf.GetByteMatchSetInput{
			ByteMatchSetId: aws.String(rs.Primary.ID),
		})

		if err != nil {
			return err
		}

		if *resp.ByteMatchSet.ByteMatchSetId == rs.Primary.ID {
			*v = *resp.ByteMatchSet
			return nil
		}

		return fmt.Errorf("WAF ByteMatchSet (%s) not found", rs.Primary.ID)
	}
}

func testAccCheckByteMatchSetDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_wafregional_byte_match_set" {
				continue
			}

			conn := acctest.Provider.Meta().(*conns.AWSClient).WAFRegionalConn()
			resp, err := conn.GetByteMatchSetWithContext(ctx, &waf.GetByteMatchSetInput{
				ByteMatchSetId: aws.String(rs.Primary.ID),
			})

			if err == nil {
				if *resp.ByteMatchSet.ByteMatchSetId == rs.Primary.ID {
					return fmt.Errorf("WAF ByteMatchSet %s still exists", rs.Primary.ID)
				}
			}

			if tfawserr.ErrCodeEquals(err, waf.ErrCodeNonexistentItemException) {
				continue
			}

			return err
		}

		return nil
	}
}

func testAccByteMatchSetConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "aws_wafregional_byte_match_set" "byte_set" {
  name = "%s"

  byte_match_tuples {
    text_transformation   = "NONE"
    target_string         = "badrefer1"
    positional_constraint = "CONTAINS"

    field_to_match {
      type = "HEADER"
      data = "referer"
    }
  }

  byte_match_tuples {
    text_transformation   = "NONE"
    target_string         = "badrefer2"
    positional_constraint = "CONTAINS"

    field_to_match {
      type = "HEADER"
      data = "referer"
    }
  }
}
`, name)
}

func testAccByteMatchSetConfig_changeName(name string) string {
	return fmt.Sprintf(`
resource "aws_wafregional_byte_match_set" "byte_set" {
  name = "%s"

  byte_match_tuples {
    text_transformation   = "NONE"
    target_string         = "badrefer1"
    positional_constraint = "CONTAINS"

    field_to_match {
      type = "HEADER"
      data = "referer"
    }
  }

  byte_match_tuples {
    text_transformation   = "NONE"
    target_string         = "badrefer2"
    positional_constraint = "CONTAINS"

    field_to_match {
      type = "HEADER"
      data = "referer"
    }
  }
}
`, name)
}

func testAccByteMatchSetConfig_noDescriptors(name string) string {
	return fmt.Sprintf(`
resource "aws_wafregional_byte_match_set" "byte_match_set" {
  name = "%s"
}
`, name)
}

func testAccByteMatchSetConfig_changeTuples(name string) string {
	return fmt.Sprintf(`
resource "aws_wafregional_byte_match_set" "byte_set" {
  name = "%s"

  byte_match_tuples {
    text_transformation   = "NONE"
    target_string         = "GET"
    positional_constraint = "EXACTLY"

    field_to_match {
      type = "METHOD"
    }
  }

  byte_match_tuples {
    text_transformation   = "LOWERCASE"
    target_string         = "badrefer2+"
    positional_constraint = "ENDS_WITH"

    field_to_match {
      type = "URI"
    }
  }
}
`, name)
}
