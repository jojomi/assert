package assert

import (
	"fmt"
	"github.com/jojomi/assert/exit"
	"github.com/jojomi/assert/ranges"
	"strings"
	"time"
)

var (
	allWeekdays = []time.Weekday{
		time.Monday,
		time.Tuesday,
		time.Wednesday,
		time.Thursday,
		time.Friday,
		time.Saturday,
		time.Sunday,
	}
)

func Weekday(args []string) (error, exit.Code) {
	currentWeekday := time.Now().Weekday().String()

	shortMap := makeUniqueMap(allWeekdays)
	mappedArgs, err := mapSliceErr(args, func(elem string) (string, error) {
		// range?
		parts := strings.SplitN(elem, "-", -1)

		mappedParts, err := mapSliceErr(parts, func(elem string) (string, error) {
			weekday, ok := shortMap[elem]
			if !ok {
				return "", fmt.Errorf("invalid arg: %s", elem)
			}
			return weekday.String(), nil
		})

		if err != nil {
			return "", err
		}

		return strings.Join(mappedParts, "-"), nil
	})
	if err != nil {
		return err, exit.CodeErrorFinal
	}

	allWeekdayStrings := mapSlice(allWeekdays, func(elem time.Weekday) string {
		return elem.String()
	})
	resolvedArgs, err := ranges.ResolveIndexedRanges(mappedArgs, allWeekdayStrings)
	if err != nil {
		return err, exit.CodeErrorFinal
	}

	// one of?
	for _, arg := range resolvedArgs {
		if arg == currentWeekday {
			return nil, exit.CodeOkay
		}
	}

	return fmt.Errorf("today is not one of %s (input was \"%s\")", strings.Join(resolvedArgs, ", "), strings.Join(args, " ")), exit.CodeErrorFinal
}

func makeUniqueMap(weekdays []time.Weekday) map[string]time.Weekday {
	result := make(map[string]time.Weekday)
	for _, weekday := range weekdays {
		weekdayString := weekday.String()
		sub := ""
		for _, char := range weekdayString {
			sub += strings.ToLower(string(char))
			if _, ok := result[sub]; ok {
				delete(result, sub)
				continue
			}
			result[sub] = weekday
		}
	}
	return result
}

func mapSlice[T, S any](input []T, mapper func(elem T) S) []S {
	res, _ := mapSliceErr(input, func(elem T) (S, error) {
		res := mapper(elem)
		return res, nil
	})
	return res
}

func mapSliceErr[T, S any](input []T, mapper func(elem T) (S, error)) ([]S, error) {
	var (
		result = make([]S, len(input))
		err    error
	)
	for i, elem := range input {
		result[i], err = mapper(elem)
		if err != nil {
			return result, err
		}
	}
	return result, nil
}
