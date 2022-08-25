package index

type Transaction interface {
	AddPartition(uuid string) (int64, error)
	AddFileToPartition(key string, partitionId int64, begin int64, size int64) error
	Commit() error
}

type Index interface {
	CreateTransaction() (Transaction, error)
	PartitionExists(uuid string) (int64, bool)
}

type index struct {
}
