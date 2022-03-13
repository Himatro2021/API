# HIMATRO API
**Developed by:** Lucky Akbar (lucky.akbar105619@students.unila.ac.id)<br>
**Built with:** Golang, Postgresql, Echo, GORM <br>
**Deployed with:** Docker, AWS, Nginx, Certbot

## How To Run
- to do...

## Host
The live version of this API are already proudly hosted at: **https://api.himatro.luckyakbar.tech**

## Feature
1. CRUD absent form
2. Fill absent list
3. Login as Admin

## API CONTRACT per Feature

### Admin menu

- #### Login as Admin <br>
  - Route: **/login**
  - Method: **POST**
  - Accepted Content Type / Payload: **form url encoded**
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
  - Accepted Content Type / Payload: **form url encoded**
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
        - type: string
        - required: false
        - allowed values: **"true"** or **"false"**
        - default value: **"false"**
    7. requireExecuseProof
        - type: string
        - required: false
        - allowed values: **"true"** or **"false"**
        - default value: **"false"**
    8. participant
        - type: string
        - required: true
        - allowed values: see [here](#defined-departement-name)
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
  - Accepted Content Type / Payload: **form url encoded**
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
- #### Update Finish At from Absent Form
  - Route: **/admin/absensi/:absentID/finishAt**
  - Method: **PATCH**
  - Accepted Content Type / Payload: **form url encoded**
  - URL params: <br>
    1. absentID
        - type: numeric string
        - required: true
  - URL query: **none**
  - Payload <br>
    1. finishAtDate
        - type: string
        - required: true
        - format: `YYYY-MM-DD`
    2. finishAtTime
        - type: string
        - required: true
        - format: `HH:MM:SS`
  - Success Response Payload: <br>
    1. ok: boolean
    2. message: string
    3. fieldName: string
    4. value: string
  - Note:<br>
  Server will do several validation to values given in the payload. Some are the value when converted to date, must not due before the form start time. Also, it can't end before now. You will receive **400 Bad Request** along with error message if the validation returns error.
<br><br>
- #### Update Start At from Absent Form
  - Route: **/admin/absensi/:absentID/startAt**
  - Method: **PATCH**
  - Accepted Content Type / Payload: **form url encoded**
  - URL params: <br>
    1. absentID
        - type: numeric string
        - required: true
  - URL query: **none**
  - Payload <br>
    1. startAtDate
        - type: string
        - required: true
        - format: `YYYY-MM-DD`
    2. startAtTime
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
  - Accepted Content Type / Payload: **form url encoded**
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
  - Accepted Content Type / Payload: **form url encoded**
  - URL params: <br>
    1. absentID
        - type: numeric string
        - required: true
  - URL query: **none**
  - Payload <br>
    1. status
        - type: string
        - required: true
        - allowed values: **"true"** or **""false""**
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
  - Accepted Content Type / Payload: **form url encoded**
  - URL params: <br>
    1. absentID
        - type: numeric string
        - required: true
  - URL query: **none**
  - Payload <br>
    1. status
        - type: string
        - required: true
        - allowed values: **"true"** or **""false""**
  - Success Response Payload: <br>
    1. ok: boolean
    2. message: string
    3. fieldName: string
    4. value: string
  - Note:<br>
  You can change the status of whether the participant must send image or not when they are can't attend the event. This change will not affect people who are already fill the form. This feature will automatically active when the backend already enable this feature. You will receive **400 Bad Request** along with error message if the validation returns error.
  <br><br>

## Defined Departement Name
1. Pengurus Harian -> PH
2. Departemen Pendidikan dan Pengembangan Diri -> PPD
3. Departemen Kaderisasi Dan Pengembangan Organisasi -> KPO
4. Departemen Komunikasi Dan Informasi -> KOMINFO
5. Departemen Sosial dan Kewirausahaan -> KWU
6. Departemen Pengembangan Keteknikan -> BANGTEK

