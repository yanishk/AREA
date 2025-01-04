# Developper Documentation

## Implement a new service

Here is a simple explanation on how to implement a new service in AREA.

First of all, you must update the database used in order to implement a new service. You only have to update the "service" table with the following informations:
- The name of the service in the "name" column which must be a string
- The color you wish the service to be displayed with in the "color" column which must be an hex code
- The logo of the service in the "logo" column which must be a path toward a webp file
- You must set the columns "hasactions" and "hasreactions" to true or false depending on the available feature linked to the service
- You must set the column "isauthneeded" to true or false if the service relies on a OAuth2 flow to authenticate a user

