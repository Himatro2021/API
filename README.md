# HIMATRO API

**Developed by:** Lucky Akbar<br>
**Built with:** Golang, Postgresql, Echo, GORM <br>
**Deployed with:** Docker, AWS, Nginx, Certbot

## How To Run

prequisite:

1. Make sure you already created .env file in the root folder. Use key from file called .env.example for example
2. Make sure you have installed docker and docker-compose. We recommend using Linux-based OS
3. Make sure you have created docker-compose.yml file. Use the docker-compose-local.yml file as the default architecture
4. Make sure you already have Postgresql instance running
5. Make sure you already create a database and a user with all privileges on that database. Set all the database user credentials in .env file.
6. Make sure you already setup a Redis instance. Set the credentials in the .env file

Steps:
1. Build docker image by running: `docker-compose up -d --build`
2. Get your container id by running command: `docker ps`. Find this server container id.
3. Execute migration command by running: `docker exec <container-id> -it go run main.go migrate`
4. Initialize super admin user by running command: `docker exec <container-id> -it go run main.go init-admin <your-email-address> ,your-admin-super-secure-password`. Replace email and password with yours. Make sure the password is strong and store it somewhere safe.
5. Your server should be ready.

## Host
The live version of this API are already proudly hosted at: **https://api.himatro.luckyakbar.tech** <br>
if you just want to test the API, you can use the staging version hosted at: **https://staging.api.himatro.luckyakbar.tech**

## Feature

1. CRUD absent
2. invite members
3. Authentication

## API CONTRACT per Feature

### CRUD Absent

1. Get absent result
  - method:
    - GET
  - description
    - Get absent result from a given form ID
  - path
    - {{ host }}/rest/absent/form/`formID`/result/
  - admin only: `false`
  - example result
    - <pre>
      {
        "title": "test cache",
        "start_at": "2002-10-11T12:30:00Z",
        "finished_at": "2022-10-12T12:45:00Z",
        "participants": [
          {
            "name": "Super Admin",
            "filled_at": "2022-08-02T18:10:33.200086Z",
            "status": "EXECUSE",
            "reason": "anjayyy"
          }
        ]
      }
    </pre>

2. Get all absent forms
  - method:
    - GET
  - description
    - get all absent form created by admin
  - path
    - {{ _.host }}/rest/absent/form/?limit=1&offset=1
    - url query: `limit` and `offset` are optional. Default to return all forms created.
  - admin only: `true`
  - example result
  - <pre>
      [
        {
          "id": 1658651284505448791,
          "participant_group_id": 1,
          "start_at": "2002-10-11T12:30:00Z",
          "finished_at": "2022-10-12T12:45:00Z",
          "title": "Test lagi ya ges",
          "allow_update_by_attendee": false,
          "allow_create_confirmation": false,
          "created_at": "2022-07-24T15:28:04.505443Z",
          "updated_at": "2022-07-24T16:09:09.557455Z",
          "deleted_at": null,
          "created_by": 1658649680660609624,
          "updated_by": 1658649680660609624,
          "deleted_by": null
        },
        {
          "id": 1659337818967544069,
          "participant_group_id": 1,
          "start_at": "2002-10-11T12:30:00Z",
          "finished_at": "2022-10-12T12:45:00Z",
          "title": "pertamax",
          "allow_update_by_attendee": false,
          "allow_create_confirmation": false,
          "created_at": "2022-08-01T14:10:18.967543Z",
          "updated_at": "2022-08-01T14:10:18.967543Z",
          "deleted_at": null,
          "created_by": 1658649680660609624,
          "updated_by": 1658649680660609624,
          "deleted_by": null
        },
        {
          "id": 1659353943732036080,
          "participant_group_id": 1,
          "start_at": "2002-10-11T12:30:00Z",
          "finished_at": "2022-10-12T12:45:00Z",
          "title": "pertamax",
          "allow_update_by_attendee": true,
          "allow_create_confirmation": false,
          "created_at": "2022-08-01T18:39:03.73203Z",
          "updated_at": "2022-08-01T18:39:03.73203Z",
          "deleted_at": null,
          "created_by": 1658649680660609624,
          "updated_by": 1658649680660609624,
          "deleted_by": null
        }
      ]
    </pre>

3. Get form by ID
  - method:
    - GET
  - description
    - Get form details by given ID
  - path
    - {{ _.host }}/rest/absent/form/`formID`
  - admin only: `false`
  - example result
    - <pre> 
      {
        "id": 1659353943732036080,
        "participant_group_id": 1,
        "start_at": "2002-10-11T12:30:00Z",
        "finished_at": "2022-10-12T12:45:00Z",
        "title": "pertamax",
        "allow_update_by_attendee": false,
        "allow_create_confirmation": false,
        "created_at": "2022-08-01T18:39:03.73203Z",
        "updated_at": "2022-08-01T18:39:03.73203Z",
        "deleted_at": null,
        "created_by": 1658649680660609624,
        "updated_by": 1658649680660609624,
        "deleted_by": null
      }
      </pre>


