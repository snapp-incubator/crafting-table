# Crafting Table

## UNRELEASED

* Replace sprintf with text templates. (2022-11-09, @bamdadnouri, !57)
* Add query builder layer. (2022-11-13, @n25a, !58)
* Add function builder layer. (2022-11-17, @n25a, !59)
* Refactor query builder layer and function builder layer and merge them. (2022-11-18, @n25a, !60)
* Remove cli and terminal UI commands. (2022-11-18, @n25a, !61)
* Add custom function name. (2022-11-21, @parsaeisa, !69)
* Add database type to layer. (2022-11-20, @nemati21, !71)
* Fix order by in build query layer. (2022-11-27, @nemati21, !75)
* Fix bug for using table name in manifest. (2022-12-05, @n25a, !76)  
* Fix bug for checking `where` in build query layer and fix naming function. (2022-12-05, @n25a, !77)
* Add aggregate set for validation. (2022-12-05, @n25a, !78)
* Update README and add new documents. (2023-01-15, @n25a, !82)

# v2.0.0 - Nov 08 2022 

* Support yaml file for creating multiple repositories. (2022-08-05, @n25a, !29)
* Add simple terminal UI that can be used instead of CLI args. (2022-08-06, @anvari1313, !28)
* Remove additional comments in created test file. (2022-08-10, @n25a, !33) 
* Add `struct name` flag for define which struct to use in source file. (2022-08-12, @n25a, !34)
* Add manifest file for generate repositories and fix bugs in generating update tests. (2022-08-16, @n25a, !36)
* Add manifest command for generate repositories. (2022-08-16, @n25a, !37)
* Add `tags` flag to manifest command for selecting which tags to generate. (2022-08-16, @n25a, !38) 
* Remove `app` package. (2022-08-16, @nemati21, !35)
* Add `Join` function for generating join query with tests. (2022-09-01, @n25a, !44)
* Isolate sql tests by reseting sqlmock for each test. (2022-10-29, @bamdadnouri, !46)
* Remove unnecessary sqlmock row creation in sql tests. (2022-10-29, @bamdadnouri, !48)
* Add aggregate functions and group by to manifest. (2022-11-08, @n25a, !55)

# v1.2.0 - Jun 25 2022

* Clean up packages. (2022-05-20, @n25a, !15)
* Issue #13: Renaming variables and packages. (2022-06-17, @n25a, !16)
* Fix bug in query syntax for `sqlx`. (2022-06-24, @n25a, !17)
* Create tests for generated repository. (2022-06-25, @n25a, !18)

## v1.1.1 - May 17 2022

* Change go module path. (2022-05-17, @n25a, !10)

## v1.1.0 - May 15 2022

* Update the gitignore file to include the `vendor` folder and `Project binary` file. (2022-05-15, @n25a, !3)
* Add Makefile to the project. (2022-05-15, @n25a, !4)
* Add integration test. (2022-05-15, @n25a, !5)

## v1.0.1 - Apr 06 2022

* Fix bug in release GitHub action.

## v1.0.0 - Apr 06 2022

* Initial release.
