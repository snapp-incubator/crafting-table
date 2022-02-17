package generator

import "log"

type Variables struct {
	Name []string
}

func GenerateRepository(source, destination, packageName string, getVars, updateVars *[]Variables, create bool) error {
	filename := "repository.go"
	err := copyFile(source)
	if err != nil {
		log.Println("Error in copyfile: ", err)
		return err
	}

	//structName, err := findStructName(filename)
	//if err != nil {
	//	log.Println("Error in findStructName: ", err)
	//	return err
	//}
	//
	//instance, err := createInstance(structName, filename)
	//if err != nil {
	//	log.Println("Error in createInstance: ", err)
	//	return err
	//}
	//structure, err := bindStruct(instance)
	//if err != nil {
	//	log.Println("Error in bindStruct: ", err)
	//	return err
	//}
	//fileContent, err = createTemplate(filename, packageName)
	//if err != nil {
	//	log.Println("Error in createTemplate: ", err)
	//	return err
	//}
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
	//err = removeFile(filename)
	//if err != nil {
	//	log.Println("Error in removeFile: ", err)
	//	return err
	//}
	//
	//distanationFileName := destination + filename
	//err = writeFile(fileContent, distanationFileName)
	//if err != nil {
	//	log.Println("Error in writeFile: ", err)
	//	return err
	//}
	//err = linter(distanationFileName)
	//if err != nil {
	//	log.Println("Error in linter: ", err)
	//	return err
	//}

	return nil
}
