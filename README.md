# mangonel
Overly simplified image upload HTTP microservice

> A mangonel was a type of catapult or siege engine used in the medieval period to throw projectiles at a castle's walls.
> [Wikipedia](https://en.wikipedia.org/wiki/Mangonel)

> Throw images at web server, hopefuly it'll stick.

## Motivation

For some reason [Lutim](https://github.com/ldidry/lutim) in non Perl equivalents are insanely hard to find.

This project does not aim to be as feature complete as [Pictshare](https://github.com/HaschekSolutions/pictshare), but to offer a simple web UI and some kind of backend storage.

## Getting started
### Prerequisites
* [go](https://golang.org/doc/install)
* Some storage

### Installing

    git clone https://github.com/jf-guillou/mangonel.git
    cd mangonel
    go get
    go build

### Updating

    cd mangonel
    git pull
    go get
    go build

### Configuration

Using environment variables (these values are default) :

    Mangonel_HashLength=5
    Mangonel_ListenAddr="127.0.0.1:8066"
    Mangonel_MaxFileSize=10240000
    Mangonel_StoragePath="./storage"

### Run

    ./mangonel

## FAQ

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details
