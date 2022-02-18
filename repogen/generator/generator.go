package generator

import "log"

type Variables struct {
	Name []string
}

func GenerateRepository(source, destination, packageName string, getVars, updateVars *[]Variables, create bool) error {
	fileName, err := getFileName(source)
	if err != nil {
		log.Println("Error in getFileName: ", err)
		return err
	}

	structure, err := bindStruct(source)
	if err != nil {
		log.Println("Error in bindStruct: ", err)
		return err
	}

	//interfaceSyntax, err := interfaceSyntaxCreator(structure, getVars, updateVars, create)
	//if err != nil {
	//	log.Println("Error in interfaceSyntaxCreator: ", err)
	//	return err
	//}

	fileContent := createTemplate(fileName, packageName, interfaceSyntax, structure)

	//if create {
	//	err = createFunctionRepository(fileContent, filename, structure)
	//	if err != nil {
	//		log.Println("Error in createFunctionRepository: ", err)
	//		return err
	//	}
	//}
	//
	//if getVars != nil {
	//	err = getFunctionCreator(fileContent, getVars)
	//	if err != nil {
	//		log.Println("Error in getFunctionCreator: ", err)
	//		return err
	//	}
	//}
	//
	//if updateVars != nil {
	//	err = updateFunctionCreator(fileContent, updateVars)
	//	if err != nil {
	//		log.Println("Error in updateFunctionCreator: ", err)
	//		return err
	//	}
	//}
	//

	err = writeFile(fileContent, destination)
	if err != nil {
		log.Println("Error in writeFile: ", err)
		return err
	}

	//err = linter(distanationFileName)
	//if err != nil {
	//	log.Println("Error in linter: ", err)
	//	return err
	//}

	return nil
}
