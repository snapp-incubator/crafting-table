package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueryBuilder(t *testing.T) {
	t.Run("where by age limit 1", func(t *testing.T) {
		q := UserQueryBuilder().Limit(1).WhereAgeEq(10)
		query, err := q.SQL()
		assert.NoError(t, err)
		assert.Equal(t, "SELECT * FROM users WHERE age = ? LIMIT 1", query)
		assert.Equal(t, []interface{}{10}, q.(*__UserSQLQueryBuilder).whereArgs)
	})

    t.Run("where by name offset 10", func(t *testing.T) {
        q := UserQueryBuilder().Offset(10).WhereNameEq("Hello World")
		query, err := q.SQL()
		assert.NoError(t, err)
		assert.Equal(t, "SELECT * FROM users WHERE name = ? OFFSET 10", query)
		assert.Equal(t, []interface{}{"Hello World"}, q.(*__UserSQLQueryBuilder).whereArgs)
    })

    t.Run("where by age greater 10 order by id desc", func(t *testing.T) {
        q := UserQueryBuilder().OrderByDesc(UserColumns.ID).WhereAgeGT(10)
		query, err := q.SQL()
		assert.NoError(t, err)
		assert.Equal(t, "SELECT * FROM users WHERE age > ? ORDER BY id DESC", query)
		assert.Equal(t, []interface{}{10}, q.(*__UserSQLQueryBuilder).whereArgs)
    })

    t.Run("update where by age greater 10 set age 11", func(t *testing.T) {
        q := UserQueryBuilder().WhereAgeGT(10).SetAge(11)
		query, err := q.SQL()
		assert.NoError(t, err)
		assert.Equal(t, "UPDATE users SET age = ? WHERE age > ?", query)
		assert.Equal(t, []interface{}{10}, q.(*__UserSQLQueryBuilder).whereArgs)
		assert.Equal(t, []interface{}{11}, q.(*__UserSQLQueryBuilder).setArgs)
    })

    t.Run("delete where age greater 10", func(t *testing.T) {
        q := UserQueryBuilder().WhereAgeGT(10)
        q.(*__UserSQLQueryBuilder).mode = "delete"
		query, err := q.SQL()
		assert.NoError(t, err)
		assert.Equal(t, "DELETE FROM users WHERE age > ?", query)
		assert.Equal(t, []interface{}{10}, q.(*__UserSQLQueryBuilder).whereArgs)

    })
}
