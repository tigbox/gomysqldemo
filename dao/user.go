package dao

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/chenchenyu/gomysqldemo/model"
)

const (
	userTableName = "user"
)

func (dao *MysqlDao) UserInsert(ctx context.Context, u *model.User) error {
	sqlStr := genInsertSQL(ctx, userTableName, u)
	result, err := dao.db.NamedExec(sqlStr, u)
	if err != nil {
		return err
	}
	u.ID, _ = result.LastInsertId()
	return err
}

func (dao *MysqlDao) UserUpdateByID(ctx context.Context, u *model.User) error {
	dao.updateByCond(ctx, userTableName, u, sq.Eq{"id": u.ID})
	return nil
}

func (dao *MysqlDao) UserDeleteByID(ctx context.Context, u *model.User) error {
	deleteCond := sq.Eq{"id": u.ID}
	deleteSql, args, err := sq.Delete(userTableName).Where(deleteCond).ToSql()
	if err != nil {
		return err
	}
	_, err = dao.db.Exec(deleteSql, args...)
	return err
}

func (dao *MysqlDao) UsersSelectByName(ctx context.Context, names []string) ([]*model.User, error) {
	selectExp := []string{"*"}
	cond := sq.Or{}
	for _, name := range names {
		cond = append(cond, sq.Eq{"name": name})
	}
	users := make([]*model.User, 0)
	dao.selectByCond(ctx, userTableName, selectExp, nil, nil, cond, &users)
	return users, nil
}
