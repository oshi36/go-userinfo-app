# Go UserInfo Application
This is a Go application that uses MongoDB as its database and a [gorilla/mux go package](https://pkg.go.dev/github.com/gorilla/mux) for handling routes.

## Steps To Run Application
- Install Golang ,Docker and Docker Compose on your machine.
- Clone this GitHub repository
```
git clone https://github.com/cloudyuga/go-userinfo-app.git
```
- Run `go mod tidy` to download dependencies inside the clone project directory.
- Run mongodb as a container.
```
docker container run -d -p 27017:27017 --name=mongodb mongo:3.6
```
- Run the application locally.
```
go run main.go
```
- Access the application with `http://localhost:80` on your browser, and you will application as shown below.
  
![go-user-data](https://github.com/oshi36/go-userinfo-app/assets/47573417/43f87e53-0e98-426e-83ca-4241bc1b5f7f)

**NOTE**: If port 80 is occupied then change it to some other port in the `main.go` file.

- Check the `Dockerfile` on how  to build its docker image and `compose.yaml` file for running the whole application stack together.
