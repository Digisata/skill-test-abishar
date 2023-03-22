# Skill test abishar
>
> A simple RESTFUL API that provides endpoints to batch insert data.

![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)

## Prerequisites

- `Go 1.19.3` or higher.
- [Migrate](https://github.com/golang-migrate/migrate).

## Installation

- Clone this repository.
- Duplicate [.env.example](.env.example) file to `.env` file and change the `variable value` with your desired value.
- Next, run below command to migrate the database.

```shell
migrate -database "postgres://{DB_USERNAME}:{DB_PASSWORD}}@{DB_HOST}:{DB_PORT}/{DB_NAME}?sslmode=disable" -path db/migrations up 
```

> **_NOTE:_**  Replace `{DB_USERNAME}, {DB_PASSWORD}, {DB_HOST}`, etc with the corresponding value in your `.env` file.

- Then run this command to synchronize all the dependencies.

```shell
go mod tidy
```

- Finally, run the project.

```shell
go run .
```

## Usage

Do the post request as below:

```shell
curl -X POST \
  'localhost:{PORT}/sales' \
  --header 'Accept: */*' \
  --header 'User-Agent: Thunder Client (https://www.thunderclient.com)' \
  --header 'Content-Type: application/json' \
  --data-raw '{
  "request_id": 12345,
  "data": [
    {
      "id": 0,
      "customer": "Vazquez Rowland",
      "quantity": 1,
      "price": 10,
      "timestamp": "2022-01-01 22:10:44"
    }
  ]
}'
```

> **_NOTE:_**  Replace the `{PORT}` with your defined `PORT` value in your `.env` file.

## Meta

Hanif Naufal – [@Digisata](https://twitter.com/Digisata) – [hnaufal123@gmail.com](mailto:hnaufal123@gmail.com)

Distributed under the MIT license. See [LICENSE](LICENSE.md) for more information.

## Contributing

1. Fork this repository.
2. Create your own branch (`git checkout -b fooBar`).
3. Commit your changes (`git commit -am 'Add some fooBar'`).
4. Push to the branch (`git push origin fooBar`).
5. Create your awesome Pull Request.
