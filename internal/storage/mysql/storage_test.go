package mysql

import (
	"context"
	"fmt"
	"testing"

	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/vench/cryptocompare/internal/entities"

	"github.com/stretchr/testify/require"
	"github.com/vench/cryptocompare/internal/config"

	"github.com/testcontainers/testcontainers-go"
)

type mysqlContainer struct {
	testcontainers.Container
	ConnectionString string
}

func setupMysql(ctx context.Context) (*mysqlContainer, error) {
	req := testcontainers.ContainerRequest{
		Image: "mysql:8.0.20",
		Env: map[string]string{
			"MYSQL_ROOT_PASSWORD": "admin",
			"MYSQL_DATABASE":      "test",
			"MYSQL_USER":          "username",
			"MYSQL_PASSWORD":      "password",
		},
		ExposedPorts: []string{"3306/tcp"},
		WaitingFor:   wait.ForListeningPort("3306/tcp"),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	ip, err := container.Host(ctx)
	if err != nil {
		return nil, err
	}

	mappedPort, err := container.MappedPort(ctx, "3306/tcp")
	if err != nil {
		return nil, err
	}

	connectionString := fmt.Sprintf("username:password@tcp(%s:%s)/test", ip, mappedPort.Port())

	return &mysqlContainer{Container: container, ConnectionString: connectionString}, nil
}

func TestStorage_StoreCurrency(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	mysqlC, err := setupMysql(ctx)
	require.NoError(t, err)

	defer t.Cleanup(func() {
		require.NoError(t, mysqlC.Terminate(ctx))
	})

	s, err := New(&config.Mysql{
		ConnectionString: mysqlC.ConnectionString,
	})
	require.NoError(t, err)

	// create table
	createTable := "create table currency (    " +
		"`key`      varchar(8)                          not null  primary key," +
		"value      json                                null," +
		"created_at timestamp default CURRENT_TIMESTAMP null," +
		"updated_at timestamp default CURRENT_TIMESTAMP null)"

	_, err = s.conn.Exec(createTable)
	require.NoError(t, err)

	val := &entities.Currency{
		FromSymbol: "ABC",
		ToSymbol:   "XYZ",

		PRICE:           1000.10,
		CHANGEPCT24HOUR: 1999.2,
		CHANGE24HOUR:    4000.1,
	}

	require.NoError(t, s.StoreCurrency(val))

	cur, err := s.GetCurrencyBy([]string{val.FromSymbol}, []string{val.ToSymbol})
	require.NoError(t, err)
	require.Equal(t, 1, len(cur))
	require.Equal(t, val, cur[0])
}
