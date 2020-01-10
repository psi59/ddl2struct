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

#### Result
```sh
#   ____ ____ __      ____ __     ____ ____ ____ _  _  ___ ____ 
#  (    (    (  )    (_  _/  \   / ___(_  _(  _ / )( \/ __(_  _)
#   ) D () D / (_/\    )((  O )  \___ \ )(  )   ) \/ ( (__  )(  
#  (____(____\____/   (__)\__/   (____/(__)(__\_\____/\___)(__)

===============================================================
TableName: ddl2struct
ColumnCount: 5
Save your times: About 10 seconds
===============================================================

type ddl2struct struct {
        Personid  int    `json:"personid" gorm:"Column:personid"`
        Lastname  string `json:"lastname" gorm:"Column:lastname"`
        Firstname string `json:"firstname" gorm:"Column:firstname"`
        Address   string `json:"address" gorm:"Column:address"`
        City      string `json:"city" gorm:"Column:city"`
}
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