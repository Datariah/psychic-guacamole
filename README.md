# psychic-guacamole
What happens when you need to rotate secrets in your application and database?

## Requirements
- Terraform (recommended: terraform version > 1.0.0, aws provider ~> 4.0.0)
- Go 1.18
- AWS Credentials
- Internet Access

## Usage

### Terraform

In order for things to function, you will need to match the region in [terraform/providers.tf](terraform/providers.tf)
and the region in the application (hardcoded to a constant in [internal/constants.go](internal/constants.go)).

This workspace will create VPC resources (VPC, Routes, Route Tables, NAT Gateway, Internet Gateway, etc), RDS Instance, 
Lambda functions, IAM Roles, etc.

However, once the RDS Instance is created, you will need to put the credentials on a  MySQL RDS secret in Secrets Manager.

Afterwards, please change the `count` parameter from `0` to `1` in `aws_secretsmanager_secret_rotation.rds_secret` and
`data.aws_secretsmanager_secret.rds` located in [terraform/secrets.tf](terraform/secrets.tf). This will enable the
rotation on the secret you create and specify in the `data.aws_secrets_manager_secret.rds` resource.

### Application

If you happen to change the secret name in the data source, make sure you update the secret name in 
[pkg/server/helpers/database.go](pkg/server/helpers/database.go) for the RDS secret and in 
[pkg/server/entrypoint.go](pkg/server/entrypoint.go) for the JWT Token secret.

The JWT token secret needs to have a key called `JWT_SECRET` with a random string. This will be used to validate the
signature from the JWT Token you use to authenticate. 

You can create a token using that same secret in [jwt.io](https://jwt.io) to authenticate your requests 
(`Authorization: Bearer <JWT_TOKEN>`).

The application will listen for requests on port 8081 and the following routes are implemented:

- `GET /api/v1/persons`: Retrieve all persons
```shell
curl -X http://localhost:8081/api/v1/persons
```
- `POST /api/v1/persons`: Create a person
```shell
curl -X POST http://localhost:8081/api/v1/persons --data-raw '{"first_name": "John","last_name": "Doe"}'
```
- `PUT /api/v1/persons`: Update a person
```shell
curl -X PUT http://localhost:8081/api/v1/persons --data-raw '{"id": 1, "first_name": "John","last_name": "Doe"}'
```
- `DELETE /api/v1/persons/{id}`: Delete a person by ID
```shell
curl -X DELETE http://localhost:8081/api/v1/persons/1
```
