package driver

type DbClient struct{}

func NewDBClient() *DbClient {
	return &DbClient{}
}
