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
	"runtime"
	"syscall"
	"testing"
)

func WriteCsvHeader(w http.ResponseWriter, filename string, titleLine []string) error {
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
	n := 0
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		if err = NormalParse(csvWriter, scanner.Bytes()); err != nil {
			return err
		}
		n++
		if n == 4000 {
			csvWriter.Flush()
		}
		if _, err = writer.Write(csvBuff.Bytes()); err != nil {
			return err
		}
		csvBuff.Reset()
	}
	csvWriter.Flush()
	if _, err = writer.Write(csvBuff.Bytes()); err != nil {
		return err
	}
	csvBuff.Reset()
	if scanner.Err() != nil {
		return scanner.Err()
	}
	return io.EOF
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
	return nil
}

func TransforCsv(w http.ResponseWriter, r *http.Request) {
	var (
		err error
	)
	defer func() {
		if r := recover(); r != nil {
			buf := make([]byte, 64<<10)
			buf = buf[:runtime.Stack(buf, false)]
			fmt.Println(fmt.Sprintf("TransforCsv panic error(%v)  \n%s", r, buf))
		}

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("error file"))
			return
		}
	}()

	filePath := "/Users/yangminghui/Desktop/aid.csv"
	dir, err := os.OpenFile(filePath, os.O_RDONLY, 0)
	if err != nil {
		fmt.Printf("OpenFile error(%v)", err)
		return
	}
	if err = WriteCsvHeader(w, "test", []string{"test", "time"}); err != nil {
		fmt.Printf("WriteCsvHeader error(%v)\n", err)
		return
	}

	if err = ParallelStreamIO(w, dir, WriteCsvResp, ReadToCsv); err != nil {
		fmt.Printf("ParallelStreamIO error(%v)\n", err)
		return
	}
	return
}

func TestHttpServer(t *testing.T) {
	t.Run("run server", func(t *testing.T) {
		http.HandleFunc("/csv/", TransforCsv)
		if err := http.ListenAndServe(":8000", nil); err != nil {
			t.Logf("ListenAndServe error(%v)", err)
		}
	})
}
