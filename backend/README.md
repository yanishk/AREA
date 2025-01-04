# Developper Documentation

## Implement a new service

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

In the file ```/home/Chaegnal/Epitech/Tek3/B5/area/B-DEV-500-PAR-5-2-area-yanis.harkouk-mirror/backend/src/service/domain/service/service.go```:
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
