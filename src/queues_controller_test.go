package main

import (
	"testing"
)

func TestQueues_Pop(t *testing.T) {
	type fields struct {
		maxQueuesCount int
		maxQueueSize   int
		queues         map[string][]string
	}
	type args struct {
		name string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "test1",
			fields: fields{
				maxQueuesCount: 1,
				maxQueueSize:   1,
				queues: map[string][]string{
					"test": []string{"AAAAAAAAAAAAAAAA"},
				},
			},
			args: args{
				name: "test",
			},
			want: "AAAAAAAAAAAAAAAA",
		},
		{
			name: "test2",
			fields: fields{
				maxQueuesCount: 1,
				maxQueueSize:   1,
				queues: map[string][]string{
					"test": []string{"test", "AAAAAAAAAAAAAAAA"},
				},
			},
			args: args{
				name: "test",
			},
			want: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &QueuesController{
				maxQueuesCount: tt.fields.maxQueuesCount,
				maxQueueSize:   tt.fields.maxQueueSize,
				queues:         tt.fields.queues,
			}
			if got, _ := q.Pop(nil, tt.args.name); got != tt.want {
				t.Errorf("Pop() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQueues_Push(t *testing.T) {
	type fields struct {
		maxQueuesCount int
		maxQueueSize   int
		queues         map[string][]string
	}
	type args struct {
		name  string
		value string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "test1",
			fields: fields{
				maxQueueSize:   1,
				maxQueuesCount: 1,
				queues:         map[string][]string{},
			},
			args: args{
				name:  "test",
				value: "AAAAAAAAAA",
			},
			wantErr: false,
		},
		{
			name: "test2",
			fields: fields{
				maxQueueSize:   1,
				maxQueuesCount: 1,
				queues: map[string][]string{
					"test": []string{"existElement"},
				},
			},
			args: args{
				name:  "test",
				value: "AAAAAAAAAA",
			},
			wantErr: true,
		},
		{
			name: "test3",
			fields: fields{
				maxQueueSize:   1,
				maxQueuesCount: 1,
				queues: map[string][]string{
					"test": []string{"existElement"},
				},
			},
			args: args{
				name:  "NE_TEST",
				value: "AAAAAAAAAA",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &QueuesController{
				maxQueuesCount: tt.fields.maxQueuesCount,
				maxQueueSize:   tt.fields.maxQueueSize,
				queues:         tt.fields.queues,
			}
			if err := q.Push(tt.args.name, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Push() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
