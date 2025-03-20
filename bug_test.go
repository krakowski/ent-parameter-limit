package bug

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"testing"

	"entgo.io/ent/dialect"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/samber/lo"

	"entgo.io/bug/ent"
	"entgo.io/bug/ent/enttest"
)

func TestBugSQLite(t *testing.T) {
	client := enttest.Open(t, dialect.SQLite, "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()
	test(t, client)
}

func TestBugMySQL(t *testing.T) {
	for version, port := range map[string]int{"56": 3306, "57": 3307, "8": 3308} {
		addr := net.JoinHostPort("localhost", strconv.Itoa(port))
		t.Run(version, func(t *testing.T) {
			client := enttest.Open(t, dialect.MySQL, fmt.Sprintf("root:pass@tcp(%s)/test?parseTime=True", addr))
			defer client.Close()
			test(t, client)
		})
	}
}

func TestBugPostgres(t *testing.T) {
	for version, port := range map[string]int{"10": 5430, "11": 5431, "12": 5432, "13": 5433, "14": 5434} {
		t.Run(version, func(t *testing.T) {
			client := enttest.Open(t, dialect.Postgres, fmt.Sprintf("host=localhost port=%d user=postgres dbname=test password=pass sslmode=disable", port))
			defer client.Close()
			test(t, client)
		})
	}
}

func TestBugMaria(t *testing.T) {
	for version, port := range map[string]int{"10.5": 4306, "10.2": 4307, "10.3": 4308} {
		t.Run(version, func(t *testing.T) {
			addr := net.JoinHostPort("localhost", strconv.Itoa(port))
			client := enttest.Open(t, dialect.MySQL, fmt.Sprintf("root:pass@tcp(%s)/test?parseTime=True", addr))
			defer client.Close()
			test(t, client)
		})
	}
}

func cleanup(ctx context.Context, client *ent.Client) error {
	if _, err := client.Profile.Delete().Exec(ctx); err != nil {
		return err
	}

	if _, err := client.User.Delete().Exec(ctx); err != nil {
		return err
	}

	return nil
}

const parameterLimit = 65536
const chunkSize = 8192

func test(t *testing.T, client *ent.Client) {
	ctx := context.Background()

	// Clean up old data
	if err := cleanup(ctx, client); err != nil {
		t.Error("creating users failed", err)
	}

	// Create builders for all users
	users := lo.Times(parameterLimit, func(index int) *ent.UserCreate {
		return client.User.Create().
			SetID(index + 1).
			SetUsername(fmt.Sprintf("user-%d", index))
	})

	// Create users using a bulk insert operation
	for _, chunk := range lo.Chunk(users, chunkSize) {
		if err := client.User.CreateBulk(chunk...).Exec(ctx); err != nil {
			t.Error("creating users failed", err)
		}
	}

	// Create builders for all profiles
	profiles := lo.Times(parameterLimit, func(index int) *ent.ProfileCreate {
		return client.Profile.Create().
			SetID(index + 1).
			SetUserID(index + 1).
			SetFirstname(fmt.Sprintf("firstname-%d", index)).
			SetLastname(fmt.Sprintf("lastname-%d", index))
	})

	// Create profiles using a bulk insert operation
	for _, chunk := range lo.Chunk(profiles, chunkSize) {
		if err := client.Profile.CreateBulk(chunk...).Exec(ctx); err != nil {
			t.Error("creating profiles failed", err)
		}
	}

	// Query all users with profiles. This will fail as we have more than 65535 parameters.
	_, err := client.Debug().User.Query().
		WithProfile().
		All(ctx)

	// Check if there was an error
	if err != nil {
		t.Error("querying users failed", err)
	}
}
