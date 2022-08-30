package index

import (
	"artifacts-cache/pkg/repository/database"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

var (
	CantInitializeIndexError         = errors.New("can't initialize index")
	CantCheckPartitionExistenceError = errors.New("can't check partition existence")
	CantFindContentEmplacementError  = errors.New("can't find content emplacement")
	CantAddFileToPartitionError      = errors.New("can't add file to partition")
	CantAddPartitionError            = errors.New("can't add partition")
	CantGetPartitionIdError          = errors.New("can't get partition id")
)

type Index interface {
	AddPartition(uuid string) (int64, error)
	AddFileToPartition(partitionId int64, subset string, path string, offset int64, size int64) error
	PartitionExists(uuid string) (int64, bool, error)
	FindContentEmplacement(subset string, filter string) (*ContentEmplacement, error)
}

type index struct {
	database                        database.Database
	addPartitionStatement           *sql.Stmt
	addFileToPartitionStatement     *sql.Stmt
	partitionExistsStatement        *sql.Stmt
	getPartitionIdStatement         *sql.Stmt
	findContentEmplacementStatement *sql.Stmt
}

func NewIndex(database database.Database) *index {
	return &index{database: database}
}

func (i *index) Init() (err error) {
	if i.addPartitionStatement, err = i.database.Statement("insert into partition (uuid, time) values ($1, $2)"); err != nil {
		return fmt.Errorf("%w. %s", CantInitializeIndexError, err)
	}
	if i.addFileToPartitionStatement, err = i.database.Statement("insert into file (partition_id, subset, path, \"offset\", size) values ($1, $2, $3, $4, $5)"); err != nil {
		return fmt.Errorf("%w. %s", CantInitializeIndexError, err)
	}
	if i.partitionExistsStatement, err = i.database.Statement("select uuid from partition where uuid = $1"); err != nil {
		return fmt.Errorf("%w. %s", CantInitializeIndexError, err)
	}
	if i.getPartitionIdStatement, err = i.database.Statement("select id from partition where uuid = $1"); err != nil {
		return fmt.Errorf("%w. %s", CantInitializeIndexError, err)
	}
	if i.findContentEmplacementStatement, err = i.database.Statement("select p.uuid, f.path, f.offset, f.size from partition p join file f on p.id = f.partition_id where f.subset = $1 and f.path like $2"); err != nil {
		return fmt.Errorf("%w. %s", CantInitializeIndexError, err)
	}
	return nil
}

func (i *index) AddPartition(uuid string) (int64, error) {
	_, err := i.addPartitionStatement.Exec(uuid, time.Now())
	if err != nil {
		return 0, fmt.Errorf("%w. [uuid: %s]. %s", CantAddPartitionError, uuid, err)
	}
	id, err := i.GetPartitionId(uuid)
	if err != nil {
		return 0, fmt.Errorf("%w. [uuid: %s]. %s", CantGetPartitionIdError, uuid, err)
	}
	return id, nil
}

func (i *index) AddFileToPartition(partitionId int64, subset string, path string, offset int64, size int64) error {
	_, err := i.addFileToPartitionStatement.Exec(partitionId, subset, path, offset, size)
	if err != nil {
		return fmt.Errorf("%w. [partition: %d, path: %s, offset: %d, size: %d]. %s", CantAddFileToPartitionError, partitionId, path, offset, size, err)
	}
	return nil
}

func (i *index) PartitionExists(uuid string) (int64, bool, error) {
	id, err := i.GetPartitionId(uuid)
	if err == nil {
		return id, true, nil
	} else if err == sql.ErrNoRows {
		return 0, false, nil
	} else {
		return 0, false, fmt.Errorf("%w. [uuid: %s]. %s", CantCheckPartitionExistenceError, uuid, err)
	}
}

func (i *index) FindContentEmplacement(subset string, filter string) (*ContentEmplacement, error) {
	rows, err := i.findContentEmplacementStatement.Query(subset, filter)
	if err != nil {
		return nil, fmt.Errorf("%w. [subset: %s, filter: %s]. %s", CantFindContentEmplacementError, subset, filter, err)
	}
	defer rows.Close()
	contentEmplacement := &ContentEmplacement{}
	for rows.Next() {
		emplacement := &Emplacement{}
		if err := rows.Scan(&emplacement.Partition, &emplacement.Path, &emplacement.Offset, &emplacement.Size); err != nil {
			return contentEmplacement, fmt.Errorf("%w. [subset: %s, filter: %s]. %s", CantFindContentEmplacementError, subset, filter, err)
		}
		contentEmplacement.Emplacements = append(contentEmplacement.Emplacements, emplacement)
	}
	if err = rows.Err(); err != nil {
		return contentEmplacement, fmt.Errorf("%w. [subset: %s, filter: %s]. %s", CantFindContentEmplacementError, subset, filter, err)
	}
	return contentEmplacement, nil
}

func (i *index) GetPartitionId(uuid string) (int64, error) {
	row := i.getPartitionIdStatement.QueryRow(uuid)
	var id int64
	err := row.Scan(&id)
	return id, err
}
