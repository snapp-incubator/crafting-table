package sqlx

import (
	"errors"
	"fmt"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/snapp-incubator/crafting-table/internal/structure"
)

type SqlxTest interface {
	Insert(structure *structure.Structure) (syntax string)
	UpdateAll(structure *structure.Structure) (syntax string)
	UpdateBy(structure *structure.Structure, vars *[]structure.UpdateVariables) (syntax string)
	SelectBy(structure *structure.Structure, vars *[]structure.GetVariable) (syntax string)
	Join(structure *structure.Structure, vars *structure.JoinVariables) (syntax string)
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

	joinTestSuccess string
	joinTestFailure string
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

	s.joinTestSuccess = `
func (suite *%sRepositoryTestSuite) TestGetJoined%s_Success() {
	require := suite.Require()
	limit := uint(1)	

	var %s %s.%s
	errFakeData := faker.FakeData(&%s)
	require.NoError(errFakeData)

	rows := sqlmock.NewRows([]string{
		%s
	}).
		AddRow(
			%s
		)

	query := "SELECT " +
		%s
		"FROM %s AS %s " +
		%s +
		"LIMIT ?"

	suite.mock.ExpectQuery(query).
		WithArgs(limit).
		WillReturnRows(rows)

	data, err := suite.repo.GetJoined%s(context.Background(), limit)
	require.NoError(err)
	require.Equal(&%s, data)
	require.NoError(suite.mock.ExpectationsWereMet())
}

`
	s.joinTestFailure = `
func (suite *%sRepositoryTestSuite) TestGetJoined%s_Failure() {
	require := suite.Require()
	limit := uint(1)	
	expectedError := errors.New("something went wrong")

	query := "SELECT " +
		%s
		"FROM %s AS %s " +
		%s +
		"LIMIT ?"

	suite.mock.ExpectQuery(query).
		WithArgs(limit).
		WillReturnError(expectedError)

	data, err := suite.repo.GetJoined%s(context.Background(), limit)
	require.Equal(expectedError, err)
	require.Nil(data)
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
			functionNameList = append(functionNameList, structure.FieldMapDBFlagToName[name])
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
			s.addPrefix(execContextVariables(v, structure, false), strcase.ToLowerCamel(structure.Name)+"."),

			functionName,
			s.addPrefix(execContextVariables(v, structure, true), strcase.ToLowerCamel(structure.Name)+"."),
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
			s.addPrefix(execContextVariables(v, structure, false), strcase.ToLowerCamel(structure.Name)+"."),
			functionName,

			s.addPrefix(execContextVariables(v, structure, true), strcase.ToLowerCamel(structure.Name)+"."),
		)
	}

	return syntax
}

func (s *sqlxTest) SelectBy(structure *structure.Structure, vars *[]structure.GetVariable) (syntax string) {

	for _, v := range *vars {
		functionNameList := make([]string, 0)
		for _, conditions := range v.Conditions {
			functionNameList = append(functionNameList, structure.FieldMapDBFlagToName[conditions])
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
			s.addPrefix(contextVariables(v.Conditions, structure), strcase.ToLowerCamel(structure.Name)+"."),

			functionName,
			s.addPrefix(contextVariables(v.Conditions, structure), strcase.ToLowerCamel(structure.Name)+"."),
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
			s.addPrefix(contextVariables(v.Conditions, structure), strcase.ToLowerCamel(structure.Name)+"."),

			functionName,
			s.addPrefix(contextVariables(v.Conditions, structure), strcase.ToLowerCamel(structure.Name)+"."),
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
			s.addPrefix(contextVariables(v.Conditions, structure), strcase.ToLowerCamel(structure.Name)+"."),

			functionName,
			s.addPrefix(contextVariables(v.Conditions, structure), strcase.ToLowerCamel(structure.Name)+"."),
		)
	}

	return syntax
}

func (s *sqlxTest) Join(structure *structure.Structure, joinVariables *structure.JoinVariables) (syntax string) {

	syntax += fmt.Sprintf(
		s.joinTestSuccess,
		structure.Name,
		structure.Name,
		strcase.ToLowerCamel(structure.Name),
		structure.PackageName,
		structure.Name,
		strcase.ToLowerCamel(structure.Name),

		joinedStringRows(structure, joinVariables),
		joinedVariablesRows(structure, joinVariables),

		joinField(structure, joinVariables),

		structure.TableName,
		string(structure.TableName[0]),

		joins(structure, joinVariables),

		structure.Name,
		strcase.ToLowerCamel(structure.Name),
	)

	syntax += fmt.Sprintf(
		s.joinTestFailure,
		structure.Name,
		structure.Name,

		joinField(structure, joinVariables),

		structure.TableName,
		string(structure.TableName[0]),

		joins(structure, joinVariables),

		structure.Name,
	)
	return syntax
}

func joinedStringRows(s *structure.Structure, joinVariables *structure.JoinVariables) (fields string) {
	// add first struct fields
	for dbName, _ := range s.FieldMapDBFlagToName {
		isEqual := false
		for _, joinVariable := range joinVariables.Fields {
			if dbName == joinVariable.JoinFieldAs {
				isEqual = true
				break
			}
		}
		if isEqual {
			continue
		}
		fields += fmt.Sprintf("\"%s\", \n\t\t", dbName)
	}

	// add join fields
	for _, joinVariable := range joinVariables.Fields {
		source := strings.Replace(joinVariable.JoinStructPath, " ", "", -1)

		ss, err := structure.BindStruct(source, joinVariable.JoinStructName)
		if err != nil {
			err = errors.New(fmt.Sprintf("Error in bindStruct: %s", err.Error()))
			panic(err)
		}

		for dbName, _ := range ss.FieldMapDBFlagToName {
			fields += fmt.Sprintf("\"%s.%s\", \n\t\t ",
				joinVariable.JoinFieldAs, dbName)
		}

		fields = strings.TrimSuffix(fields, "\n\t\t")
	}

	return fields
}

func joinedVariablesRows(s *structure.Structure, joinVariables *structure.JoinVariables) (fields string) {
	// add first struct fields
	for dbName, _ := range s.FieldMapDBFlagToName {
		isEqual := false
		for _, joinVariable := range joinVariables.Fields {
			if dbName == joinVariable.JoinFieldAs {
				isEqual = true
				break
			}
		}
		if isEqual {
			continue
		}
		fields += fmt.Sprintf("\"%s.%s\", \n\t\t",
			strcase.ToLowerCamel(s.Name), dbName)
	}

	// add join fields
	for _, joinVariable := range joinVariables.Fields {
		source := strings.Replace(joinVariable.JoinStructPath, " ", "", -1)

		ss, err := structure.BindStruct(source, joinVariable.JoinStructName)
		if err != nil {
			err = errors.New(fmt.Sprintf("Error in bindStruct: %s", err.Error()))
			panic(err)
		}

		for dbName, _ := range ss.FieldMapDBFlagToName {
			fields += fmt.Sprintf("\"%s.%s.%s\", \n\t\t ",
				strcase.ToLowerCamel(s.Name),
				strcase.ToCamel(joinVariable.JoinFieldAs),
				dbName,
			)
		}

		fields = strings.TrimSuffix(fields, "\n\t\t ")
	}

	return fields
}

func (s *sqlxTest) addPrefix(str, prefix string) string {
	tmp := strings.Split(str, ",")
	result := ""
	for _, v := range tmp {
		result += prefix + strcase.ToCamel(strings.Replace(v, " ", "", -1)) + ",\n\t\t\t"
	}

	return result
}
