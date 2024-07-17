package utility

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"reflect"
	"testing"
)

func TestNewStreamMQ(t *testing.T) {

	client := redis.NewClient(&redis.Options{Addr: "127.0.0.1:26379"})
	type args struct {
		client *redis.Client
		maxLen int
		approx bool
	}
	tests := []struct {
		name string
		args args
		want *StreamMQ
	}{
		{
			name: "new",
			args: args{
				client: client,
				maxLen: 100,
				approx: true,
			},
			want: &StreamMQ{
				client: nil,
				maxLen: 100,
				approx: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewStreamMQ(tt.args.client, tt.args.maxLen, tt.args.approx)

			if got.client == nil {
				t.Errorf("%v", got.client)
			}

			if !reflect.DeepEqual(got.approx, tt.want.approx) {
				t.Errorf("NewStreamMQ() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStreamMQ_Consume(t *testing.T) {
	client := redis.NewClient(&redis.Options{Addr: "127.0.0.1:26379"})
	type fields struct {
		client *redis.Client
		maxLen int64
		approx bool
	}
	type args struct {
		ctx       context.Context
		topic     string
		group     string
		consumer  string
		start     string
		batchSize int
		h         Handler
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "消费",
			fields: fields{
				client: client,
				maxLen: 100,
				approx: true,
			},
			args: args{
				ctx:       context.TODO(),
				topic:     "email",
				group:     "",
				consumer:  "",
				start:     "1706256914275-0",
				batchSize: 10,
				h: func(msg *Msg) error {
					fmt.Println(string(msg.Body))
					return nil
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &StreamMQ{
				client: tt.fields.client,
				maxLen: tt.fields.maxLen,
				approx: tt.fields.approx,
			}

			// 阻塞
			if err := q.Consume(tt.args.ctx, tt.args.topic, tt.args.group, tt.args.consumer, tt.args.start, tt.args.batchSize, tt.args.h); (err != nil) != tt.wantErr {
				t.Errorf("Consume() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStreamMQ_SendMsg(t *testing.T) {

	client := redis.NewClient(&redis.Options{Addr: "127.0.0.1:26379"})

	type fields struct {
		client *redis.Client
		maxLen int64
		approx bool
	}
	type args struct {
		ctx context.Context
		msg *Msg
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "发送消息",
			fields: fields{
				client: client,
				maxLen: 100,
				approx: true,
			},
			args: args{
				ctx: context.Background(),
				msg: &Msg{
					ID:        "0",
					Topic:     "email",
					Body:      []byte("xfzheng@163.com"),
					Partition: 1,
					Group:     "1",
					Consumer:  "1",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &StreamMQ{
				client: tt.fields.client,
				maxLen: tt.fields.maxLen,
				approx: tt.fields.approx,
			}
			if err := q.SendMsg(tt.args.ctx, tt.args.msg); (err != nil) != tt.wantErr {
				t.Errorf("SendMsg() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
