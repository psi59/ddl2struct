# Welcome to ddl2struct 👋
![Version](https://img.shields.io/badge/version-0.0.1-blue.svg?cacheSeconds=2592000)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)


> Create golang struct from ddl

## Install

```sh
go get github.com/realsangil/ddl2struct
```

## Usage
```sh
ddl2struct [flags]
```

#### Flags
```sh
-c, --copy            copy to clipboard
-h, --help            help for ddl2struct
-i, --input string    sql file path
-o, --output string   output file path
```

#### Example
```sh
ddl2struct -i example.sql
```
#### example.sql
```sql
CREATE TABLE ddl2struct (
    PersonID int,
    LastName varchar(255),
    FirstName varchar(255),
    Address varchar(255),
    City varchar(255)
);
```

## Author

👤 **realsangil**

* Website: [blog.realsangil.net](https://blog.realsangil.net)
* Github: [@realsangil](https://github.com/realsangil)

## Show your support

Give a ⭐️ if this project helped you!


## 📝 License

Copyright © 2020 [realsangil](https://github.com/realsangil).

This project is [Apache License 2.0](https://opensource.org/licenses/Apache-2.0) licensed.

***
_This README was generated with ❤️ by [readme-md-generator](https://github.com/kefranabg/readme-md-generator)_