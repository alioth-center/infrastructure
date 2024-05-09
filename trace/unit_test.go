package trace

import (
	"context"
	"fmt"
	"net"
	"net/netip"
	"testing"

	"google.golang.org/grpc/peer"
)

func TestStack(t *testing.T) {
	for i := 0; i < 10; i++ {
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			if i == 9 {
				t.Log("\n" + string(Stack(0)))
			}
		})
	}
}

func TestContext(t *testing.T) {
	t.Run("SetTraceIDKey", func(t *testing.T) {
		t.Log(ContextKey())
		SetTraceIDKey("tid")
		if traceIDKey != "tid" {
			t.Error("trace id key is not equal")
		}
		SetTraceIDKey("trace_id")
		if traceIDKey != "tid" {
			t.Error("trace id key is not equal")
		}
	})

	t.Run("TransformContext", func(t *testing.T) {
		ctx := NewContextWithTid("114514")
		tid, newCtx := TransformContext(ctx)
		if tid != "114514" {
			t.Error("trace id is not equal")
		}
		if newCtx != ctx {
			t.Error("context is not equal")
		}
	})

	t.Run("GetTid", func(t *testing.T) {
		ctx := NewContextWithTid("114514")
		if GetTid(ctx) != "114514" {
			t.Error("trace id is not equal")
		}
		if GetTid(context.Background()) != "" {
			t.Error("trace id is not empty")
		}
		if GetTid(context.WithValue(context.Background(), traceIDKey, map[string]string{})) != "" {
			t.Error("trace id is not empty")
		}
	})

	t.Run("FormContext", func(t *testing.T) {
		t.Run("Create", func(t *testing.T) {
			ctx := context.Background()
			traced := FromContext(ctx)
			if GetTid(traced) == "" {
				t.Error("trace id is empty")
			}
			if traced == ctx {
				t.Error("context is equal")
			}
		})

		t.Run("Existed", func(t *testing.T) {
			ctx := NewContext()
			traced := FromContext(ctx)
			if GetTid(traced) != GetTid(ctx) {
				t.Error("trace id is not equal")
			}
			if traced != ctx {
				t.Error("context is not equal")
			}
		})
	})

	t.Run("ForkContext", func(t *testing.T) {
		ctx := NewContext()
		ctx = context.WithValue(ctx, "key", "value")
		tid := GetTid(ctx)
		newCtx := ForkContext(ctx)
		if GetTid(newCtx) != tid {
			t.Error("trace id is not equal")
		}
		if newCtx.Value("key") != nil {
			t.Error("value is not nil")
		}
	})

	t.Run("ForkContextWithOpts", func(t *testing.T) {
		ctx := NewContext()
		ctx = context.WithValue(ctx, "key", "value")
		tid := GetTid(ctx)
		newCtx := ForkContextWithOpts(ctx, "key")
		if GetTid(newCtx) != tid {
			t.Error("trace id is not equal")
		}
		if newCtx.Value("key") != "value" {
			t.Error("value is not equal")
		}
	})

	t.Run("NewContext", func(t *testing.T) {
		ctx := NewContext()
		tid := GetTid(ctx)
		if tid == "" {
			t.Error("trace id is empty")
		}
	})

	t.Run("NewContextWithTid", func(t *testing.T) {
		tid := "114514"
		ctx := NewContextWithTid(tid)
		if GetTid(ctx) != tid {
			t.Error("trace id is not equal")
		}
	})

	t.Run("AttachTraceID", func(t *testing.T) {
		ctx := context.Background()
		tid := ""
		tid, ctx = AttachTraceID(ctx)
		if tid == "" {
			t.Error("trace id is empty")
		}
	})

	t.Run("Context", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, "key", "value")
		tid := "114514"
		newCtx := Context(ctx, tid)
		if GetTid(newCtx) != "114514" {
			t.Error("trace id is not equal")
		}
		if newCtx.Value("key") != "value" {
			t.Error("value is not equal")
		}
	})

	t.Run("GetClientIPFromPeer", func(t *testing.T) {
		t.Run("ConvertNotSuccess", func(t *testing.T) {
			ctx := context.Background()
			if GetClientIPFromPeer(ctx) != "" {
				t.Error("client ip is not empty")
			}
		})

		t.Run("NilAddr", func(t *testing.T) {
			ctx := peer.NewContext(context.Background(), &peer.Peer{})
			if GetClientIPFromPeer(ctx) != "" {
				t.Error("client ip is not empty")
			}
		})

		t.Run("FailedToParse", func(t *testing.T) {
			ctx := peer.NewContext(context.Background(), &peer.Peer{
				Addr: net.Addr(&net.IPAddr{}),
			})
			if GetClientIPFromPeer(ctx) != "" {
				t.Error("client ip is not empty")
			}
		})

		t.Run("Success", func(t *testing.T) {
			addr, _ := netip.ParseAddr("114.114.114.114")
			ctx := peer.NewContext(context.Background(), &peer.Peer{
				Addr: net.Addr(net.TCPAddrFromAddrPort(netip.AddrPortFrom(addr, 80))),
			})
			if GetClientIPFromPeer(ctx) != "114.114.114.114" {
				t.Error("client ip is not equal")
			}
		})
	})
}
