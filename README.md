# DockerHttpClient

A web interface for the [Docker engine api](https://docs.docker.com/engine/api/v1.30/)

## Start

Create the JWT signing keys:

```bash
./keys.sh
```

Create a folder `/certs` and place your SSL certificates or create self-signed certificates (this creates the folder as well):

```bash
./certs.sh
```


## TODO
- JWT
- add SemanticUI
- Improve layout for Images and Containers
- add a toolbar with menu
- add Logout functionality
- Complete the containers page
- Add links in the menu for variaous operations: remember that the URLS are passed to docker engine so they need to respect the specification. Therefore everything should be done with links and forms