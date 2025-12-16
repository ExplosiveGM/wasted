# Using 
1. Install golang
2. Install postgres
3. Create .env file with settings
4. Install make
5. Use `make migration-up`
6. After that use `make server-start`
7. Currently you can use 
* POST localhost:8080/api/v1/auth/request-code
{
	"login": "Email|phoneNumber"
}
* POST localhost:8080/api/v1/auth/verify
{
	"login": "Email|phoneNumber",
	"code": 123456
}
* POST localhost:8080/api/v1/auth/refresh
{
	"refresh_token":"refreshTokenValue"
}