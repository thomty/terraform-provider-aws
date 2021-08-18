package aws

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/appstream"
	"github.com/hashicorp/aws-sdk-go-base/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAwsAppStreamImageBuilder_basic(t *testing.T) {
	var imageBuilderOutput appstream.ImageBuilder
	resourceName := "aws_appstream_image_builder.test"
	instanceType := "stream.standard.small"
	rName := acctest.RandomWithPrefix("tf-acc-test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckAwsAppStreamImageBuilderDestroy,
		ErrorCheck:        testAccErrorCheck(t, appstream.EndpointsID),
		Steps: []resource.TestStep{
			{
				Config: testAccAwsAppStreamImageBuilderConfig(instanceType, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsAppStreamImageBuilderExists(resourceName, &imageBuilderOutput),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					testAccCheckResourceAttrRfc3339(resourceName, "created_time"),
					resource.TestCheckResourceAttr(resourceName, "state", appstream.ImageBuilderStateRunning),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"image_name"},
			},
		},
	})
}

func TestAccAwsAppStreamImageBuilder_disappears(t *testing.T) {
	var imageBuilderOutput appstream.ImageBuilder
	resourceName := "aws_appstream_image_builder.test"
	instanceType := "stream.standard.medium"
	rName := acctest.RandomWithPrefix("tf-acc-test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckAwsAppStreamImageBuilderDestroy,
		ErrorCheck:        testAccErrorCheck(t, appstream.EndpointsID),
		Steps: []resource.TestStep{
			{
				Config: testAccAwsAppStreamImageBuilderConfig(instanceType, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsAppStreamImageBuilderExists(resourceName, &imageBuilderOutput),
					testAccCheckResourceDisappears(testAccProvider, resourceAwsAppStreamImageBuilder(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAwsAppStreamImageBuilder_complete(t *testing.T) {
	var imageBuilderOutput appstream.ImageBuilder
	resourceName := "aws_appstream_image_builder.test"
	rName := acctest.RandomWithPrefix("tf-acc-test")
	description := "Description of a test"
	descriptionUpdated := "Updated Description of a test"
	instanceType := "stream.standard.small"
	instanceTypeUpdate := "stream.standard.medium"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckAwsAppStreamImageBuilderDestroy,
		ErrorCheck:        testAccErrorCheck(t, appstream.EndpointsID),
		Steps: []resource.TestStep{
			{
				Config: testAccAwsAppStreamImageBuilderConfigComplete(rName, description, instanceType),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsAppStreamImageBuilderExists(resourceName, &imageBuilderOutput),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "state", appstream.ImageBuilderStateRunning),
					resource.TestCheckResourceAttr(resourceName, "instance_type", instanceType),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					testAccCheckResourceAttrRfc3339(resourceName, "created_time"),
				),
			},
			{
				Config: testAccAwsAppStreamImageBuilderConfigComplete(rName, descriptionUpdated, instanceTypeUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsAppStreamImageBuilderExists(resourceName, &imageBuilderOutput),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "state", appstream.ImageBuilderStateRunning),
					resource.TestCheckResourceAttr(resourceName, "instance_type", instanceTypeUpdate),
					resource.TestCheckResourceAttr(resourceName, "description", descriptionUpdated),
					testAccCheckResourceAttrRfc3339(resourceName, "created_time"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"image_name"},
			},
		},
	})
}

func TestAccAwsAppStreamImageBuilder_withTags(t *testing.T) {
	var imageBuilderOutput appstream.ImageBuilder
	resourceName := "aws_appstream_image_builder.test"
	rName := acctest.RandomWithPrefix("tf-acc-test")
	description := "Description of a test"
	descriptionUpdated := "Updated Description of a test"
	instanceType := "stream.standard.small"
	instanceTypeUpdate := "stream.standard.medium"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckAwsAppStreamImageBuilderDestroy,
		ErrorCheck:        testAccErrorCheck(t, appstream.EndpointsID),
		Steps: []resource.TestStep{
			{
				Config: testAccAwsAppStreamImageBuilderConfigWithTags(rName, description, instanceType),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsAppStreamImageBuilderExists(resourceName, &imageBuilderOutput),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags_all.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags_all.Key", "value"),
					resource.TestCheckResourceAttr(resourceName, "state", appstream.ImageBuilderStateRunning),
					resource.TestCheckResourceAttr(resourceName, "instance_type", instanceType),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					testAccCheckResourceAttrRfc3339(resourceName, "created_time"),
				),
			},
			{
				Config: testAccAwsAppStreamImageBuilderConfigWithTags(rName, descriptionUpdated, instanceTypeUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsAppStreamImageBuilderExists(resourceName, &imageBuilderOutput),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.Key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags_all.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags_all.Key", "value"),
					resource.TestCheckResourceAttr(resourceName, "state", appstream.ImageBuilderStateRunning),
					resource.TestCheckResourceAttr(resourceName, "instance_type", instanceTypeUpdate),
					resource.TestCheckResourceAttr(resourceName, "description", descriptionUpdated),
					testAccCheckResourceAttrRfc3339(resourceName, "created_time"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"image_name"},
			},
		},
	})
}

func testAccCheckAwsAppStreamImageBuilderExists(resourceName string, appStreamImageBuilder *appstream.ImageBuilder) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		conn := testAccProvider.Meta().(*AWSClient).appstreamconn
		resp, err := conn.DescribeImageBuilders(&appstream.DescribeImageBuildersInput{Names: []*string{aws.String(rs.Primary.ID)}})

		if err != nil {
			return err
		}

		if resp == nil && len(resp.ImageBuilders) == 0 {
			return fmt.Errorf("appstream imageBuilder %q does not exist", rs.Primary.ID)
		}

		*appStreamImageBuilder = *resp.ImageBuilders[0]

		return nil
	}
}

func testAccCheckAwsAppStreamImageBuilderDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*AWSClient).appstreamconn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_appstream_image_builder" {
			continue
		}

		resp, err := conn.DescribeImageBuilders(&appstream.DescribeImageBuildersInput{Names: []*string{aws.String(rs.Primary.ID)}})

		if tfawserr.ErrCodeEquals(err, appstream.ErrCodeResourceNotFoundException) {
			continue
		}

		if err != nil {
			return err
		}

		if resp != nil && len(resp.ImageBuilders) > 0 {
			return fmt.Errorf("appstream imageBuilder %q still exists", rs.Primary.ID)
		}
	}

	return nil

}

func testAccAwsAppStreamImageBuilderConfig(instanceType, name string) string {
	return fmt.Sprintf(`
resource "aws_appstream_image_builder" "test" {
  image_name    = "Amazon-AppStream2-Sample-Image-02-04-2019"
  instance_type = %[1]q
  name          = %[2]q
}
`, instanceType, name)
}

func testAccAwsAppStreamImageBuilderConfigComplete(name, description, instanceType string) string {
	return fmt.Sprintf(`
data "aws_availability_zones" "available" {
  state = "available"
}

resource "aws_vpc" "example" {
  cidr_block = "192.168.0.0/16"
}


resource "aws_subnet" "example" {
  availability_zone = data.aws_availability_zones.available.names[0]
  cidr_block        = "192.168.0.0/24"
  vpc_id            = aws_vpc.example.id
}

resource "aws_appstream_image_builder" "test" {
  image_name                     = "Amazon-AppStream2-Sample-Image-02-04-2019"
  name                           = %[1]q
  description                    = %[2]q
  enable_default_internet_access = false
  instance_type                  = %[3]q
  vpc_config {
    subnet_ids = [aws_subnet.example.id]
  }
}
`, name, description, instanceType)
}

func testAccAwsAppStreamImageBuilderConfigWithTags(name, description, instanceType string) string {
	return fmt.Sprintf(`
data "aws_availability_zones" "available" {
  state = "available"
}

resource "aws_vpc" "example" {
  cidr_block = "192.168.0.0/16"
}

resource "aws_subnet" "example" {
  availability_zone = data.aws_availability_zones.available.names[0]
  cidr_block        = "192.168.0.0/24"
  vpc_id            = aws_vpc.example.id
}

resource "aws_appstream_image_builder" "test" {
  image_name                     = "Amazon-AppStream2-Sample-Image-02-04-2019"
  name                           = %[1]q
  description                    = %[2]q
  enable_default_internet_access = false
  instance_type                  = %[3]q
  vpc_config {
    subnet_ids = [aws_subnet.example.id]
  }
  tags = {
    Key = "value"
  }
}
`, name, description, instanceType)
}
