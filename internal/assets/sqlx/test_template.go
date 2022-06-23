package sqlx

import (
	"fmt"

	"github.com/iancoleman/strcase"
	"github.com/snapp-incubator/crafting-table/internal/structure"
)

type SqlxTest interface {
	Insert(structure *structure.Structure) (syntax string)
	UpdateAll(structure *structure.Structure) (syntax string)
	UpdateBy(structure *structure.Structure, vars *[]structure.UpdateVariables) (syntax string)
	SelectAll(structure *structure.Structure) (syntax string)
	SelectBy(structure *structure.Structure, vars *[]structure.Variables) (syntax string)
}

type sqlxTest struct {
	insertTestSuccess string
	insertTestFailure string

	// TODO : add tests for SELECT All
	//selectAllTestSuccess string
	//selectAllTestFailure string

	selectByTestSuccess            string
	selectByTestErrorNoRowsFailure string
	selectByTestFailure            string

	updateAllTestSuccess string
	updateAllTestFailure string

	updateByTestSuccess string
	updateByTestFailure string
}

func NewSqlxTest() SqlxTest {
	s := sqlxTest{}

	s.insertTestSuccess = `
func (suite *%sRepositoryTestSuite) TestInsert_Success() {
	require := suite.Require()

	%s

	%s

	sqlmock.NewRows([]string{
		%s
	}).
		AddRow(
			%s
		)

	syntax := "INSERT INTO %s .+"
	suite.mock.ExpectExec(syntax).
		WithArgs(
			%s
		).
		WillReturnResult(sqlmock.NewResult(int64(1), expRowsAffected))

	err := suite.repo.Insert(
		context.Background(),
		&%s,
	)
	require.NoError(err)
	require.NoError(suite.mock.ExpectationsWereMet())
}

`

	s.insertTestFailure = `

func (suite *%sRepositoryTestSuite) TestInsert_Failure() {
	require := suite.Require()

	expectedError := errors.New("something went wrong")

	%s

	%s

	syntax := "INSERT INTO %s (.+) VALUES (.+)"
	suite.mock.ExpectExec(syntax).
		WithArgs(
			%s
		).
		WillReturnError(expectedError)

	err := suite.repo.Insert(
		context.Background(),
		&%s,
	)
	require.Equal(expectedError, err)
	require.NoError(suite.mock.ExpectationsWereMet())
}

`

	s.selectByTestSuccess = `
func (suite *%sRepositoryTestSuite) Test%s_Success() {
	require := suite.Require()

	%s

	%s

	rows := sqlmock.NewRows([]string{
		%s
	}).
		AddRow(
			%s
		)

	syntax := "SELECT .+ FROM %s WHERE .+"
	suite.mock.ExpectQuery(syntax).
		WithArgs(
			%s
		).
		WillReturnRows(rows)

	data, err := suite.repo.%s(context.Background(), %s)
	require.NoError(err)
	require.Equal(%s, data)
	require.NoError(suite.mock.ExpectationsWereMet())
}

`

	s.selectByTestErrorNoRowsFailure = `
func (suite *%sRepositoryTestSuite) Test%s_NotFoundErr_Failure() {
	require := suite.Require()

	expectedError := Err%sNotFound

	%s

	syntax := "SELECT (.+) FROM %s WHERE (.+) "
	suite.mock.ExpectQuery(syntax).
		WithArgs(
			%s
		).
		WillReturnError(expectedError)

	data, err := suite.repo.%s(context.Background(), %s)
	require.Equal(expectedError, err)
	require.Nil(data)
	require.NoError(suite.mock.ExpectationsWereMet())
}

`

	s.selectByTestFailure = `
func (suite *%sRepositoryTestSuite) Test%s_NotFoundErr_Failure() {
	require := suite.Require()

	expectedError := errors.New("something went wrong")

	%s

	syntax := "SELECT (.+) FROM %s WHERE (.+) "
	suite.mock.ExpectQuery(syntax).
		WithArgs(
			%s
		).
		WillReturnError(expectedError)

	data, err := suite.repo.%s(context.Background(), %s)
	require.Equal(expectedError, err)
	require.Nil(data)
	require.NoError(suite.mock.ExpectationsWereMet())
}

`

	s.updateAllTestSuccess = `
func (suite *%sRepositoryTestSuite) TestUpdate_Success() {
	require := suite.Require()
	
	%s

	expRowsAffected := int64(1)

	%s

	sqlmock.NewRows([]string{
		%s		
	}).
		AddRow(
			%s
		)

	syntax := "UPDATE %s SET .+"
	suite.mock.ExpectExec(syntax).
		WithArgs(
			%s
		).
		WillReturnResult(sqlmock.NewResult(int64(rideID), expRowsAffected))

	rowsAffected, err := suite.repo.Update(context.Background(), %s, %s)
	require.NoError(err)
	require.Equal(expRowsAffected, rowsAffected)
	require.NoError(suite.mock.ExpectationsWereMet())
}

`

	s.updateAllTestFailure = `
func (suite *%sRepositoryTestSuite) TestUpdate_Failure() {
	require := suite.Require()

	expectedError := errors.New("something went wrong")

	%s

	expectedRowsAffected := int64(0)

	%s

	syntax := "UPDATE %s SET .+"
	suite.mock.ExpectExec(syntax).
		WithArgs(
			%s
		).
		WillReturnError(expectedError)

	rowsAffected, err := suite.repo.Update(context.Background(), rideID, cancellationEvent)
	require.EqualError(err, expectedError.Error())
	require.Equal(expectedRowsAffected, rowsAffected)
	require.NoError(suite.mock.ExpectationsWereMet())
}

`

	s.updateByTestSuccess = `
func (suite *%sRepositoryTestSuite) TestUpdate%s_Success() {
	require := suite.Require()

	%s

	expRowsAffected := int64(1)

	%s

	sqlmock.NewRows([]string{
		%s
	}).
		AddRow(
			%s
		)

	syntax := "UPDATE %s SET (.+) WHERE (.+)"
	suite.mock.ExpectExec(syntax).
		WithArgs(
			%s
		).
		WillReturnResult(sqlmock.NewResult(int64(1), expRowsAffected))

	rowsAffected, err := suite.repo.Update%s(
		context.Background(),
		%s,
		%s,
	)
	require.NoError(err)
	require.Equal(rowsAffected, expRowsAffected)
	require.NoError(suite.mock.ExpectationsWereMet())
}

`

	s.updateByTestFailure = `
func (suite *%sRepositoryTestSuite) TestUpdate%s_Failure() {
	require := suite.Require()

	expectedError := errors.New("something went wrong")

	%s

	expRowsAffected := int64(0)

	syntax := "UPDATE %s SET (.+) WHERE (.+)"
	suite.mock.ExpectExec(syntax).
		WithArgs(
			%s
		).
		WillReturnError(expectedError)

	rowsAffected, err := suite.repo.Update%s(
		context.Background(),
		%s,
		%s,
	)
	require.Equal(expectedError, err)
	require.Equal(rowsAffected, expRowsAffected)
	require.NoError(suite.mock.ExpectationsWereMet())
}

`

	return s
}

func (s *sqlxTest) Insert(structure *structure.Structure) (syntax string) {
	syntax = fmt.Sprintf(
		s.insertTestSuccess,
		structure.Name,
		structure.FakeDataVariables(),
		structure.FakeStructVariables(),
		structure.GetDBFieldsInQutation(),
		strcase.ToSnake(structure.Name),
		structure.GetVariableFields(strcase.ToLowerCamel(structure.Name)),
		strcase.ToLowerCamel(structure.Name),
	)

	syntax += fmt.Sprintf(
		s.insertTestFailure,
		structure.Name,
		structure.FakeDataVariables(),
		structure.FakeStructVariables(),
		strcase.ToSnake(structure.Name),
		structure.GetVariableFields(strcase.ToLowerCamel(structure.Name)),
		strcase.ToLowerCamel(structure.Name),
	)

	return syntax
}
