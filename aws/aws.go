package aws

import (
	"github.com/aws/aws-sdk-go-v2/aws"
)

const (
	WaitTimeSecs = 10
)

var (
	AwsTypeString = aws.String("String")
	AwsTypeNumber = aws.String("Number")
	AwsTypeBinary = aws.String("Binary")
	AwsEmailValue = aws.String("email")
)
