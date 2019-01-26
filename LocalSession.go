package GoMybatis

import (
	"database/sql"
	"github.com/zhuxiujia/GoMybatis/utils"
)

//本地直连session
type LocalSession struct {
	SessionId              string
	db                     *sql.DB
	stmt                   *sql.Stmt
	tx                     *sql.Tx
	isCommitedOrRollbacked bool
	isClosed               bool
}

func (this *LocalSession) Id() string {
	return this.SessionId
}

func (this *LocalSession) Rollback() error {
	if this.isClosed == true {
		return utils.NewError("LocalSession", " can not Rollback() a Closed Session!")
	}
	if this.tx != nil {
		var err = this.tx.Rollback()
		if err == nil {
			this.isCommitedOrRollbacked = true
		} else {
			return err
		}
	}
	return nil
}

func (this *LocalSession) Commit() error {
	if this.isClosed == true {
		return utils.NewError("LocalSession", " can not Commit() a Closed Session!")
	}
	if this.tx != nil {
		var err = this.tx.Commit()
		if err == nil {
			this.isCommitedOrRollbacked = true
		}
	}
	return nil
}

func (this *LocalSession) Begin() error {
	if this.isClosed == true {
		return utils.NewError("LocalSession", " can not Begin() a Closed Session!")
	}
	if this.tx == nil {
		var tx, err = this.db.Begin()
		if err == nil {
			this.tx = tx
		} else {
			return err
		}
	}
	return nil
}

func (this *LocalSession) Close() {
	if this.db != nil {
		if this.stmt != nil {
			this.stmt.Close()
		}
		// When Close be called, if session is a transaction and do not call
		// Commit or Rollback, then call Rollback.
		if this.tx != nil && !this.isCommitedOrRollbacked {
			this.tx.Rollback()
		}
		this.tx = nil
		this.db = nil
		this.stmt = nil
		this.isClosed = true
	}
}

func (this *LocalSession) Query(sqlorArgs string) ([]map[string][]byte, error) {
	if this.isClosed == true {
		return nil, utils.NewError("LocalSession", " can not Query() a Closed Session!")
	}
	var rows *sql.Rows
	var err error
	if this.tx != nil {
		rows, err = this.tx.Query(sqlorArgs)
	} else {
		rows, err = this.db.Query(sqlorArgs)
	}
	if err != nil {
		return nil, err
	} else {
		defer rows.Close()
		return rows2maps(rows)
	}
	return nil, nil
}

func (this *LocalSession) Exec(sqlorArgs string) (*Result, error) {
	if this.isClosed == true {
		return nil, utils.NewError("LocalSession", " can not Exec() a Closed Session!")
	}
	var result sql.Result
	var err error
	if this.tx != nil {
		if this.isCommitedOrRollbacked {
			return nil, utils.NewError("LocalSession", " Exec() sql fail!, session isCommitedOrRollbacked!")
		}
		result, err = this.tx.Exec(sqlorArgs)
	} else {
		result, err = this.db.Exec(sqlorArgs)
	}
	if err != nil {
		return nil, err
	} else {
		var LastInsertId, _ = result.LastInsertId()
		var RowsAffected, _ = result.RowsAffected()
		return &Result{
			LastInsertId: LastInsertId,
			RowsAffected: RowsAffected,
		}, nil
	}
}
