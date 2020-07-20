# iti-challenge
This project is an challenge for job as software engineer at bank Itau. I'm need implement API Rest to receive a password and apply validation rules.

## Business rules
- [x] 9 or more than characters
- [x] 1 or more digits
- [x] 1 or more upper char
- [x] 1 or more lower char
- [x] 1 or more special char (!@#$%^&*()-+)
- [x] Not have repeated char 

## Requirements 
- [x] API Rest
- [x] Unit test
- [x] Integration test

## Endpoints

**URL** : `/v1/api/password-validate`

**Method** : `POST`

```json
{
    "value": "string"
}
```

## Success Response

**Code** : `200 OK`

```json
{
    "valid": false,
}
```

## Error Responses

**Condition** : if send unmarshal json

**Code** : `400 BAD REQUEST`

**Code** : `500 INTERNAL SERVER ERROR`



## Requirements for run
* install go (only for run tests)
* install docker

## For run tests
```shell
go test ./... -coverprofile test_coverage.out
```

## For run application
```shell
docker build -t app .
docker run -p 8080:8080 --name my-app app
curl -XPOST -d '{"value":"your-password"}'  http://localhost:8080/v1/api/password-validate
``` 




  


