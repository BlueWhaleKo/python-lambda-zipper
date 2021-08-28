package util

import (
	"io"
	"os"
)

func Print(r io.ReadCloser) error {
	defer r.Close()

	_, err := io.Copy(os.Stdout, r)
	if err != nil {
		return err
	}

	return nil
}

//func print(rd io.Reader) error {
//	var lastLine string
//
//	scanner := bufio.NewScanner(rd)
//	for scanner.Scan() {
//		lastLine = scanner.Text()
//		fmt.Println(scanner.Text())
//	}
//
//	errLine := &ErrorLine{}
//	json.Unmarshal([]byte(lastLine), errLine)
//	if errLine.Error != "" {
//		return errors.New(errLine.Error)
//	}
//
//	if err := scanner.Err(); err != nil {
//		return err
//	}
//
//	return nil
//}
