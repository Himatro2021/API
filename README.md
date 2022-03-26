# HIMATRO API

**Developed by:** Lucky Akbar (lucky.akbar105619@students.unila.ac.id)<br>
**Built with:** Golang, Postgresql, Echo, GORM <br>
**Deployed with:** Docker, AWS, Nginx, Certbot

## How To Run

prequisite:

1. Make sure you use Linux based operating system. This service only intended to run in Linux environment. If you use another OS, feel free to make a pull request to add the extra steps
2. Install docker on your machine
3. Install docker-compose.
4. Have required private data to initialize the database. This file locations are defined in your .env file. So feel free to store your private data. This location should be accessed by the API server. Please refer to [here](#defined-private-data-to-initialize-database) to create your own.
5. To create super admin credentials, you can utilize our encryptor utility. To use is, just prepare your admin password, and then run this command `./cmd/encryptor`. After that, the program will ask you to type your password. The result of this program is the encrypted version of the password. So you can use that result and store it in the superAdmin.csv file.

Steps:

1. Clone or pull this repository
2. Populate your .env file based on structure defined on .env.example<br>
   For example, copy all the content of .env.example, and paste it on a new .env file. Fill all the variable value with approprieate value.
3. Create docker-compose.yml file. Copy content of docker-compose-local.yml, and paste it to docker-compose.yml. Uncomment all the commented config in docker-compose-local.yml. _note_ you should change the value of the environment in postgres service. The value there should match with the value on your .env file, such as the value in .env, PGDATABASE should match with the environment value defined in POSTGRES_DB.
4. Run this command from the terminal on the root file of the source code -> `docker-compose up`. You should see the last output similar to this: <br>
   `postgres_1 | 2022-03-18 01:20:43.446 UTC [1] LOG: database system is ready to accept connections`<br>
   `himatro-api_1 | 2022/03/18 08:20:44 Successfully connected to Postgres Server`<br>
   `himatro-api_1 | 2022/03/18 08:20:44 Server listening on port :8080`
5. After that, you should know your container's id for Himatro API. Type `docker ps`
6. Find Himatro Api container id, and execute `docker exec -it your_container_id /app/bin/main migrate` to run db migrator
7. Execute `docker exec -it your_container_id /app/bin/main seeder` to run db seeder
8. Your API should be accessible from post 8080 in your machine.

## Host

The live version of this API are already proudly hosted at: **https://api.himatro.luckyakbar.tech** <br>
if you just want to test the API, you can use the staging version hosted at: **https://staging.api.himatro.luckyakbar.tech**

## Feature

1. CRUD absent form
2. Fill absent list
3. Login as Admin

## API CONTRACT per Feature

### Admin menu

- #### Login as Admin <br>

  - Route: **/login**
  - Method: **POST**
  - Accepted Content Type / Payload: **application/json**
  - URL params: **none**
  - URL query: **none**
  - Payload <br>
    1. NPM
       - type: numeric string
       - required: true
    2. password
       - type: string
       - required: true
  - Return type: **JSON**
  - Success Response Payload:<br>
    1. ok
       - type: boolean
    2. token
       - type: string
  - Failed Response Payload:<br>
    1. ok
       - type: boolean
    2. message
       - type: string
  - Note:<br>
    When using this route, remember to use the received token as bearer authorization to access restricted resource
    <br> <br>

- #### Create Absent Form <br>

  - Route: **/admin/absensi**
  - Method: **POST**
  - Accepted Content Type / Payload: **application/json**
  - URL params: **none**
  - URL query: **none**
  - Payload <br>
    1. title
       - type: string
       - required: true
       - case sensitive: yes
    2. startAtDate
       - type: string
       - required: true,
       - format: `YYYY-MM-DD`
    3. startAtTime
       - type: string
       - required: true,
       - format: `HH:MM:SS`
    4. finishtAtDate
       - type: string
       - required: true,
       - format: `YYYY-MM-DD`
    5. finishtAtTime
       - type: string
       - required: true,
       - format: `HH:MM:SS`
    6. requireAttendanceProof
       - type: boolean
       - required: false
       - allowed values: **"true"** or **"false"**
       - default: false
    7. requireExecuseProof
       - type: boolean
       - required: false
       - allowed values: **"true"** or **"false"**
       - default: false
    8. participant
       - type: string
       - required: true
       - allowed values: see [here](#defined-departement-name)
       - case sensitive: no
  - Success Response Payload: <br>
    1. ok: boolean
    2. absentID: int
    3. title: string
    4. participant: int
    5. startAt: date
    6. finishAt: date
    7. requireAttendanceImageProof: boolean
    8. requireExecuseImageProof: boolean
       <br> <br>
  - Note:<br>
    You have to strictly follow the rules, format or allowed values defined in each payload. If there is some validation error, server will return error message regarding what is error and will give you **404 Bad Request** response.
    <br><br>

- #### Update Title from Absent Form

  - Route: **/admin/absensi/:absentID/title**
  - Method: **PATCH**
  - Accepted Content Type / Payload: **application/json**
  - URL params: <br>
    1. absentID
       - type: numeric string
       - required: true
  - URL query: **none**
  - Payload <br>
    1. title
       - type: string
       - required: true
  - Success Response Payload: <br>
    1. ok: boolean
    2. message: string
    3. fieldName: string
    4. value: string
  - Note:<br>
    Updating title of an Absent Form can only be done with existing absent form and accessed via it's absentID. If you're tempting to perform update on non existing absent form, server will return error message with **400 Bad Request** response. The title will have the same as the value you send, and uppercase / lowercase character will not modified.
    <br><br>

- #### Get Form Absent Details <br>
  - Route: **/admin/absensi**
  - Method: **GET**
  - Accepted Content-Type / Payload: **none**
  - URL Params: **none**
  - URL Query: <br>
    1. limit
    - type: numeric string
    - required: false
    - default: null (no limit)
  - Payload: **none**
  - Success Response Payload: 1. ok: boolean 2. message: string 3. list:
    <br> array of: - form_id: int - title: string - created_at: date string - updated_at: date string - participant_code: int - require_attendance_image_proof: boolean - require_execuse_image_proof: boolean - total_participant: int - hadir: int - izin: int - tanpa_keterangan: int
    <br><br>
- #### Update Finish At from Absent Form
  - Route: **/admin/absensi/:absentID/finishAt**
  - Method: **PATCH**
  - Accepted Content Type / Payload: **application/json**
  - URL params: <br>
    1. absentID
       - type: numeric string
       - required: true
  - URL query: **none**
  - Payload <br>
    1. date
       - type: string
       - required: true
       - format: `YYYY-MM-DD`
    2. time
       - type: string
       - required: true
       - format: `HH:MM:SS`
  - Success Response Payload: <br>
    1. ok: boolean
    2. message: string
    3. fieldName: string
    4. value: string
  - Note:<br>
    Server will do several validation to values given in the payload. Some are the value when converted to date, must not due before the form start time. You will receive **400 Bad Request** along with error message if the validation returns error.
    <br><br>
- #### Update Start At from Absent Form
  - Route: **/admin/absensi/:absentID/startAt**
  - Method: **PATCH**
  - Accepted Content Type / Payload: **application/json**
  - URL params: <br>
    1. absentID
       - type: numeric string
       - required: true
  - URL query: **none**
  - Payload <br>
    1. date
       - type: string
       - required: true
       - format: `YYYY-MM-DD`
    2. time
       - type: string
       - required: true
       - format: `HH:MM:SS`
  - Success Response Payload: <br>
    1. ok: boolean
    2. message: string
    3. fieldName: string
    4. value: string
  - Note:<br>
    Server will do several validation to values given in the payload. Some are the value when converted to date, must not comes after the form end date. You will receive **400 Bad Request** along with error message if the validation returns error.
    <br><br>
- #### Update Participant from Absent Form
  - Route: **/admin/absensi/:absentID/participant**
  - Method: **PATCH**
  - Accepted Content Type / Payload: **application/json**
  - URL params: <br>
    1. absentID
       - type: numeric string
       - required: true
  - URL query: **none**
  - Payload <br>
    1. participant
       - type: string
       - required: true
       - allowed values: see [here](#defined-departement-name)
  - Success Response Payload: <br>
    1. ok: boolean
    2. message: string
    3. fieldName: string
    4. value: string
  - Note:<br>
    Server will validate if you supply the right value in the payload. The form participant can't be changed if you are trying to change the participant when one or more participants already fill the absent. You will receive **400 Bad Request** along with error message if the validation returns error.

<br><br>

- #### Update Attendance Proof from Absent Form
  - Route: **/admin/absensi/:absentID/attendanceImageProof**
  - Method: **PATCH**
  - Accepted Content Type / Payload: **application/json**
  - URL params: <br>
    1. absentID
       - type: numeric string
       - required: true
  - URL query: **none**
  - Payload <br>
    1. status
       - type: boolean
       - required: true
       - allowed values: **"true"** or **""false""**
       - default to: false
  - Success Response Payload: <br>
    1. ok: boolean
    2. message: string
    3. fieldName: string
    4. value: string
  - Note:<br>
    You can change the status of whether the participant must send image or not when they are attend the event. This change will not affect people who are already fill the form. This feature will automatically active when the backend already enable this feature. You will receive **400 Bad Request** along with error message if the validation returns error.
    <br><br>
- #### Update Execuse Proof from Absent Form
  - Route: **/admin/absensi/:absentID/execuseImageProof**
  - Method: **PATCH**
  - Accepted Content Type / Payload: **application/json**
  - URL params: <br>
    1. absentID
       - type: numeric string
       - required: true
  - URL query: **none**
  - Payload <br>
    1. status
       - type: boolean
       - required: true
       - allowed values: **"true"** or **""false""**
       - default to: false
  - Success Response Payload: <br>
    1. ok: boolean
    2. message: string
    3. fieldName: string
    4. value: string
  - Note:<br>
    You can change the status of whether the participant must send image or not when they are can't attend the event. This change will not affect people who are already fill the form. This feature will automatically active when the backend already enable this feature. You will receive **400 Bad Request** along with error message if the validation returns error.
    <br><br>

### Non Admin Menu

- ### Check Form Absent is Writeable

  - Route: **/absensi/:absentID**
  - Method: **GET**
  - Accepted Content Type / Payload: **application/json**
  - URL params: <br>
    1. absentID
       - type: numeric string
       - required: true
  - URL query: **none**
  - Payload: **none**
  - Success Response Payload: **none**
  - Note:<br>
    This endpoint used to check wheter the user can fill the absent form. If the server returns **200 OK**, you should render absent form filling page. Otherwise, it means that the form is not writeable, so server will response with **400 Bad Request** alongside with error message.

- ### Get Absent Form Result

  - Route: **/absensi/:absentID/result**
  - Method: **GET**
  - Accepted Content Type / Payload: **application/json**
  - URL params: <br>
    1. absentID
       - type: numeric string
       - required: true
  - URL query: **none**
  - Payload: **none**
  - Success Response Payload: <br>
    1. ok: boolean
    2. formID: int
    3. total: int
    4. list: array -> npm, updatedAt, keterangan, nama, departemen (all string)
  - Note:<br>
    This endpoint will give you absent result no matter if the form it self is already closed or not even open yet. The field total in the response payload represent how many participants are in the list. If you're trying to request inexisting / request with invalid _absentID_, server will response with **400 Bad Request** alongside with error message.

- ### Fill Absent Form

  - Route: **/absensi/:absentID**
  - Method: **POST**
  - Accepted Content Type / Payload: **application/json**
  - URL params: <br>
    1. absentID
       - type: numeric string
       - required: true
  - URL query: **none**
  - Payload <br>
    1. NPM
       - type: string
       - required: true
    2. keterangan
       - type: string
       - required: true
       - allowed values: **"h"** or **"i"**
  - Success Response Payload: **none**
  - Note:<br>
    This endpoint will only accept right NPM as a proof that this NPM are owned by registered himatro's member and also one of the expected absent attendance. This endpoint can only be used to fill the absent form if the participant never filled the absent before. If you need to change the absent list after filling, you should use the **PATCH** method. Also this endpoint will give you update absent list token to be use when you want to change / update your presence status.

- ### Update Absent List
  - Route: **/absensi/:absentID**
  - Method: **PATCH**
  - Accepted Content Type / Payload: \*_application/json_
  - URL params: <br>
    1. absentID
       - type: numeric string
       - required: true
  - URL query: **none**
  - Payload <br>
    1. keterangan
       - type: string
       - required: true
       - allowed values: **"h"** or **"i"**
  - Success Response Payload: **none**
  - Note:<br>
    This endpoint will only accept your payload and read your update absent list token cookie. If there is error or absence in your token, you will not able to update your presence status. If server accepts your request, it will give you only **202 Accepted** response.

## Defined Departement Name

1. Pengurus Harian -> PH
2. Departemen Pendidikan dan Pengembangan Diri -> PPD
3. Departemen Kaderisasi Dan Pengembangan Organisasi -> KPO
4. Departemen Komunikasi Dan Informasi -> KOMINFO
5. Departemen Sosial dan Kewirausahaan -> KWU
6. Departemen Pengembangan Keteknikan -> BANGTEK

extra: If you want to create absent form for all members, use ALL in the participant payload.

## Defined private data to initialize database

There are several data required to run the server. Simply, you can look at .env.example content, and see the variable named like this: `ANGGOTA_BIASA_SEEDER_DATA_PATH`, or any variable ends with `SEEDER_DATA_PATH`. This file needs to be a csv file. <br>
Because csv file need to have header, the header are defined in the .env.example file which named something like this `CSV_HEADER_CONFIG`. You should strictly follow this headers.
<br><br>

## API TEST RESULT

the result of testing against this API is located at: [here](https://docs.google.com/document/d/1klvQ0hIYvDtcUUbGUpmdMFxWZQV6ky8As_LxDNOqFkc/edit?usp=sharing)

## TO DO LIST

1. Add order_by query params in get absent forms details
2. Add length validator in field validation time (it always the same regardless the year)
3. Add release on github
4. Add security feature to store sesion token and use extra data e.g ip and any other to make as unique as possible to one user only
5. Update how to run section
