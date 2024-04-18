package main

//go:generate sqlboiler --wipe mysql

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/oklog/ulid/v2"
	"github.com/rensawamo/sqlboiler/boiler/gen/models"
	pb "github.com/rensawamo/sqlboiler/boiler/gen/buf/message"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func main() {
	boil.DebugMode = true

	// mysql 接続
	dsn, ok := os.LookupEnv("DSN")
	dieFalse(ok, "env DSN not found")

	db, err := sql.Open("mysql", dsn)
	dieIf(err)

	err = db.Ping()
	dieIf(err)

	fmt.Println("connected")

	// トランザクション
	tx, err := db.BeginTx(context.Background(), nil)
	dieIf(err)
	defer tx.Rollback()

	now := time.Now()
	id := newULID()
	// migrateされたdbから model が sqlboilerに	よって生成される

	r := &models.Resource{
		ID:        id,
    // from status enum('STATUS_ACTIVE', 'STATUS_DELETED') not null, 
		Status:    models.ResourcesStatusSTATUS_ACTIVE,
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: null.Time{},
	}
	err = r.Insert(context.Background(), tx, boil.Infer())
	dieIf(err)

	p := models.Property{
		ResourceID: id,
		FirstName:  "fleelen",
		LastName:   "sousouno",
	}
	// db に insert
	err = p.Insert(context.Background(), tx, boil.Infer())
	dieIf(err)
	tx.Commit()

	fmt.Println("insert resources:", id)

	id = newULID()
	now = time.Now()
	pbr := &pb.Resource{
		Name: id.String(),
		Properties: &pb.Properties{
			FirstName: "hanako",
			LastName:  "yamada",
		},
		Status:     pb.Status_STATUS_ACTIVE,
		CreateTime: timestamppb.New(now),
		UpdateTime: timestamppb.New(now),
		DeleteTime: nil,
	}
	// スキーマと合わせるとdb のモデルが即座に生成可能
	// message Resource {
  //   // 名前
  //   string name = 1;
  //   // プロパティ
  //   Properties properties = 2;
  //   // 状態
  //   Status status = 3;
  //   // 作成
  //   google.protobuf.Timestamp create_time = 12;
  //   // 更新
  //   google.protobuf.Timestamp update_time = 13;
  //   // 削除
  //   google.protobuf.Timestamp delete_time = 14;
  // }
  
	// プロトコルバッファからモデルを生成してmockデータを挿入
	mr := NewModelResourceFromPb(pbr)
	mp := NewModelPropertyFromPb(id, pbr.Properties)

	tx, err = db.BeginTx(context.Background(), nil)
	dieIf(err)
	defer tx.Rollback()

	err = mr.Insert(context.Background(), tx, boil.Infer())
	dieIf(err)

	err = mp.Insert(context.Background(), tx, boil.Infer())
	dieIf(err)
	tx.Commit()

	fmt.Println("insert resources:", id)

}

func newULID() ulid.ULID {
	return ulid.MustNew(ulid.Now(), ulid.DefaultEntropy())
}

func newTimestamppbFromNullTime(t null.Time) *timestamppb.Timestamp {
	if t.Valid {
		return timestamppb.New(t.Time)
	}
	return nil
}

func newNullTimeFromTimestamppb(t *timestamppb.Timestamp) null.Time {
	if t == nil {
		return null.NewTime(time.Time{}, false)
	}
	return null.NewTime(t.AsTime(), true)
}

func NewPbPropertyFromModel(source *models.Property) *pb.Properties {
	return &pb.Properties{
		FirstName: source.FirstName,
		LastName:  source.LastName,
	}
}

func NewModelPropertyFromPb(parent ulid.ULID, source *pb.Properties) *models.Property {
	return &models.Property{
		ID:         0,
		ResourceID: parent,
		FirstName:  source.FirstName,
		LastName:   source.LastName,
	}
}

func NewModelResourceFromPb(source *pb.Resource) *models.Resource {
	return &models.Resource{
		ID:        ulid.MustParse(source.Name),
		Status:    models.ResourcesStatus(pb.Status_name[int32(source.Status)]),
		CreatedAt: source.CreateTime.AsTime(),
		UpdatedAt: source.UpdateTime.AsTime(),
		DeletedAt: newNullTimeFromTimestamppb(source.DeleteTime),
	}
}

func NewPbResourceFromModel(source *models.Resource) *pb.Resource {
	return &pb.Resource{
		Name: source.ID.String(),
		Properties: &pb.Properties{
			FirstName: source.R.Property.FirstName,
			LastName:  source.R.Property.LastName,
		},
		Status:     pb.Status(pb.Status_value[source.Status.String()]),
		CreateTime: timestamppb.New(source.CreatedAt),
		UpdateTime: timestamppb.New(source.UpdatedAt),
		DeleteTime: newTimestamppbFromNullTime(source.DeletedAt),
	}
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
