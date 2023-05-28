# goboyle

Boiler plate for golang API / Psql DB using [go-chi](https://github.com/go-chi/chi) for routing and [bun](https://github.com/uptrace/bun) for ORM. 

### Demo
1. Start by setting up the DB, instructions can be found [here](https://github.com/edstef/goboyle/tree/master/models)

<br>

2. Build the code, then run the API
```
    go build
    ./goboyle
```
>Ensure the API is running
```
curl --location 'localhost:8080/ping'
```
<br>

3. Next, create a profile
```
curl --location 'localhost:8080/profile' \
--header 'Content-Type: application/json' \
--data '{
    "name": "edstef"
}'
```

>The response should look something like:
```
{"id":"e6899d19-656b-45be-b88a-a5883157dc66","name":"edstef","picture_url":"/defaults/1","theme":"default_theme_1"}
```
<br>


4. Using the `id` from the response above, create a JWT for this id:
```
curl --location 'localhost:8080/get_jwt/e6899d19-656b-45be-b88a-a5883157dc66'
```


>The response will look something like:
```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwcm9maWxlX2lkIjoiZTY4OTlkMTktNjU2Yi00NWJlLWI4OGEtYTU4ODMxNTdkYzY2In0.19_WmcsGhRaFHtcZ0DWYY9Ct685TyKRUMzTyF1XF6z4
```
<br>


5. Finally, decode the jwt
```
curl --location 'localhost:8080/decode_jwt' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwcm9maWxlX2lkIjoiZTY4OTlkMTktNjU2Yi00NWJlLWI4OGEtYTU4ODMxNTdkYzY2In0.19_WmcsGhRaFHtcZ0DWYY9Ct685TyKRUMzTyF1XF6z4'

```
>You should see the original `profile id` in the response:

```
{"profile_id":"e6899d19-656b-45be-b88a-a5883157dc66"}
```
<br>

6. You can then use this JWT to also refetch the profile data from the DB
```
curl --location 'localhost:8080/profile' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwcm9maWxlX2lkIjoiMzNkOTZiY2MtYzAxYi00ZTc2LWJiNmYtOTRjZGI4NmYzNGJhIn0.fJHjBpYR6BcnZUkNb4ESfBZm7uq0r7bC6NnvDf2Ms3E'
```
```
{"id":"e6899d19-656b-45be-b88a-a5883157dc66","name":"","picture_url":"/defaults/1","theme":"default_theme_1"}
```
