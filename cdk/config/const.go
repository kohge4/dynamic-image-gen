package config

var Const *ConstConfig

type ConstConfig struct {
	ApiGatewayName     string
	LamndaFunctionName string
	ResorceDescription string
}

func init() {
	Const = &ConstConfig{
		ApiGatewayName:     "dynamic-image-api",
		LamndaFunctionName: "dynamic-image-lambda",
		ResorceDescription: "Dynamic Image Generator CDK Stack",
	}
}
