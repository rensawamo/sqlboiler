package main

// go:generate sqlboiler --wipe mysql

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/oklog/ulid/v2"
	"github.com/takekazuomi/sqlboiler01/poc6/gen/models"
	"github.com/volatiletech/sqlboiler/v4/boil"

	// Dot import so we can access query mods directly instead of prefixing with "qm."
	. "github.com/volatiletech/sqlboiler/v4/queries/qm"
)


// https://github.com/volatiletech/sqlboiler?tab=readme-ov-file#query-mod-system


func main() {
	boil.DebugMode = true
	ctx := context.Background()

	dsn, ok := os.LookupEnv("DSN")
	dieFalse(ok, "env DSN not found")

	db, err := sql.Open("mysql", dsn)
	dieIf(err)

	err = db.Ping()
	dieIf(err)

	fmt.Println("connected")


	// 生成された構造体（この例では Pilot）に対してraw queryを使用します。
	// このクエリ mod が呼び出しに存在する場合、他のクエリを上書きします。
	// ここでは "?" プレースホルダーはサポートされていません。"$1, $2" などを使用してください。
	// 注: mysql だと、? になる

	q := SQL("select * from pilots where id=?", 1)
	_ = q


	ps, err := models.Pilots(SQL("select * from pilots where id=?", 1)).All(ctx, db)
	dieIf(err)
	printPilots(ps)

	
	// 以下のようなクエリを生成することができる
	// select 節の構築
	_ = Select("id", "name") // Select specific columns.
	_ = Select(models.PilotColumns.ID, models.PilotColumns.Name)
	_ = From("pilots as p") // Specify the FROM table manually, can be useful for doing complex queries.
	_ = From(models.TableNames.Pilots + " as p")

	// where 節の構築
	Where("name=?", "John")
	models.PilotWhere.Name.EQ("John")

	And("age=?", 24)
	// No equivalent type safe query yet
	Or("height=?", 183)

	// No equivalent type safe query yet
	Where("(name=? and age=?) or (age=?)", "John", 5, 6)

	// Expr allows manual grouping of statements
	// Exprでステートメントを手動でグループ化することが可能
	// Where(
	// 	Expr(
	// 		models.PilotWhere.Name.EQ("John"),
	// 		Or2(models.PilotWhere.Age.EQ(5)),
	// 	),
	// 	Or2(models.PilotWhere.Age.EQ(6)),
	// )

	// Whereの代わりにExprを使うと、ステートメントを手動でグループ化することができる
	qms := []QueryMod{
		Expr(
			models.PilotWhere.Name.EQ("John"),
			Or2(models.PilotWhere.Age.EQ(5)),
		),
		Or2(models.PilotWhere.Age.EQ(6)),
	}
	ps, err = models.Pilots(qms...).All(ctx, db)

	dieIf(err)
	printPilots(ps)

	// WHERE IN clause building いずれかに一致するものを取得する
	WhereIn("name, age in ?", "John", 24, "Tim", 33) // Generates: WHERE ("name","age") IN (($1,$2),($3,$4))
	//WhereIn(fmt.Sprintf("%s, %s in ?", models.PilotColumns.Name, models.PilotColumns.Age, "John", 24, "Tim", 33))
	AndIn("weight in ?", 84)
	//AndIn(models.PilotColumns.Weight+" in ?", 84)
	OrIn("height in ?", 183, 177, 204)
	//OrIn(models.PilotColumns.Height+" in ?", 183, 177, 204)

	InnerJoin("pilots p on jets.pilot_id=?", 10)
	InnerJoin(models.TableNames.Pilots+" p on "+models.TableNames.Jets+"."+models.JetColumns.PilotID+"=?", 10)

	GroupBy("name")
	//GroupBy("name like ? DESC, name", "John") // MySQLだとできなけど、ポスグレだとできるのかな？
	GroupBy(models.PilotColumns.Name)
	OrderBy("age, height")
	//OrderBy(models.PilotColumns.Age, models.PilotColumns.Height)

	// group by having
	Having("count(jets) > 2")
	Having(fmt.Sprintf("count(%s) > 2", models.TableNames.Jets)) 

	Limit(15)
	Offset(5) // Offset(5) は、5行目から取得する

	// Explicit locking
	For("update nowait")

	// Common Table Expressions
	With("cte_0 AS (SELECT * FROM table_0 WHERE thing=$1 AND stuff=$2)")


	//  条件を指定して取得する
	Load("Languages", models.LanguageWhere.Language.EQ("english")) // If it's a ToOne relationship it's in singular form, ToMany is plural.
	Load(models.PilotRels.Languages, models.LanguageWhere.Language.EQ("english"))

	// Relationships に書いてあった
	// https://github.com/volatiletech/sqlboiler#relationships

	// 遅いパターン
	// 飛行機を操縦できるパイロットを全部探す。これだと、
	// `SELECT `pilots`.* FROM `pilots` WHERE (`id` = ?) LIMIT 1;` が飛行機の数だけ実行される。
	// Avoid this loop query pattern, it is slow.
	jets, _ := models.Jets().All(ctx, db)
	pilots := make([]*models.Pilot, len(jets))
	for i := 0; i < len(jets); i++ {
		pilots[i], _ = jets[i].Pilot().One(ctx, db)
	}

	// Instead, use Eager Loading!
	jets, _ = models.Jets(Load("Pilot")).All(ctx, db)

	for _, j := range jets {
		_ = j.R.Pilot
		fmt.Println(j.Name, ":", j.R.Pilot.Name)
	}

	// Type safe relationship names exist too:
	jets, _ = models.Jets(Load(models.JetRels.Pilot)).All(ctx, db)
	for _, j := range jets {
		_ = j.R.Pilot
		fmt.Println(j.Name, ":", j.R.Pilot.Name)
	}

	// こんなことをすると
	//_, err = models.Jets(Load(models.JetRels.Pilot), models.PilotWhere.Name.EQ("taro yamada")).All(ctx, db)
	//dieIf(err)
	// `SELECT `jets`.* FROM `jets` WHERE (`pilots`.`name` = ?);`が生成されて、下記のエラーになる。
	// panic: models: failed to assign all query results to Jet slice: bind failed to execute query: Error 1054: Unknown column 'pilots.name' in 'where clause'

	// Example of a nested load.
	// Each jet will have its pilot loaded, and each pilot will have its languages loaded.
	_, _ = models.Jets(Load("Pilot.Languages")).All(ctx, db)
	// Note that each level of a nested Load call will be loaded. No need to call Load() multiple times.

	// こんなSQLになった。`Pilot.Languages`と指定されたので、飛行機を持ってきて、飛行機に紐づくパイロットを持ってきて、言語とパイロットはjoinで引っ掛けてもってくる。
	// ここで、パイロットと言語のJOINを選択したのは、pilot_languages が公差エンティに入っているかららしい。
	// SELECT `jets`.* FROM `jets`;
	// []
	// SELECT * FROM `pilots` WHERE (`pilots`.`id` IN (?,?));
	// [1 2]
	// SELECT `languages`.`id`, `languages`.`language`, `a`.`pilot_id` FROM `languages` INNER JOIN `pilot_languages` as `a` on `languages`.`id` = `a`.`language_id` WHERE (`a`.`pilot_id` IN (?,?));
	// [1 2]

	//英語を喋るパイロットを取得する
	p := make([]*models.Pilot, 0, 100)
	qm := SQL(`select * from pilots p
	join pilot_languages pl on p.id = pl.pilot_id
	join languages l on l.id = pl.language_id
	where  l.language = ?`, "english")
	err = models.Pilots(qm).Bind(ctx, db, &p) //Bind the results of a query to your own struct object.
	dieIf(err)

	printPilots(p)
}

func dieFalse(ok bool, msg string) {
	if !ok {
		panic(msg)
	}
}

func dieIf(err error) {
	if err != nil {
		panic(err)
	}
}

func newULID() ulid.ULID {
	return ulid.MustNew(ulid.Now(), ulid.DefaultEntropy())
}

func printPilots(ps models.PilotSlice) {
	for i, p := range ps {
		fmt.Println(i, ":", p)
	}
}