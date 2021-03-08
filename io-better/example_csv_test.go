package io_better

import (
	"bernieyangmh.com/TinyComponent/utils"
	"bufio"
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"syscall"
)

func WriteCsvMetadata(w http.ResponseWriter, filename string, titleLine []string) error {
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment;filename=%s.csv", filename))
	w.Header().Set("Content-Type", "application/csv")
	w.Header().Set("Transfer-Encoding", "chunked")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("\xEF\xBB\xBF")) // fix ms乱码
	if err := csv.NewWriter(w).Write(titleLine); err != nil {
		return err
	}
	return nil
}

func WriteCsvResp(w io.Writer, reader io.Reader) error {

	_, err := io.Copy(w, reader)
	if err != nil {
		if errors.Is(err, syscall.EPIPE) || errors.Is(err, syscall.EPROTOTYPE) {
			return errors.New("client close connection")
		}
		return err
	}
	return nil
}

//, parse func(r io.Reader, w io.Writer)
func ReadToCsv(reader io.Reader, writer io.Writer) (err error) {
	var (
		csvBuff = &bytes.Buffer{}

		csvWriter = csv.NewWriter(csvBuff)
		scanner   = bufio.NewScanner(reader)
	)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		if err = NormalParse(csvWriter, scanner.Bytes()); err != nil {
			return err
		}
		if _, err = writer.Write(csvBuff.Bytes()); err != nil {
			return err
		}
	}
	return nil
}

// eg: fields split by "\001"
func NormalParse(w *csv.Writer, data []byte) error {
	fields := bytes.Split(data, []byte("\001"))
	lines := []string{}
	for _, f := range fields {
		lines = append(lines, utils.ByteToString(f))
	}
	if err := w.Write(lines); err != nil {
		return err
	}
	w.Flush()
	return nil
}

func TransforCsv(w http.ResponseWriter, r *http.Request) {
	var (
		err error
	)
	defer func() {
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("error file"))
			return
		}
	}()

	filePath := ""
	dir, err := os.OpenFile(filePath, os.O_WRONLY, 0)

	if err = WriteCsvMetadata(w, "test", []string{"test", "time"}); err != nil {
		fmt.Println("WriteCsvMetadata error(%v)", err)
	}

	if err = ParallelStreamIO(w, dir, WriteCsvResp, ReadToCsv); err != nil {
		fmt.Println("ParallelStreamIO error(%v)", err)
	}
	return
}

func main() {
	http.HandleFunc("/csv/", TransforCsv)
	http.ListenAndServe("127.0.0.0:8000", nil)
}
