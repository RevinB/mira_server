package config

type Config struct {
	AppUrl       string
	FinalUrlBase string
	S3BucketName string
	JWTSecret    []byte
}
