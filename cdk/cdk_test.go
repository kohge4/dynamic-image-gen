package main

import (
	"encoding/json"
	"log"
	"os"
	"testing"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/assertions"
	"github.com/aws/jsii-runtime-go"

	"github.com/sebdah/goldie/v2"
)

// 参考 リソース名 https://docs.aws.amazon.com/ja_jp/AWSCloudFormation/latest/UserGuide/aws-template-resource-type-ref.html
func TestCdkStackHasResource(t *testing.T) {
	app := awscdk.NewApp(nil)
	stack := NewImageGenAppCdkStack(app, "TestCdkStack", nil)

	template := assertions.Template_FromStack(stack, nil)

	template.ResourceCountIs(jsii.String("AWS::Lambda::Function"), jsii.Number(1))
	template.HasResourceProperties(jsii.String("AWS::Lambda::Function"), map[string]interface{}{
		"PackageType": "Image",
		"Code": map[string]interface{}{
			"ImageUri": assertions.Match_AnyValue(),
		},
	})

	template.ResourceCountIs(jsii.String("AWS::ApiGatewayV2::Api"), jsii.Number(1))
	template.HasResourceProperties(jsii.String("AWS::ApiGatewayV2::Api"), map[string]interface{}{
		"ProtocolType": "HTTP",
	})
	template.HasResourceProperties(jsii.String("AWS::ApiGatewayV2::Integration"), map[string]interface{}{
		"IntegrationType": "AWS_PROXY",
		"IntegrationUri":  assertions.Match_AnyValue(),
	})
	template.HasResourceProperties(jsii.String("AWS::ApiGatewayV2::Route"), map[string]interface{}{
		"RouteKey": "ANY /{proxy+}",
		"Target":   assertions.Match_AnyValue(),
	})

	template.ResourceCountIs(jsii.String("AWS::CloudFront::CachePolicy"), jsii.Number(1))

	template.ResourceCountIs(jsii.String("AWS::CloudFront::Distribution"), jsii.Number(1))
}

func TestCdkStackSnapshot(t *testing.T) {
	if os.Getenv("CI") == "true" {
		t.Skip("Skip snapshot test in CI environment")
	}

	app := awscdk.NewApp(nil)
	stack := NewImageGenAppCdkStack(app, "TestCdkStack", nil)

	template := assertions.Template_FromStack(stack, nil)
	templateJSON, err := json.Marshal(template.ToJSON())
	if err != nil {
		log.Fatalf("Error marshaling map: %v", err)
	}

	data, err := os.ReadFile("./testdata/cdkStackTest.golden")
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	g := goldie.New(t)
	g.AssertWithTemplate(t, "cdkStackTest", data, templateJSON)
}
