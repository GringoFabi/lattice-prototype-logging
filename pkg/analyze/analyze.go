package analyze

import (
	"bufio"
	"encoding/json"
	"os"
)

func ReadJson(path string) ActionCount {
	bytes, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	logs := new(Logs)
	if err = json.Unmarshal(bytes, &logs); err != nil {
		panic(err)
	}

	return logs.CountActions()
}

func ReadJsonLines(path string) ActionCount {
	file, err := os.Open(path)

	if err != nil {
		panic(err)
	}

	defer func(file *os.File) {
		if err = file.Close(); err != nil {
			panic(err)
		}
	}(file)

	logs := make(Logs, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		log := new(Log)
		if err = json.Unmarshal(scanner.Bytes(), &log); err != nil {
			panic(err)
		}
		logs = append(logs, *log)
	}

	return logs.CountActions()
}

func ReadJsonLinesArray(path string) ActionCount {
	file, err := os.Open(path)

	if err != nil {
		panic(err)
	}

	defer func(file *os.File) {
		if err = file.Close(); err != nil {
			panic(err)
		}
	}(file)

	allLogs := make(Logs, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		logs := new([]Log)
		if err = json.Unmarshal(scanner.Bytes(), &logs); err != nil {
			panic(err)
		}
		allLogs = append(allLogs, *logs...)
	}

	return allLogs.CountActions()
}

func Zip(a ActionCount, b ActionCount, c ActionCount) ActionCount {
	zip := make(ActionCount)
	merge(a, zip)
	merge(b, zip)
	merge(c, zip)
	return zip
}

func merge(actionCount ActionCount, zip ActionCount) {
	for action, count := range actionCount {
		if _, ok := zip[action]; !ok {
			zip[action] = count
		} else {
			zip[action] += count
		}
	}
}
