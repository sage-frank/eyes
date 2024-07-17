package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"time"
	// read and parse mysql binlog
	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-mysql-org/go-mysql/replication"
)

// MyWriter 是自定义的io.Writer实现
type MyWriter struct {
	content *bytes.Buffer
}

// Write 实现io.Writer接口的Write方法
func (m *MyWriter) Write(p []byte) (n int, err error) {
	return m.content.Write(p)
}

func main() {
	cfg := replication.BinlogSyncerConfig{
		ServerID: 100,
		Flavor:   "mysql",
		Host:     "127.0.0.1",
		Port:     33061,
		User:     "root",
		Password: "root",
	}
	syncer := replication.NewBinlogSyncer(cfg)
	streamer, _ := syncer.StartSync(mysql.Position{Name: "", Pos: 3025278})

	//for {
	//	ev, _ := streamer.GetEvent(context.Background())
	//	// Dump event
	//	ev.Dump(os.Stdout)
	//}

	// or we can use a timeout context
	for {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		ev, err := streamer.GetEvent(ctx)
		cancel()
		//
		//select {
		//case <-ctx.Done():
		//	continue
		//default:
		//	ev.Dump(os.Stdout)
		//}

		if errors.Is(err, context.DeadlineExceeded) {
			// meet timeout
			continue
		}
		buf := bytes.Buffer{}
		writer := MyWriter{content: &buf}
		var myWriter io.Writer = &writer

		ev.Dump(myWriter)
		fmt.Println("ev.rawdata", buf.String())
		fmt.Println("--------------------------------------------------")
		// ret := make([]byte, 100)
		// ev.Event.Decode(ret)
		// fmt.Println(ret)
		// ev.Dump(os.Stdout)
	}
}
