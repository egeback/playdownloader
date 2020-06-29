# Play downloader

Microservice to download play serivice media files (series, news, movies etc.). The service wraps arround svtplay-dl, https://svtplay-dl.se/


## Deployment options

The microservice can be deployed as standlone application or in a docker container

### Standalone golang application

#### Clone repo

``` bash
git clone https://github.com/egeback/playdownloader.git
```

Run build script in from root director

``` bash
./cmd/build.sh
```

#### Update configuration

Update ./config/config.yaml

#### Run application

``` bash
./mediadownloader
```

### Docker container

#### Clone repo

``` bash
git clone https://github.com/egeback/playdownloader.git
```

#### Update configuration

1. Update ./config/config.yaml
2. Update Dockefile (update ports and location of media)

#### Using docker-compose ([link](https://www.google.com/url?sa=t&rct=j&q=&esrc=s&source=web&cd=&cad=rja&uact=8&ved=2ahUKEwi06f-GpafqAhXLo4sKHVWeA3UQFjAAegQIBBAC&url=https%3A%2F%2Fdocs.docker.com%2Fcompose%2F&usg=AOvVaw02oes91geDSZ-H__u_XMxc))


``` bash
docker-compose up -d --no-deps --build
```

This will run swag, build golang code and deploy container

## Using API

Swagger documenation available at http://localhost:8082/api/swagger/index.html

## TODO

* [x] Basic Authentication
* [x] Fix swag in docker
* [ ] Should svtplay-dl be removed or move to python for wrapper to reduce size?
* [ ] Test cases
* [x] Update README.md with documentation
* [x] Update code documentation
* [x] Add delete to /jobs to stop job (delete with normal clean up)
