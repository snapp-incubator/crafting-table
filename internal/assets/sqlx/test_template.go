package sqlx

import (
	"fmt"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/snapp-incubator/crafting-table/internal/structure"
)

type SqlxTest interface {
	Insert(structure *structure.Structure) (syntax string)
	UpdateAll(structure *structure.Structure) (syntax string)
	UpdateBy(structure *structure.Structure, vars *[]structure.UpdateVariables) (syntax string)
	SelectBy(structure *structure.Structure, vars *[]structure.Variables) (syntax string)
	//SelectAll(structure *structure.Structure) (syntax string)
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

	expRowsAffected := int64(1)

	var %s %s.%s
	errFakeData := faker.FakeData(&%s)
	require.NoError(errFakeData)

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

	var %s %s.%s
	errFakeData := faker.FakeData(&%s)
	require.NoError(errFakeData)

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
func (suite *%sRepositoryTestSuite) TestGetBy%s_Success() {
	require := suite.Require()

	var %s %s.%s
	errFakeData := faker.FakeData(&%s)
	require.NoError(errFakeData)

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

	data, err := suite.repo.GetBy%s(
		context.Background(), 
		%s
	)
	require.NoError(err)
	require.Equal(%s, data)
	require.NoError(suite.mock.ExpectationsWereMet())
}

`

	s.selectByTestErrorNoRowsFailure = `
func (suite *%sRepositoryTestSuite) TestGetBy%s_NotFoundErr_Failure() {
	require := suite.Require()

	expectedError := Err%sNotFound

	var %s %s.%s
	errFakeData := faker.FakeData(&%s)
	require.NoError(errFakeData)

	syntax := "SELECT (.+) FROM %s WHERE (.+) "
	suite.mock.ExpectQuery(syntax).
		WithArgs(
			%s
		).
		WillReturnError(expectedError)

	data, err := suite.repo.GetBy%s(
		context.Background(),
		%s
	)
	require.Equal(expectedError, err)
	require.Nil(data)
	require.NoError(suite.mock.ExpectationsWereMet())
}

`

	s.selectByTestFailure = `
func (suite *%sRepositoryTestSuite) TestGetBy%s_OtherErr_Failure() {
	require := suite.Require()

	expectedError := errors.New("something went wrong")

	var %s %s.%s
	errFakeData := faker.FakeData(&%s)
	require.NoError(errFakeData)

	syntax := "SELECT (.+) FROM %s WHERE (.+) "
	suite.mock.ExpectQuery(syntax).
		WithArgs(
			%s
		).
		WillReturnError(expectedError)

	data, err := suite.repo.GetBy%s(
		context.Background(),
		%s
	)
	require.Equal(expectedError, err)
	require.Nil(data)
	require.NoError(suite.mock.ExpectationsWereMet())
}

`

	s.updateAllTestSuccess = `
func (suite *%sRepositoryTestSuite) TestUpdate_Success() {
	require := suite.Require()

	expRowsAffected := int64(1)

	var %s %s.%s
	errFakeData := faker.FakeData(&%s)
	require.NoError(errFakeData)

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
		WillReturnResult(sqlmock.NewResult(int64(1), expRowsAffected))

	rowsAffected, err := suite.repo.Update(context.Background(), %s.%s, %s)
	require.NoError(err)
	require.Equal(expRowsAffected, rowsAffected)
	require.NoError(suite.mock.ExpectationsWereMet())
}

`

	s.updateAllTestFailure = `
func (suite *%sRepositoryTestSuite) TestUpdate_Failure() {
	require := suite.Require()

	expectedError := errors.New("something went wrong")

	var %s %s.%s
	errFakeData := faker.FakeData(&%s)
	require.NoError(errFakeData)

	expectedRowsAffected := int64(0)

	syntax := "UPDATE %s SET .+"
	suite.mock.ExpectExec(syntax).
		WithArgs(
			%s
		).
		WillReturnError(expectedError)

	rowsAffected, err := suite.repo.Update(context.Background(), %s.%s, %s)
	require.EqualError(err, expectedError.Error())
	require.Equal(expectedRowsAffected, rowsAffected)
	require.NoError(suite.mock.ExpectationsWereMet())
}

`

	s.updateByTestSuccess = `
func (suite *%sRepositoryTestSuite) TestUpdate%s_Success() {
	require := suite.Require()

	expRowsAffected := int64(1)

	var %s %s.%s
	errFakeData := faker.FakeData(&%s)
	require.NoError(errFakeData)

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
		%s
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

	var %s %s.%s
	errFakeData := faker.FakeData(&%s)
	require.NoError(errFakeData)

	expRowsAffected := int64(0)

	syntax := "UPDATE %s SET (.+) WHERE (.+)"
	suite.mock.ExpectExec(syntax).
		WithArgs(
			%s
		).
		WillReturnError(expectedError)

	rowsAffected, err := suite.repo.Update%s(
		context.Background(),
		%s
	)
	require.Equal(expectedError, err)
	require.Equal(rowsAffected, expRowsAffected)
	require.NoError(suite.mock.ExpectationsWereMet())
}

`

	return &s
}

func (s *sqlxTest) Insert(structure *structure.Structure) (syntax string) {
	syntax = fmt.Sprintf(
		s.insertTestSuccess,
		structure.Name,

		strcase.ToLowerCamel(structure.Name),
		structure.PackageName,
		structure.Name,
		strcase.ToLowerCamel(structure.Name),

		structure.GetDBFieldsInQuotation(),
		structure.GetVariableFields(strcase.ToLowerCamel(structure.Name)+"."),

		strcase.ToSnake(structure.Name),
		structure.GetVariableFields(strcase.ToLowerCamel(structure.Name)+"."),

		strcase.ToLowerCamel(structure.Name),
	)

	syntax += fmt.Sprintf(
		s.insertTestFailure,
		structure.Name,

		strcase.ToLowerCamel(structure.Name),
		structure.PackageName,
		structure.Name,
		strcase.ToLowerCamel(structure.Name),

		strcase.ToSnake(structure.Name),
		structure.GetVariableFields(strcase.ToLowerCamel(structure.Name)+"."),

		strcase.ToLowerCamel(structure.Name),
	)

	return syntax
}

func (s *sqlxTest) UpdateAll(structure *structure.Structure) (syntax string) {
	syntax = fmt.Sprintf(
		s.updateAllTestSuccess,
		structure.Name,

		strcase.ToLowerCamel(structure.Name),
		structure.PackageName,
		structure.Name,
		strcase.ToLowerCamel(structure.Name),

		structure.GetDBFieldsInQuotation(),
		structure.GetVariableFields(strcase.ToLowerCamel(structure.Name)+"."),

		strcase.ToSnake(structure.Name),
		structure.GetVariableFields(strcase.ToLowerCamel(structure.Name)+"."),

		strcase.ToLowerCamel(structure.Name),
		structure.Fields[0].Name,
		strcase.ToLowerCamel(structure.Name),
	)

	syntax += fmt.Sprintf(
		s.updateAllTestFailure,
		structure.Name,

		strcase.ToLowerCamel(structure.Name),
		structure.PackageName,
		structure.Name,
		strcase.ToLowerCamel(structure.Name),

		strcase.ToSnake(structure.Name),
		structure.GetVariableFields(strcase.ToLowerCamel(structure.Name)+"."),

		strcase.ToLowerCamel(structure.Name),
		structure.Fields[0].Name,
		strcase.ToLowerCamel(structure.Name),
	)

	return syntax
}

func (s *sqlxTest) UpdateBy(structure *structure.Structure, vars *[]structure.UpdateVariables) (syntax string) {
	for _, v := range *vars {
		functionNameList := make([]string, 0)
		for _, name := range v.Fields {
			functionNameList = append(functionNameList, structure.FieldDBNameToName[name])
		}
		functionName := strings.Join(functionNameList, "And")

		syntax = fmt.Sprintf(
			s.updateByTestSuccess,
			structure.Name,
			functionName,

			strcase.ToLowerCamel(structure.Name),
			structure.PackageName,
			structure.Name,
			strcase.ToLowerCamel(structure.Name),

			structure.GetDBFieldsInQuotation(),
			structure.GetVariableFields(strcase.ToLowerCamel(structure.Name)+"."),

			strcase.ToSnake(structure.Name),
			s.addPrefix(execContextVariables(v, structure, true), strcase.ToLowerCamel(structure.Name)+"."),

			functionName,
			s.addPrefix(execContextVariables(v, structure, false), strcase.ToLowerCamel(structure.Name)+"."),
		)

		syntax += fmt.Sprintf(
			s.updateByTestFailure,
			structure.Name,
			functionName,

			strcase.ToLowerCamel(structure.Name),
			structure.PackageName,
			structure.Name,
			strcase.ToLowerCamel(structure.Name),

			strcase.ToSnake(structure.Name),
			s.addPrefix(execContextVariables(v, structure, true), strcase.ToLowerCamel(structure.Name)+"."),
			functionName,

			s.addPrefix(execContextVariables(v, structure, false), strcase.ToLowerCamel(structure.Name)+"."),
		)
	}

	return syntax
}

func (s *sqlxTest) SelectBy(structure *structure.Structure, vars *[]structure.Variables) (syntax string) {

	for _, v := range *vars {
		functionNameList := make([]string, 0)
		for _, name := range v.Name {
			functionNameList = append(functionNameList, structure.FieldDBNameToName[name])
		}
		functionName := strings.Join(functionNameList, "And")

		syntax = fmt.Sprintf(
			s.selectByTestSuccess,
			structure.Name,
			functionName,

			strcase.ToLowerCamel(structure.Name),
			structure.PackageName,
			structure.Name,
			strcase.ToLowerCamel(structure.Name),

			structure.GetDBFieldsInQuotation(),
			structure.GetVariableFields(strcase.ToLowerCamel(structure.Name)+"."),

			strcase.ToSnake(structure.Name),
			s.addPrefix(contextVariables(v.Name, structure), strcase.ToLowerCamel(structure.Name)+"."),

			functionName,
			s.addPrefix(contextVariables(v.Name, structure), strcase.ToLowerCamel(structure.Name)+"."),
			strcase.ToLowerCamel(structure.Name),
		)

		syntax += fmt.Sprintf(
			s.selectByTestErrorNoRowsFailure,
			structure.Name,
			functionName,
			structure.Name,

			strcase.ToLowerCamel(structure.Name),
			structure.PackageName,
			structure.Name,
			strcase.ToLowerCamel(structure.Name),

			strcase.ToSnake(structure.Name),
			s.addPrefix(contextVariables(v.Name, structure), strcase.ToLowerCamel(structure.Name)+"."),

			functionName,
			s.addPrefix(contextVariables(v.Name, structure), strcase.ToLowerCamel(structure.Name)+"."),
		)

		syntax += fmt.Sprintf(
			s.selectByTestFailure,
			structure.Name,
			functionName,

			strcase.ToLowerCamel(structure.Name),
			structure.PackageName,
			structure.Name,
			strcase.ToLowerCamel(structure.Name),

			strcase.ToSnake(structure.Name),
			s.addPrefix(contextVariables(v.Name, structure), strcase.ToLowerCamel(structure.Name)+"."),

			functionName,
			s.addPrefix(contextVariables(v.Name, structure), strcase.ToLowerCamel(structure.Name)+"."),
		)
	}

	return syntax
}

func (s *sqlxTest) addPrefix(str, prefix string) string {
	tmp := strings.Split(str, ",")
	result := ""
	for _, v := range tmp {
		result += prefix + strcase.ToCamel(strings.Replace(v, " ", "", -1)) + ",\n\t\t\t"
	}

	return result
}
