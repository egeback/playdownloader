# Play downloader

Microservice to download play serivice media files (series, news, movies etc.). The service wraps arround [svtplay-dl](https://svtplay-dl.se/) executable.

## Installation

### Clone repo

``` bash
git clone https://github.com/egeback/playdownloader.git
```

#### Update configuration

Update ./config/config.yaml

## Deployment options

The microservice can be deployed as standlone application or in a docker container

### Standalone golang application

#### Install Svtplay-dl

Follow [install page](https://svtplay-dl.se/install/)  
  
#### Run build script in from root director

``` bash
./cmd/build.sh
```

#### Run application

``` bash
./mediadownloader
```

### Docker container

#### Configure Docker Container

Update Dockefile (update ports and location of media)

#### 1. Using docker-compose ([link](https://www.google.com/url?sa=t&rct=j&q=&esrc=s&source=web&cd=&cad=rja&uact=8&ved=2ahUKEwi06f-GpafqAhXLo4sKHVWeA3UQFjAAegQIBBAC&url=https%3A%2F%2Fdocs.docker.com%2Fcompose%2F&usg=AOvVaw02oes91geDSZ-H__u_XMxc))

``` bash
docker-compose up -d --no-deps --build
```

#### 2. Using docker build

``` bash
docker build -t egeback_playdownloader .
```

Both options will run swag, build golang code and deploy container

## Using API

Swagger documenation available at [http://localhost:8082/api/swagger/index.html](http://localhost:8082/api/swagger/index.html)

## TODO

* [x] Basic Authentication
* [x] Fix swag in docker
* [ ] Should svtplay-dl be removed or move to python for wrapper to reduce size?
* [ ] Test cases
* [x] Update README.md with documentation
* [x] Update code documentation
* [x] Add delete to /jobs to stop job (delete with normal clean up)
