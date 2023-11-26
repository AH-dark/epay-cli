## epay cli
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FAH-dark%2Fepay-cli.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2FAH-dark%2Fepay-cli?ref=badge_shield)


A command line tool for epay testing.

### Build

```bash
make build
```

### Usage

Based on previous build, you can find the binary file in `bin/epay-cli`.

We provide a simple command line tool for epay testing and maintenance.

#### Test for API

```bash
./bin/epay-cli test submit # test submit api
./bin/epay-cli test mapi # test mapi submit api
```

#### Migration

```bash
./bin/epay-cli migrate # generate sql schema
```

#### Help

```bash
./bin/epay-cli help
```


## License
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FAH-dark%2Fepay-cli.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2FAH-dark%2Fepay-cli?ref=badge_large)