package secrets

// Use this code snippet in your app.
// If you need more information about configurations or implementing the sample code, visit the AWS docs:
// https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/setting-up.html

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/smithy-go"
	log "github.com/sirupsen/logrus"
)

type errBinaryNotStringSecret struct{}
type errStringNotBinarySecret struct{}

func (m *errBinaryNotStringSecret) Error() string {
	return "secret is of type binary, expected string"
}

func (m *errStringNotBinarySecret) Error() string {
	return "secret is of type string, expected binary"
}

func NewBinaryNotStringSecretError() *errBinaryNotStringSecret {
	return &errBinaryNotStringSecret{}
}

func NewStringNotBinarySecretError() *errStringNotBinarySecret {
	return &errStringNotBinarySecret{}
}

func GetBinarySecret(secretName, region string) (string, error) {
	ctx := context.TODO()

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		log.Errorf("configuration error, " + err.Error())
		return "", err
	}

	client := secretsmanager.NewFromConfig(cfg)

	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	}

	// In this sample we only handle the specific exceptions for the 'GetSecretValue' API.
	// See https://docs.aws.amazon.com/secretsmanager/latest/apireference/API_GetSecretValue.html

	result, err := client.GetSecretValue(ctx, input)
	if err != nil {
		var oe *smithy.OperationError
		if errors.As(err, &oe) {
			log.Errorf("aws: %v", err)
		} else {
			log.Errorf("general error: %v", err)
		}
		return "", err
	}

	// Decrypts secret using the associated KMS key.
	// Depending on whether the secret is a string or binary, one of these fields will be populated.
	var decodedBinarySecret string
	if result.SecretBinary != nil {
		decodedBinarySecretBytes := make([]byte, base64.StdEncoding.DecodedLen(len(result.SecretBinary)))
		len, err := base64.StdEncoding.Decode(decodedBinarySecretBytes, result.SecretBinary)
		if err != nil {
			log.Errorf("Base64 Decode Error: %v", err)
			return "", err
		}
		decodedBinarySecret = string(decodedBinarySecretBytes[:len])
	} else {
		log.Error(NewStringNotBinarySecretError())
		return "", err
	}

	return decodedBinarySecret, nil
}

func GetSecretValues(secretName, region string) (*map[string]interface{}, error) {
	ctx := context.TODO()

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		log.Errorf("configuration error, " + err.Error())
		return nil, err
	}

	client := secretsmanager.NewFromConfig(cfg)

	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	}

	// In this sample we only handle the specific exceptions for the 'GetSecretValue' API.
	// See https://docs.aws.amazon.com/secretsmanager/latest/apireference/API_GetSecretValue.html

	result, err := client.GetSecretValue(ctx, input)
	if err != nil {
		var oe *smithy.OperationError
		if errors.As(err, &oe) {
			log.Errorf("aws: %v", err)
		} else {
			log.Errorf("general error: %v", err)
		}
		return nil, err
	}

	// Decrypts secret using the associated KMS key.
	// Depending on whether the secret is a string or binary, one of these fields will be populated.
	var secretString string
	if result.SecretString != nil {
		secretString = *result.SecretString
	} else {
		log.Error(NewBinaryNotStringSecretError())
	}

	secretValues := &map[string]interface{}{}

	err = json.Unmarshal([]byte(secretString), secretValues)

	if err != nil {
		log.Errorf("error unmarshaling secret string to map: %v", err)
	}

	// Your code goes here.
	return secretValues, nil
}
