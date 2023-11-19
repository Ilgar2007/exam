package helpers

import (
	"database/sql"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func NewNullString(s string) sql.NullString {

	if len(s) == 0 {
		return sql.NullString{}
	}

	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

func NewIncrementId(db *sql.DB, tableName string, Length int) (func() string, error) {
	var (
		id            sql.NullString
		capitalLetter = string(strings.ToUpper(tableName)[0])
		query         = fmt.Sprintf("SELECT id FROM %s ORDER BY created_at DESC LIMIT 1", tableName)
	)
	resp := db.QueryRow(query)

	resp.Scan(&id)

	if !id.Valid {
		fmt.Printf("$$$$$$$$$$$$$$     %+v     $$$$$$$$$$$$$$", id)
		id.String = ""
	}
	return func() string {
		idNumber := idToInt(id.String)
		idNumber++
		var (
			numberLenght = idNumber
			count        = 0
		)
		for numberLenght > 0 {
			numberLenght /= 10
			count++
		}
		if count == 0 {
			count++
		}
		fmt.Printf("COUNT %d  dnLength %d  &&&&&&&&&&&\n", count, idNumber)

		id := fmt.Sprintf("%s-%s%d", capitalLetter, strings.Repeat("O", Length-count), idNumber)
		return id

	}, nil
}

func idToInt(DatabaseId string) int {
	pattern := regexp.MustCompile("[0-9]+")
	firstMatchSubstring := pattern.FindString(DatabaseId)
	id, err := strconv.Atoi(firstMatchSubstring)
	if err != nil {
		return 0
	}
	fmt.Println("#####################@@@@@@  ", id, "   ", firstMatchSubstring, "  @@@@@@####################")

	return id
}
