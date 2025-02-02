package main

import (
	"fmt"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigatewayv2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigatewayv2integrations"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscloudfront"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscloudfrontorigins"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsecr"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/constructs-go/constructs/v10"

	"github.com/aws/jsii-runtime-go"

	"cdk/config"
)

type CdkStackProps struct {
	awscdk.StackProps
}

func NewImageGenAppCdkStack(scope constructs.Construct, id string, props *CdkStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	// ECRリポジトリを参照
	ecrRepo := awsecr.Repository_FromRepositoryName(stack, jsii.String(config.Env.ECRRegistryName), jsii.String(config.Env.ECRRepositoryName))

	// Lambda関数をECRイメージから作成
	lambdaFunction := NewLambdaFunctionByECR(stack, ecrRepo)

	// API Gateway HTTP APIを作成してLambdaを統合
	api := NewApiGatewayAndLambdaIntegration(stack, lambdaFunction)

	// CloudFrontのキャッシュポリシーを作成
	cachePolicy := NewCachePolicy(stack)

	// CloudFrontのログ用S3バケットを作成
	logBucket := awss3.NewBucket(stack, jsii.String("LogBucket"), &awss3.BucketProps{
		ObjectOwnership: awss3.ObjectOwnership_BUCKET_OWNER_PREFERRED,
	})

	// CloudFrontディストリビューションを作成
	distribution := NewCloudFrontDistributionByApiGateway(stack, api, cachePolicy, logBucket)

	awscdk.NewCfnOutput(stack, jsii.String("ApiUrl"), &awscdk.CfnOutputProps{
		Value: api.ApiEndpoint(),
	})
	awscdk.NewCfnOutput(stack, jsii.String("CloudFrontUrl"), &awscdk.CfnOutputProps{
		Value: distribution.DomainName(),
	})

	return stack
}

func NewLambdaFunctionByECR(stack awscdk.Stack, ecrRepo awsecr.IRepository) awslambda.IFunction {
	containerFunction := awslambda.NewDockerImageFunction(stack, jsii.String("ECRLambdaFunction"), &awslambda.DockerImageFunctionProps{
		FunctionName: jsii.String(fmt.Sprintf("%s_%s", config.Const.LamndaFunctionName, *stack.StackName())),
		Description:  jsii.String(config.Const.ResorceDescription),
		Code: awslambda.DockerImageCode_FromEcr(ecrRepo, &awslambda.EcrImageCodeProps{
			Tag: jsii.String("latest"),
		}),
		MemorySize:           jsii.Number(1024),                        // メモリサイズを1024MBに設定
		EphemeralStorageSize: awscdk.Size_Mebibytes(jsii.Number(1024)), // 一時ストレージサイズを1024MBに設定
		Timeout:              awscdk.Duration_Seconds(jsii.Number(30)), // タイムアウトを30秒に設定
	})
	return containerFunction
}

// API Gateway HTTP APIを作成してLambdaを統合
func NewApiGatewayAndLambdaIntegration(stack awscdk.Stack, lambdaFunction awslambda.IFunction) awsapigatewayv2.HttpApi {
	api := awsapigatewayv2.NewHttpApi(stack, jsii.String("HttpApi"), &awsapigatewayv2.HttpApiProps{
		ApiName:     jsii.String(fmt.Sprintf("%s_%s", config.Const.ApiGatewayName, *stack.StackName())),
		Description: jsii.String(config.Const.ResorceDescription),
	})
	integration := awsapigatewayv2integrations.NewHttpLambdaIntegration(
		jsii.String("LambdaIntegration"),
		lambdaFunction,
		&awsapigatewayv2integrations.HttpLambdaIntegrationProps{},
	)
	api.AddRoutes(&awsapigatewayv2.AddRoutesOptions{
		Path:        jsii.String("/{proxy+}"),
		Integration: integration,
	})
	return api
}

func NewCachePolicy(stack awscdk.Stack) awscloudfront.ICachePolicy {
	cachePolicy := awscloudfront.NewCachePolicy(stack, jsii.String("CustomCachePolicy"), &awscloudfront.CachePolicyProps{
		CachePolicyName:     jsii.String(fmt.Sprintf("%s_%s", "DynamicImageCachePolicy", *stack.StackName())),
		Comment:             jsii.String(config.Const.ResorceDescription),
		DefaultTtl:          awscdk.Duration_Seconds(jsii.Number(60)), // デフォルトTTLを60秒に設定
		MinTtl:              awscdk.Duration_Seconds(jsii.Number(10)),
		MaxTtl:              awscdk.Duration_Seconds(jsii.Number(300)),
		QueryStringBehavior: awscloudfront.CacheQueryStringBehavior_All(), // 全てのクエリパラメータをキャッシュキーに含める
	})
	return cachePolicy
}

func NewCloudFrontDistributionByApiGateway(stack awscdk.Stack, api awsapigatewayv2.HttpApi, cachePolicy awscloudfront.ICachePolicy, logBucket awss3.IBucket) awscloudfront.Distribution {
	distribution := awscloudfront.NewDistribution(stack, jsii.String("CloudFrontDistribution"), &awscloudfront.DistributionProps{
		DefaultBehavior: &awscloudfront.BehaviorOptions{
			Origin: awscloudfrontorigins.NewHttpOrigin(
				jsii.String("dynamic-image-dummy-origin.com"), // APIGatewayのdomain名を後から設定する必要があるためダミーの値を設定
				&awscloudfrontorigins.HttpOriginProps{
					ProtocolPolicy: awscloudfront.OriginProtocolPolicy_HTTPS_ONLY,
				}),
			CachePolicy:          cachePolicy,
			ViewerProtocolPolicy: awscloudfront.ViewerProtocolPolicy_REDIRECT_TO_HTTPS,
		},
		LogBucket: logBucket,
		Comment:   jsii.String(config.Const.ResorceDescription),
	})

	// OriginをAPI Gatewayに変更
	// CloudFrontリソースを直接取得
	cfnDistribution := distribution.Node().DefaultChild().(awscloudfront.CfnDistribution)
	// CloudFormationテンプレートでオリジンのDomainNameを動的に設定
	cfnDistribution.AddOverride(jsii.String("Properties.DistributionConfig.Origins.0.DomainName"), map[string]interface{}{
		"Fn::Select": []interface{}{
			2, // URLの3番目の部分（スラッシュで分割したホスト名部分）
			map[string]interface{}{
				"Fn::Split": []interface{}{
					"/",
					api.ApiEndpoint(),
				},
			},
		},
	})
	return distribution
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewImageGenAppCdkStack(app, config.Env.CDKStackID, &CdkStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	// If unspecified, this stack will be "environment-agnostic".
	// Account/Region-dependent features and context lookups will not work, but a
	// single synthesized template can be deployed anywhere.
	//---------------------------------------------------------------------------
	return nil

	// Uncomment if you know exactly what account and region you want to deploy
	// the stack to. This is the recommendation for production stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String("123456789012"),
	//  Region:  jsii.String("us-east-1"),
	// }

	// Uncomment to specialize this stack for the AWS Account and Region that are
	// implied by the current CLI configuration. This is recommended for dev
	// stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
	//  Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	// }
}
