package index

import (
	"database/sql"
	"errors"
	"fmt"
	"gitlab-cache/pkg/database"
)

var (
	CantInitializeIndexError         = errors.New("can't initialize index")
	CantCheckPartitionExistenceError = errors.New("can't check partition existence")
	CantAddFileToPartitionError      = errors.New("can't add file to partition")
	CantAddPartitionError            = errors.New("can't add partition")
)

type Index interface {
	AddPartition(uuid string) (int64, error)
	AddFileToPartition(key string, partitionId int64, begin int64, size int64) error
	PartitionExists(uuid string) (int64, bool, error)
}

type index struct {
	database                    database.Database
	addPartitionStatement       *sql.Stmt
	addFileToPartitionStatement *sql.Stmt
	partitionExistsStatement    *sql.Stmt
	getPartitionIdStatement     *sql.Stmt
}

func NewIndex(database database.Database) *index {
	return &index{database: database}
}

func (i *index) Init() (err error) {
	if i.addPartitionStatement, err = i.database.Statement("insert into partition (uuid, time) values ($1, $2)"); err != nil {
		return fmt.Errorf("%w. %s", CantInitializeIndexError, err)
	}
	if i.addFileToPartitionStatement, err = i.database.Statement("insert into file (key, partition_id, begin, size) values ($1, $2, $3, $4)"); err != nil {
		return fmt.Errorf("%w. %s", CantInitializeIndexError, err)
	}
	if i.partitionExistsStatement, err = i.database.Statement("select uuid from partition where uuid = $1"); err != nil {
		return fmt.Errorf("%w. %s", CantInitializeIndexError, err)
	}
	if i.getPartitionIdStatement, err = i.database.Statement("select id from partition where uuid = $1"); err != nil {
		return fmt.Errorf("%w. %s", CantInitializeIndexError, err)
	}
	return nil
}

func (i *index) AddPartition(uuid string) (int64, error) {
	result, err := i.addPartitionStatement.Exec(uuid)
	if err != nil {
		return 0, fmt.Errorf("%w. [uuid: %s]. %s", CantAddPartitionError, uuid, err)
	}
	return result.LastInsertId()
}

func (i *index) AddFileToPartition(key string, partitionId int64, begin int64, size int64) error {
	_, err := i.addFileToPartitionStatement.Exec(key, partitionId, begin, size)
	if err != nil {
		return fmt.Errorf("%w. [key: %s, partition: %d, begin: %d, size: %d]. %s", CantAddFileToPartitionError, key, partitionId, begin, size, err)
	}
	return nil
}

func (i *index) PartitionExists(uuid string) (int64, bool, error) {
	row := i.getPartitionIdStatement.QueryRow(uuid)
	var id int64
	err := row.Scan(&id)
	if err == nil {
		return id, true, nil
	} else if err == sql.ErrNoRows {
		return 0, false, nil
	} else {
		return 0, false, fmt.Errorf("%w. [uuid: %s]. %s", CantCheckPartitionExistenceError, uuid, err)
	}
}
