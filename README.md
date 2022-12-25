# [e-inwork.com](https://e-inwork.com)

## Getting Started
### Settings up Docker environment & run this microservice together with the [the Golang User Service](https://github.com/e-inwork-com/go-user-service)
#### This application will run two different microservices in the two different ports 4000 & 40001. For more detail, it can check on [docker-compose.py](https://github.com/e-inwork-com/go-user-service/blob/master/docker-compose.yml).
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
   docker-compose up -d
   ```
5. Create a user in the User API with CURL command line:
    ```
    curl -d '{"email":"jon@doe.com", "password":"pa55word", "first_name": "Jon", "last_name": "Doe"}' -H "Content-Type: application/json" -X POST http://localhost:4000/api/users
    ```
6. Login to the User API:
   ```
   curl -d '{"email":"jon@doe.com", "password":"pa55word"}' -H "Content-Type: application/json" -X POST http://localhost:4000/api/authentication
   ```
7. You will get a token from the response login, and set it as a token variable for example like below:
   ```
   token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjhhY2NkNTUzLWIwZTgtNDYxNC1iOTY0LTA5MTYyODhkMmExOCIsImV4cCI6MTY3MjUyMTQ1M30.S-G5gGvetOrdQTLOw46SmEv-odQZ5cqqA1KtQm0XaL4
   ```
8. Create a profile for current user, you can use any image or use the image on the folder test:
   ```
   curl -F "profile_picture=@/YourRootFolder/.../go-profile-service/api/test/profile.jpg" -H "Authorization: Bearer $token"  -X POST http://localhost:4001/api/profiles
   ```
9. The response will show a profile picture, open it on the browser for example like this http://localhost:4001/api/profiles/pictures/926d610c-fd54-450e-aa83-030683227072.jpg
10. Good luck!