package utility

import (
	"testing"
)

func TestRMQ_Consume(t *testing.T) {
	//
	//client, _ := amqp.Dial("amqp://guest:guest@localhost:5672/")
	//
	//type fields struct {
	//	client *amqp.Connection
	//	logger *log.Logger
	//	err    error
	//}
	//type args struct {
	//	name string
	//}
	//tests := []struct {
	//	name   string
	//	fields fields
	//	args   args
	//}{
	//	{
	//		name: "消费消息",
	//		fields: fields{
	//			client: client,
	//			logger: log.New(os.Stdout, "log-", 0777),
	//			err:    errors.New("err"),
	//		},
	//		args: args{
	//			name: "orders",
	//		},
	//	},
	//}
	//
	//for _, tt := range tests {
	//	t.Run(tt.name, func(t *testing.T) {
	//		r := &RMQ{
	//			client: tt.fields.client,
	//			logger: tt.fields.logger,
	//			err:    tt.fields.err,
	//		}
	//		r.Consume(tt.args.name)
	//	})
	//}
}

func TestRMQ_Publish(t *testing.T) {
	//
	//client, _ := amqp.Dial("amqp://guest:guest@localhost:5672/")
	//
	//type fields struct {
	//	client *amqp.Connection
	//	logger *log.Logger
	//	err    error
	//}
	//type args struct {
	//	name string
	//	msg  []byte
	//}
	//tests := []struct {
	//	name   string
	//	fields fields
	//	args   args
	//}{
	//	{
	//		name: "生产订单",
	//		fields: fields{
	//			client: client,
	//			logger: log.New(os.Stdout, "log-", 0777),
	//			err:    errors.New("err"),
	//		},
	//		args: args{
	//			name: "orders",
	//			msg:  []byte("KPD000000001"),
	//		},
	//	},
	//}
	//
	//for _, tt := range tests {
	//	t.Run(tt.name, func(t *testing.T) {
	//		r := &RMQ{
	//			client: tt.fields.client,
	//			logger: tt.fields.logger,
	//			err:    tt.fields.err,
	//		}
	//		r.Publish(tt.args.name, tt.args.msg)
	//	})
	//}
}
