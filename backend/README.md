# Developper Documentation

## Implement a new service

### Database

Here is a simple explanation on how to implement a new service in AREA.

First of all, you must update the database used in order to implement a new service. You only have to update the "service" table with the following informations:
- The name of the service in the "name" column which must be a string
- The color you wish the service to be displayed with in the "color" column which must be an hex code
- The logo of the service in the "logo" column which must be a path toward a webp file
- You must set the columns "hasactions" and "hasreactions" to true or false depending on the available feature linked to the service
- You must set the column "isauthneeded" to true or false if the service relies on a OAuth2 flow to authenticate a user

### OAuh2

If your newly implemented service requires OAuth2, here are the step to follow to implement it in AREA.

#### Require the access code:

In the file ```backend/src/service/domain/service/service.go```:
- Implement a function that returns the link to request an access code, you can use the function
```go
func getCallbackAndClientId(callbackType string, serviceName string, isIdNecessary bool) (string, string)
```
to get the callback and the <service>_CLIENT_ID and the <service>_CLIENT_SECRET if necessary.

> [!NOTE]
> To use our helper functions, such as getCallbackAndClientId and many other, the environment variable must be named a certain way and must contain two callback, one for our mobile plateform and one for our website. Those naming convention can be found in the .env.example at the root of the "backend" folder.

> [!NOTE]
> The serviceName variable must be entirely capitalized.

- In the same file, go to
```go
func (self *ServiceService) OAuth2Service(serviceName, callbackType, appType string) (string, error) {
```
and modify the switch case to compare to the name of your service and add the function you did in the step above.

#### Exchange the access code for an access token

The second step of any OAuth2 flow is the exchange of the access code for an access token that can be stored in our database.

In the file ```backend/src/service/domain/service/service.go```:
- Create the function that will handle the logic of the exchange.
- Get the callback URI from the .env according to the plateform used: "web" or "mobile"
- Use the function:
```go
func genericAccessTokenRequest(tokenUrl, code, callbackUrl string) (*http.Request, error) {
```
-Set your headers and return the *http.Request returned by "genericAccessTokenRequest".
> [!NOTE]
> In some specific cases, depending on the service you are trying to implement, this function cannot be used but its logic can.

#### Refresh the token

> [!NOTE]
> Some services do not implement a refresh token, be sure to read the documentation of the service you are implementing. If the service does not provide refresh token you can skip this step.

In the file ```backend/src/service/domain/service/service.go```:
- Create a function that will handle the logic of the refreshing of the token, it must be prototyped this way
```go 
func (self *ServiceService) <YOUR DESIRED FUNCTION NAME>(refreshToken string) (*http.Request, error)
```
- You can use one of two function depending on your service requirement:
```go
func (self *ServiceService) genericRefreshTokenRequest(tokenUrl, refreshToken string) (*http.Request, error)
```
or
```go
func (self *ServiceService) genericJsonBodyRefreshToken(clientId, secretId, refreshToken string) string
```
> [!NOTE]
> The second function is to be used when the refresh method requires the CLIENT_ID and CLIENT_SECRET variable.

- Set your header and return the "*http.Request".

In the file ```backend/src/service/service.go```:
- In the ServiceService interface, add the newly created function prototyp to the list.

In the file ```backend/src/service/domain/userservice/userservice.go```:
- In the function
```go
func (self *UserServiceService) refreshToken(refreshToken, userId, serviceName string) (entities.ResultToken, error)
```
add your service to the switch case and return the function that you created in the first step.

## Implement a new action

### Database

First of all, to implement a new action you must update the database. In the table "actions":
- Fill the "serviceid" column with the ID of the service that you want this action to belon to
> [!NOTE]
> The ID of a service can be found in the "service" table in the column "id"
- Fill the "name" column with the name you wish to give to the action. Beware, this will be the name displayed on the front and mobile plateform
- Fill the "description" column, this column will be displayed in the front.
> [!NOTE]
> Think about what you want your user to know, as it will be displayed to the user. A good practice is to put the limitation of the action in the description.
- Fill the "nbparam" column with the number of input expected of the user in order to create the action correctly
- Fill the "parameters" column. In this column, you must give the name, type (string, int), route that must be called, any pre-conceived value and if those values are exhaustive for each parameter.
