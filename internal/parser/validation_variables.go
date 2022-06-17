package parser

import (
	"errors"
	"regexp"
	"strings"
)

func ValidateGetFlag(flag string) error {
	// TODO : change validation to regex

	if string(flag[0]) != "[" && string(flag[len(flag)-1]) != "]" {
		return errors.New("you must set get variables in format of [ (var1,var2), (var2,var4), var3 ]")
	}

	openParentheses := false
	for index, char := range flag {
		if openParentheses && char == '(' {
			return errors.New("open parentheses are not closed")
		}

		if !openParentheses && char == ')' {
			return errors.New("close parentheses are not opened")
		}

		if openParentheses && char == ')' && flag[index-1] == ',' {
			return errors.New("close parentheses must not be followed by comma")
		}

		if openParentheses && char == ')' && flag[index+1] != ']' && flag[index+1] != ',' {
			return errors.New("close parentheses must be followed by comma")
		}

		if char == '(' {
			openParentheses = true
		}

		if char == ')' {
			openParentheses = false
		}
	}
	return nil
}

func ValidateUpdateFlag(flag string) error {
	cleanFlag := flag
	for strings.Contains(cleanFlag, "  ") {
		cleanFlag = strings.Replace(flag, "  ", " ", -1)
	}
	cleanFlag = strings.Replace(flag, " ", "", -1)

	regex := `^\[(\[\(([a-zA-Z]+([a-zA-Z]*[0-9]*\_*)*,)*|([a-zA-Z]+([a-zA-Z]*[0-9]*\_*)*)?\),\(([a-zA-Z]+([a-zA-Z]*[0-9]*\_*)*,)*|([a-zA-Z]+([a-zA-Z]*[0-9]*\_*)*)?\)\],)+|(\[\(([a-zA-Z]+([a-zA-Z]*[0-9]*\_*)*,)*|([a-zA-Z]+([a-zA-Z]*[0-9]*\_*)*)?\),\(([a-zA-Z]+([a-zA-Z]*[0-9]*\_*)*,)*|([a-zA-Z]+([a-zA-Z]*[0-9]*\_*)*)?\)\])+\]$`
	found, err := regexp.MatchString(regex, cleanFlag)

	if err != nil {
		return err
	}

	if !found {
		return errors.New("you must set update variables in format of [ [(byPar1,byPar2,...), (field1, field2,...)], ... ]")
	}

	return nil
}
