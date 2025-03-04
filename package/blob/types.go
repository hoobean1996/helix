package blob

type BlobSpecification struct {
	ProjectName string
	Env         string
	AccessKey   string
	SecretKey   string

	RequestCPU    string
	RequestMemory string
	LimitCPU      string
	LimitMemory   string
}
