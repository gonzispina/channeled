package context

import (
	"context"
	"time"

	"github.com/gonzispina/channeled/kit/uuid"
)

// Upgrade from the standard's lib context to our context
func Upgrade(ctx context.Context) Context {
	return &c{Context: ctx, trackingID: uuid.New()}
}

// Background context
func Background() Context {
	return &c{Context: context.Background(), trackingID: uuid.New()}
}

// WithTrackingID context
func WithTrackingID(id string) Context {
	if id == "" {
		id = uuid.New()
	}
	return &c{Context: context.Background(), trackingID: id}
}

// CancelFunc shadow
type CancelFunc = context.CancelFunc

// WithValue context
func WithValue(ctx Context, key interface{}, value interface{}) Context {
	n := context.WithValue(ctx, key, value)
	return &c{Context: n, trackingID: ctx.TrackingID()}
}

// WithCancel context
func WithCancel(ctx Context) (Context, CancelFunc) {
	n, f := context.WithCancel(ctx)
	return &c{Context: n, trackingID: ctx.TrackingID()}, f
}

// WithTimeout context
func WithTimeout(ctx Context, d time.Duration) (Context, context.CancelFunc) {
	n, f := context.WithTimeout(ctx, d)
	return &c{Context: n, trackingID: ctx.TrackingID()}, f
}

// Merge context
func Merge(newCtx context.Context, oldCtx Context) Context {
	ctx := oldCtx.(*c)
	return &c{Context: newCtx, trackingID: ctx.trackingID}
}

// Context interface
type Context interface {
	context.Context
	TrackingID() string
}

type c struct {
	context.Context
	trackingID string
}

// TrackingID to track everything
func (c c) TrackingID() string {
	return c.trackingID
}

// Value of the context
func (c c) Value(key interface{}) interface{} {
	return c.Context.Value(key)
}
