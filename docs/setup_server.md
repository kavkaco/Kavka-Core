#  Setup Server
Let's begin to clone and configure Kavka!
 
### Prerequisites
`go-version`: `1.22`   
`docker-version`: `^24.0`   
`docker-compose-version`: `^1.29`   
 
### Installation

1. Clone `Kavka-Core` repository!

 ```bash
 git clone --depth 1 https://github.com/kavkaco/Kavka-Core.git
 ```

2. Install dependencies

 ```bash
 go mod tidy
 ```

3. Edit configuration
```bash
vim ./config/configs.yml
```

5. Start services

Lets build and start services with `docker-compose`

 ```bash
 sudo docker-compose up -d mongo redis minio

./scripts/run_devel.sh # For development

sudo docker-compose up -d app # For deployment
 ```

✅ Everything almost done.   
Kavka's back-end server is up now!

### Api Docs

You can easily read the documentation of back-end api and test it here on **Postman**!   

[https://www.postman.com/crimson-equinox-208211/workspace/kavka](https://www.postman.com/crimson-equinox-208211/workspace/kavka)
