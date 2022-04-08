package config

type Config struct {
	AppUrl           string
	FinalUrlBase     string
	S3BucketName     *string
	CloudfrontDistID *string
	JWTSecret        []byte
}
