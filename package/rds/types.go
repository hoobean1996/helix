package rds

type MySQLSpecification struct {
	ProjectName string
	// test, staging, prod
	Env      string
	Password string
	Database string

	RequestCPU    string
	RequestMemory string
	LimitCPU      string
	LimitMemory   string
}
