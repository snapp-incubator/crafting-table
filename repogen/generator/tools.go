package generator

//
//func ErrNotFound(tableName string) error {
//	text := tableName + " not found"
//	return errors.New(text)

//}

//func getSyntaxCreator(tableName string, parms map[string]interface{}) string {
//
//	syntax := "SELECT * FROM " + tableName
//
//	whereFlage := true
//	for key := range parms {
//		if whereFlage {
//			syntax += " WHERE "
//			whereFlage = false
//		}
//		syntax += key + " = ?" + " AND "
//	}
//
//	syntax = syntax[:len(syntax)-4] + ";"
//
//	return syntax
//}
//
//func getValues(parms map[string]interface{}) []interface{} {
//
//	values := make([]interface{}, 0)
//
//	for _, value := range parms {
//		values = append(values, value)
//	}
//
//	return values
//}
