package io_better

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
	"syscall"
)

func streamWrite(r io.Reader, w io.Writer) {
	eg := &sync.WaitGroup{}
	// init csvWriter, BufioScanner, Io.Pipe, Chunk
	chunkBuf = &bytes.Buffer{}
	csvWriter = csv.NewWriter(csvBuf)

	rp, wp := io.Pipe()
	defer wp.Close()

	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)
	// 读取brk文件，分割成csv格式，写入pipe
	eg.Go(func() error {
		lineNum := 0
		inter := 1
		for scanner.Scan() {
			var (
				items []string
			)
			lineNum++
			if lineNum/50000 == inter {
				inter++
				log.Warn("TaskCsv process line(%d)", lineNum)
			}

			if lineNum > 700000 {
				log.Error("文件太大了(%d) | params (%+v) ", lineNum, v)
				csvBuf.Reset()
				wp.Close()
				return fmt.Errorf("文件太大了 (%d)", lineNum)
			}
			//\001 brk分割符
			fields := bytes.Split(scanner.Bytes(), []byte("\001"))
			if len(fields) < 7 {
				log.Warn("TaskCsv fields lenth too small:%d", len(fields))
				continue
			}
			for _, field := range fields {
				items = append(items, string(field))
			}
			if err = csvWriter.Write(items); err != nil {
				log.Error("TaskCsv csvWriter.Write(%v) erorr(%v)", items, err)
				return err
			}
			items = nil

			if lineNum/5000 == inter {
				inter++
				log.Info("TaskCsv flush line(%d)", lineNum)
				csvWriter.Flush() // todo 多久flush合适, 目前5000行flush 约2Mb
				if _, err = wp.Write(csvBuf.Bytes()); err != nil {
					log.Error("TaskCsv Write error(%v)", err)
					csvBuf.Reset()
					wp.Close()
					return err
				}
				csvBuf.Reset()
				//todo 审核下载速度一般小于内网下载速度	time.sleep()
			}
		}
		//scan结束， Close Pipe
		csvWriter.Flush()
		if _, err = wp.Write(csvBuf.Bytes()); err != nil {
			log.Error("TaskCsv Write error(%v)", err)
			return err
		}
		csvBuf.Reset()
		wp.Close()
		if err != nil {
			log.Error("TaskCsv write error(%v)", err)
			return err
		}
		if lineNum == 0 {
			return ecode.NothingFound
		}
		return nil
	})
	if err = csvWriter.Error(); err != nil {
		log.Error("TaskCsv csvWriter error(%v)", err)
		return nil, nil, err
	}
	if scanner.Err() != nil {
		log.Error("TaskCsv scanner error(%v)", scanner.Err())
		return nil, nil, scanner.Err()
	}

	// 先写header 再写body
	c.Writer.Header().Set("Content-Type", "application/csv")
	c.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment;filename=%v.csv", 123))
	c.Writer.Header().Set("Transfer-Encoding", "chunked")
	c.Writer.WriteHeader(http.StatusOK)
	//c.Writer.Write([]byte("\xEF\xBB\xBF")) // fix ms乱码

	eg.Go(func() error {
		_, err = io.Copy(c.Writer, rp)
		if err != nil {
			log.Error("taskcsv copy error(%v)", err)
			wp.Close()
			rp.Close()
			if errors.Is(err, syscall.EPIPE) || errors.Is(err, syscall.EPROTOTYPE) {
				log.Error("taskcsv client cancel error(%v)", err)
				return ecode.Canceled
			} else {
				return err
			}
		}
		return nil
	})
	if err = eg.Wait(); err != nil {
		log.Error("taskcsv wait error(%v)", err)
		return nil, nil, err
	}
	return nil, nil, nil
}

}
