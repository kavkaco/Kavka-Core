<br />

<div align="center">
  <a href="https://github.com/kavkaco">
    <img src="./docs/logo.png" alt="Logo" width="80" height="80">
  </a>

  <h3 align="center">Kavka</h3>

  <p align="center">
    Messaging Application
    <br />
    <a href="https://github.com/kavkaco/Docs"><strong>Explore the docs »</strong></a>
    <br />
    <br />
    <a href="https://github.com/kavkaco/Kavka-Backend/issues">Report Bug</a>
  </p>
</div>
 
## Introduction
 
Kavka is a **feature-rich**, **open-source** chat application developed with **GoLang**, and it is proudly hosted on GitHub. As an open-source project, Kavka welcomes contributors from all over the world to join the community and contribute their skills and expertise to make the application even better.
 
At Kavka, we believe in the power of collaboration and continuous learning. Our primary goal is to provide a platform where developers can not only connect and communicate but also enhance their programming skills and strive towards becoming senior programmers. We understand the importance of continuous growth in the ever-evolving field of technology, and Kavka aims to be a catalyst in that journey.
 
With a focus on learning, Kavka offers a range of features designed to challenge and expand developers' knowledge. From implementing real-time messaging using `websocket` to ensuring data security and encryption, Kavka provides an environment where developers can experiment with new technologies and stay up-to-date with the latest trends in the industry.
 
## Built With

[![My Skills](https://skillicons.dev/icons?i=vscode,golang,docker,nginx,git,github,postman,mongodb,redis,vuejs,nuxtjs,ts,aws)](https://skillicons.dev)
 
## Getting Started

Let's begin to clone and configure Kavka on a local machine!
 
### Prerequisites

`go-version` `1.18`   
`docker-version`: `^24.0`   
`docker-compose-version`: `^1.29`   
 
### Installation

1. Clone the repo

 ```bash
 git clone --depth 1 https://github.com/kavkaco/Kavka-Backend.git
 ```

2. Install dependencies

 ```bash
 go mod tidy
 ```

3. Build databases using by docker

 ```bash
 sudo docker-compose up -d redis
 sudo docker-compose up -d mongo
 ```

3. Build MinIO service

 ```bash
 sudo docker-compose up -d minio
 ```

### Setup

Everything almost done. You can easily run the backend server on your system!

```bash
sudo chmod +x ./scripts/run_devel.sh
./scripts/run_devel.sh
```

## Postman

[https://www.postman.com/crimson-equinox-208211/workspace/kavka](https://www.postman.com/crimson-equinox-208211/workspace/kavka)
