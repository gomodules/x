package types_test

import (
	"fmt"

	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq"
	_ "gocloud.dev/secrets/gcpkms"
	"github.com/appscode/go/encoding/sql/types"
)

// User describes a user
type Credential struct {
	Id   int64
	Name string
	Data types.SecureString `xorm:"text"`
}

func main() {
	Orm, err := newPGEngine("postgres", "postgres", "127.0.0.1", 5432, "postgres")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = Orm.CreateTables(&Credential{})
	if err != nil {
		fmt.Println(err)
		return
	}

	url := fmt.Sprintf("gcpkms://projects/%v/locations/%v/keyRings/%v/cryptoKeys/%v", "ackube", "global", "gitea", "gitea-key")
	fmt.Println(url)
	_, err = Orm.Insert(&Credential{1, "test", types.SecureString{
		Url:  url,
		Data: "this is a test credential",
	}})
	if err != nil {
		fmt.Println(err)
		return
	}

	creds := make([]Credential, 0)
	err = Orm.Find(&creds)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(creds)
}

// Connects to any databse using provided credentials
func newPGEngine(user, password, host string, port int64, dbName string) (*xorm.Engine, error) {
	cnnstr := fmt.Sprintf("user=%v password=%v host=%v port=%v dbname=%v sslmode=disable",
		user, password, host, port, dbName)
	engine, err := xorm.NewEngine("postgres", cnnstr)
	if err != nil {
		return nil, err
	}
	// engine.ShowSQL(system.Env() == system.DevEnvironment)
	engine.ShowSQL(true)
	return engine, nil
}
