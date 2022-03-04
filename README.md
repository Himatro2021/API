# API CONTRACT

1. /login
  - POST
  - Content-Type: application/x-www-form-urlencoded
  - Form:
    <br>a. NPM (string)
    <br>b. password (string)
  - Returns: token

2. /admin
  - GET
  - Require login
  - Authentication method: Bearer


## Host
api.himatro.luckyakbar.tech