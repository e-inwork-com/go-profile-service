# e-inwork.com 
# Go Profile Service

### Create a user:
```
curl -d '{"email":"jon@doe.com", "password":"pa55word", "first_name": "Jon", "last_name": "Doe"}' -H "Content-Type: application/json" -X POST http://localhost:4000/api/users
``` 

### Get a authorization token:
```
curl -d '{"email":"jon@doe.com", "password":"pa55word"}' -H "Content-Type: application/json" -X POST http://localhost:4000/api/authentication
```

### Create a profile:
```
curl -F "profile_picture=@/Users/.../go-profile-service/api/test/profile.jpg" -H "Authorization: Bearer {token}"  -X POST http://localhost:4001/api/profiles
```

### Get a user profile:
```
curl -H "Authorization: Bearer {token}" -X GET http://localhost:4001/api/profiles/me
```

### Update a user profile:
```
curl -F "profile_picture=@/Users/.../go-profile-service/api/test/profile.png" -H "Authorization: Bearer {token}"  -X PATCH http://localhost:4001/api/profiles/{id}
```
