# [e-inwork.com](https://e-inwork.com)

## Golang Profile Microservice
This microservice is responsible for managing the profile data of [the Golang User Microservice](https://github.com/e-inwork-com/go-user-service).

To run both of the microservices, follow the command below:
1. Install Docker
    - https://docs.docker.com/get-docker/
2. Git clone this repository to your localhost, and from the terminal run below command:
   ```
   git clone git@github.com:e-inwork-com/go-profile-service
   ```
3. Change the active folder to `go-user-service`:
   ```
   cd go-profile-service
   ```
4. Run Docker Compose:
   ```
   docker-compose -f docker-compose.local.yml up -d
   ```
5. Create a user in the User API with CURL command line:
    ```
    curl -d '{"email_t":"jon@doe.com", "password":"pa55word", "first_name_t": "Jon", "last_name_t": "Doe"}' -H "Content-Type: application/json" -X POST http://localhost:4001/service/users
    ```
6. Login to the User API:
   ```
   curl -d '{"email_t":"jon@doe.com", "password":"pa55word"}' -H "Content-Type: application/json" -X POST http://localhost:4001/service/users/authentication
   ```
7. You will get a token from the response login, and set it as a token variable for example like below:
   ```
   token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjhhY2NkNTUzLWIwZTgtNDYxNC1iOTY0LTA5MTYyODhkMmExOCIsImV4cCI6MTY3MjUyMTQ1M30.S-G5gGvetOrdQTLOw46SmEv-odQZ5cqqA1KtQm0XaL4
   ```
8. Create a profile for current user, you can use any image or use the image on the folder test:
   ```
   curl -F profile_name_t="Jon Doe" -F profile_picture_s=@api/test/images/profile.jpg -H "Authorization: Bearer $token"  -X POST http://localhost:4002/service/profiles
   ```
9. Copy the va
lue of `profile_picture` from the response. Open it in a browser, such as http://localhost:4002/service/profiles/pictures/926d610c-fd54-450e-aa83-030683227072.jpg
10. Good luck!