4. Fill absent form by Participant
  - method:
    - POST
  - path
    - {{ _.host }}/rest/absent/form/`formID`/
  - admin only: `false`
  - payload:
    - <pre>
        {
          "status": "PRESENT",
          "execuse_reason": ""
        }
      </pre>
    - value on `status` must be one of defined [here] (https://github.com/Himatro2021/API/blob/f7e961e62da35d28b09bd8e315ec225937d18b4f/internal/model/absent.go#L114)
    - example result
      - <pre>
          {
            "id": 1659438461337425656,
            "absent_form_id": 1659438431720765761,
            "user_id": 1658649680660609624,
            "created_at": "2022-08-02T18:07:41.337422278+07:00",
            "updated_at": "2022-08-02T18:07:41.337422Z",
            "status": "PRESENT",
            "reason": ""
          }
        </pre>

5. Create absent form
  - description
    - create absent form to be used
  - method
    - POST
  - path
    - {{ _.host }}/rest/absent/form/
  - admin only: `true`
  - payload:
    - <pre>
        {
          "title": "test cache",
          "start_at_date": "2002-10-11",
          "start_at_time": "12:30",
          "finished_at_date": "2022-10-12",
          "finished_at_time": "12:45",
          "group_member_id": 1
        }
      <pre>
  - example result
   -<pre>
      {
        "id": 1659438431720765761,
        "participant_group_id": 1,
        "start_at": "2002-10-11T12:30:00Z",
        "finished_at": "2022-10-12T12:45:00Z",
        "title": "test cache",
        "allow_update_by_attendee": false,
        "allow_create_confirmation": false,
        "created_at": "2022-08-02T18:07:11.72074499+07:00",
        "updated_at": "2022-08-02T18:07:11.72074499+07:00",
        "deleted_at": null,
        "created_by": 1658649680660609624,
        "updated_by": 1658649680660609624,
        "deleted_by": null
      }
    </pre>

6. Update absent form
  - description
   - use to update form infor, such as set form closed as of now, or set form closed 1 day later
  - path
    - {{ _.host }}/rest/absent/form/`formID`/
  - admin only: `true`
  - payload
    - <pre>
        {
          "title": "Test lagi ya ges",
          "start_at_date": "2002-10-11",
          "start_at_time": "12:30",
          "finished_at_date": "2022-10-12",
          "finished_at_time": "12:45",
          "group_member_id": 1
        }
      </pre>
  - example result
   -  <pre>
        {
          "id": 1658651284505448791,
          "participant_group_id": 1,
          "start_at": "2002-10-11T12:30:00Z",
          "finished_at": "2022-10-12T12:45:00Z",
          "title": "Test lagi ya ges",
          "allow_update_by_attendee": false,
          "allow_create_confirmation": false,
          "created_at": "2022-07-24T15:28:04.505443Z",
          "updated_at": "2022-07-24T16:09:09.55745572+07:00",
          "deleted_at": null,
          "created_by": 1658649680660609624,
          "updated_by": 1658649680660609624,
          "deleted_by": null
        }
      </pre>

7. Update absent list by attendee
  - description
    - update absent list by attendee, eg. change present status, etc.
  - method
    - PATCH
  - path
    - {{ _.host }}/rest/absent/form/attendance/`absentListID`/
  - admin only: `false`
  - payload
    - <pre>
        {
          "absent_form_id": 1659438431720765761,
          "status": "EXECUSE",
          "reason": "anjayyy"
        }
      </pre>
  - example result
    - <pre>
        {
          "id": 1659438461337425656,
          "absent_form_id": 1659438431720765761,
          "user_id": 1658649680660609624,
          "created_at": "2022-08-02T18:07:41.337422Z",
          "updated_at": "2022-08-02T18:10:33.200086361+07:00",
          "status": "EXECUSE",
          "reason": "anjayyy"
        }
      </pre>
    
### Invite Members

1. Check invitation exists
  - description
    - intended to be used by frontend to check is invitation with given id is valid / exits
  - method
    - GET
  - path
    - {{ _.host }}/rest/members/invitation/`invitationID`/
  - example result
    - HTTP 200 OK

2. Create invitation
  - description
    - used to send email invitation to an intended member
  - method
    - POST
  - path
    - {{ _.host }}/rest/members/invitations/
  - payload
    - <pre>
        {
          "name": "lucky akbar",
          "email": "your@mail.com",
          "role": "MEMBER"
        }
      </pre>
  - example result
    - <pre>
        {
          "id": 1659196275595280880,
          "email": "your email",
          "name": "lucky akbar",
          "invitation_code": "1659196275595273759",
          "role": "MEMBER",
          "invitation_status": "PENDING"
        }
      </pre>

### Authentication

1. Login by Email and Password
  - description
    - use to gain access to server by sending email and password
  - method
    - POST
  - payload
    - <pre>
        {
          "email": "your email",
          "password": "lucky1603"
        }
      </pre>
  - example result
    - <pre>
        {
          "access_token": "1yWdFlf23h2-__-kqsuxs0QHqfILmx3vMpS4Jsl1h8tjwwnC0JfqlJ7tUaCpbBwE0t6iZY6yCPb4pPrbqcNuSHqionGddLMBXLmUxY3bn",
          "refresh_token": "1yaFBaBwQz3-__-LQ2Ehg9553PEVGteC30jrIS97UG8gL3EVKeeegzZL5Z4T58qT5nYvnk48VMBDDbvjuYAmclPV1hm9FNZefwOv9soJB",
          "access_token_expired_at": "2022-08-09T18:04:44.624137133+07:00",
          "refresh_token_expired_at": "2022-09-01T18:04:44.62413871+07:00"
        }
      </pre>