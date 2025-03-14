# Commands

### Migrate

```bash
migrate create -ext=sql -dir=migrations -seq <name>
```

#### Run

```bash
migrate -path=migrations -database "mysql://root:root@tcp(localhost:3306)/bank" up

or

migrate -path=migrations -database "mysql://root:root@tcp(localhost:3306)/bank" down
```
