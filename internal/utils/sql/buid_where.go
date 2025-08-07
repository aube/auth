package sql

import (
	"fmt"
	"strings"
)

func BuildWhere(params map[string]any, orAnd string, startNum int) (whereClause string, whereParams []any) {
	if len(params) == 0 {
		return "", nil
	}

	// Условие по умолчанию AND
	if orAnd != "AND" && orAnd != "OR" {
		orAnd = "AND"
	}

	var conditions []string
	whereParams = make([]any, 0, len(params))
	paramNum := startNum
	for field, value := range params {
		if value == nil {
			continue
		}

		// Поддержка разных операторов через специальный синтаксис
		// Например: "age >" => "age > $1"
		operator := "="
		if strings.Contains(field, " ") {
			parts := strings.SplitN(field, " ", 2)
			field = parts[0]
			operator = parts[1]
		}

		conditions = append(conditions, fmt.Sprintf("%s $%d", field+" "+operator, paramNum))
		whereParams = append(whereParams, value)
		paramNum++
	}

	if len(conditions) == 0 {
		return "", nil
	}

	whereClause = strings.Join(conditions, " "+orAnd+" ")

	return whereClause, whereParams
}
