package main

// Available regions on AWS.
const (
	ApNorthEast1 = "ap-northeast-1"
	ApNorthEast2 = "ap-northeast-2"
	ApSouth1     = "ap-south-1"
	ApSouthEast1 = "ap-southeast-1"
	ApSouthEast2 = "ap-southeast-2"
	CaCentral1   = "ca-central-1"
	EuCentral1   = "eu-central-1"
	EuWest1      = "eu-west-1"
	EuWest2      = "eu-west-2"
	EuWest3      = "eu-west-3"
	SaEast1      = "sa-east-1"
	UsEast1      = "us-east-1"
	UsEast2      = "us-east-2"
	UsWest1      = "us-west-1"
	UsWest2      = "us-west-2"
)

var AllAwsRegions = []string{
	ApNorthEast1,
	ApNorthEast2,
	ApSouth1,
	ApSouthEast1,
	ApSouthEast2,
	CaCentral1,
	EuCentral1,
	EuWest1,
	EuWest2,
	EuWest3,
	SaEast1,
	UsEast1,
	UsEast2,
	UsWest1,
	UsWest2,
}

// Available instance types on AWS.
const (
	T2Nano    = "t2.nano"
	T2Micro   = "t2.micro"
	T2Small   = "t2.small"
	T2Medium  = "t2.medium"
	T2Large   = "t2.large"
	T2XLarge  = "t2.xlarge"
	T2XXLarge = "t2.xxlarge"

	M4Large    = "m4.large"
	M4XLarge   = "m4.xlarge"
	M42XLarge  = "m4.2xlarge"
	M44XLarge  = "m4.4xlarge"
	M410XLarge = "m4.10xlarge"
	M416XLarge = "m4.16xlarge"

	M5Large    = "m5.large"
	M5XLarge   = "m5.xlarge"
	M52XLarge  = "m5.2xlarge"
	M54XLarge  = "m5.4xlarge"
	M512XLarge = "m5.12xlarge"
	M524XLarge = "m5.24xlarge"
)

var AllAwsInstances = []string{
	T2Nano,
	T2Micro,
	T2Small,
	T2Medium,
	T2Large,
	T2XLarge,
	T2XXLarge,
	M4Large,
	M4XLarge,
	M42XLarge,
	M44XLarge,
	M410XLarge,
	M416XLarge,
	M5Large,
	M5XLarge,
	M52XLarge,
	M54XLarge,
	M512XLarge,
	M524XLarge,
}

var AllAwsInstancesInEuWest3 = []string{
	T2Nano,
	T2Micro,
	T2Small,
	T2Medium,
	T2Large,
	T2XLarge,
	T2XXLarge,
	M5Large,
	M5XLarge,
	M52XLarge,
	M54XLarge,
	M512XLarge,
	M524XLarge,
}

var AllAwsInstancesInApNortheast1 = []string{
	T2Nano,
	T2Micro,
	T2Small,
	T2Medium,
	T2Large,
	T2XLarge,
	T2XXLarge,
	M4Large,
	M4XLarge,
	M42XLarge,
	M44XLarge,
	M410XLarge,
	M416XLarge,
}
