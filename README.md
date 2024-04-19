作成手順

## .envrc の導入 (インフラと分岐するときは gitignore)
```sh
$ direnv edit .
```

## mysql のセットアップ
```sh
# dbconfig.ymlの 以下の出力先を変更する
dir: boiler2/migrations 
```

```sh
$ sql-migrate new MIGRATEFILENAME
```

```sh
$ make up
docker compose -f database/docker-compose.yml up --force-recreate -d
[+] Running 12/12
 ✔ db 11 layers [⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿]      0B/0B      Pulled                                                                  22.3s 
   ✔ e54b73e95ef3 Pull complete                                                                                        3.5s 
   .....
```

![alt text](image.png)

## sql boiler のセットアップ

### boiler に移動
```sh
$ cd boiler
```

```sh
$ make buf-gen
```

###  [sql-migrate](https://github.com/rubenv/sql-migrate)

```sh
$ make migrate-up
```

### boiler の mysql のコード生成
```sh
$ make boiler-gen
```

### コマンド実行
```sh
$ go run main.go
```

### boiler utility

https://github.com/volatiletech/sqlboiler?tab=readme-ov-file#query-mod-system
```sh
// Dot import so we can access query mods directly instead of prefixing with "qm."
import . "github.com/volatiletech/sqlboiler/v4/queries/qm"

// Use a raw query against a generated struct (Pilot in this example)
// If this query mod exists in your call, it will override the others.
// "?" placeholders are not supported here, use "$1, $2" etc.
SQL("select * from pilots where id=$1", 10)
models.Pilots(SQL("select * from pilots where id=$1", 10)).All()

Select("id", "name") // Select specific columns.
Select(models.PilotColumns.ID, models.PilotColumns.Name)
From("pilots as p") // Specify the FROM table manually, can be useful for doing complex queries.
From(models.TableNames.Pilots + " as p")

// WHERE clause building
Where("name=?", "John")
models.PilotWhere.Name.EQ("John")
And("age=?", 24)
// No equivalent type safe query yet
Or("height=?", 183)
// No equivalent type safe query yet

Where("(name=? and age=?) or (age=?)", "John", 5, 6)
// Expr allows manual grouping of statements
Where(
  Expr(
    models.PilotWhere.Name.EQ("John"),
    Or2(models.PilotWhere.Age.EQ(5)),
  ),
  Or2(models.PilotAge),
)

// WHERE IN clause building
WhereIn("(name, age) in ?", "John", 24, "Tim", 33) // Generates: WHERE ("name","age") IN (($1,$2),($3,$4))
WhereIn(fmt.Sprintf("(%s, %s) in ?", models.PilotColumns.Name, models.PilotColumns.Age), "John", 24, "Tim", 33)
AndIn("weight in ?", 84)
AndIn(models.PilotColumns.Weight + " in ?", 84)
OrIn("height in ?", 183, 177, 204)
OrIn(models.PilotColumns.Height + " in ?", 183, 177, 204)

InnerJoin("pilots p on jets.pilot_id=?", 10)
InnerJoin(models.TableNames.Pilots + " p on " + models.TableNames.Jets + "." + models.JetColumns.PilotID + "=?", 10)

GroupBy("name")
GroupBy("name like ? DESC, name", "John")
GroupBy(models.PilotColumns.Name)
OrderBy("age, height")
OrderBy(models.PilotColumns.Age, models.PilotColumns.Height)

Having("count(jets) > 2")
Having(fmt.Sprintf("count(%s) > 2", models.TableNames.Jets)

Limit(15)
Offset(5)

// Explicit locking
For("update nowait")

// Common Table Expressions
With("cte_0 AS (SELECT * FROM table_0 WHERE thing=$1 AND stuff=$2)")

// Eager Loading -- Load takes the relationship name, ie the struct field name of the
// Relationship struct field you want to load. Optionally also takes query mods to filter on that query.
Load("Languages", Where(...)) // If it's a ToOne relationship it's in singular form, ToMany is plural.
Load(models.PilotRels.Languages, Where(...))
```