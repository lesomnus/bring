package smb

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"strings"
	"time"

	"github.com/hirochachacha/go-smb2"
	"github.com/lesomnus/bring/bringer"
	"github.com/lesomnus/bring/bringer/option"
	"github.com/lesomnus/bring/log"
	"github.com/lesomnus/bring/thing"
)

type br struct {
	password    *string
	dialTimeout time.Duration
}

func (c *br) apply(opts []bringer.Option) {
	for _, opt := range opts {
		switch o := opt.(type) {
		case (*option.PwOption):
			c.password = &o.Value
		case (*option.DialTimeoutOpt):
			c.dialTimeout = o.Value
		}
	}
}

func (b *br) with(opts []bringer.Option) *br {
	if len(opts) == 0 {
		return b
	}

	b_ := *b
	b_.apply(opts)
	return &b_
}

func SmbBringer(opts ...bringer.Option) bringer.Bringer {
	b := &br{}
	b.apply(opts)
	return b
}

type file struct {
	*smb2.File
	l *slog.Logger

	share   *smb2.Share
	session *smb2.Session
	conn    net.Conn
}

func (f *file) Close() error {
	if f.conn == nil {
		return nil
	}

	f.l.Debug("close SMB file")
	err := errors.Join(
		f.share.Umount(),
		f.session.Logoff(),
		f.conn.Close(),
	)

	f.conn = nil
	return err
}

func (b *br) bring(ctx context.Context, t thing.Thing, opts ...bringer.Option) (v *file, err error) {
	l := log.From(ctx).With(slog.String("bringer", "smb"))
	b = b.with(opts)
	// TODO: connection pool? session pool?

	v = &file{l: l}

	host := t.Url.Host
	if !strings.Contains(host, ":") {
		// Add default port number
		host += ":445"
		l.Debug("use default por number")
	}

	ctx_dial := ctx
	if b.dialTimeout != 0 {
		ctx_, cancel := context.WithTimeout(ctx, b.dialTimeout)
		defer cancel()
		ctx_dial = ctx_
	}

	l.Info("dial TCP", slog.String("host", host))
	{
		d := net.Dialer{}
		v.conn, err = d.DialContext(ctx_dial, "tcp", host)
		if err != nil {
			e := &net.OpError{}
			if errors.As(err, &e) {
				return v, err
			}
			return v, fmt.Errorf("dial TCP: %w", err)
		}
	}

	username := t.Url.User.Username()
	password := ""
	if b.password != nil {
		password = *b.password
	}
	if v, ok := t.Url.User.Password(); ok {
		password = v
	}

	share, p := splitPath(t.Url.Path)

	l.Info("dial SMB",
		slog.String("username", username),
		slog.Bool("password", password != ""),
	)
	{
		d := &smb2.Dialer{
			Initiator: &smb2.NTLMInitiator{
				User:     username,
				Password: password,
			},
		}
		v.session, err = d.DialContext(ctx_dial, v.conn)
		if err != nil {
			return v, fmt.Errorf("dial SMB: %w", err)
		}
	}

	l.Info("mount", slog.String("share", share))
	v.share, err = v.session.Mount(share)
	if err != nil {
		return v, fmt.Errorf("mount SMB share %s: %w", share, err)
	}

	l.Info("open", slog.String("path", p))
	v.share = v.share.WithContext(ctx)
	v.File, err = v.share.Open(p)
	if err != nil {
		return v, fmt.Errorf("open: %w", err)
	}

	return v, nil
}

func (b *br) Bring(ctx context.Context, t thing.Thing, opts ...bringer.Option) (io.ReadCloser, error) {
	f, err := b.bring(ctx, t, opts...)
	if err != nil {
		errs := []error{err}
		if f.share != nil {
			errs = append(errs, f.share.Umount())
		}
		if f.session != nil {
			errs = append(errs, f.session.Logoff())
		}
		if f.conn != nil {
			errs = append(errs, f.conn.Close())
		}
		return nil, errors.Join(errs...)
	}

	return f, nil
}

// Splits `/share_name/filepath` into `[share_name, filepath]`.
func splitPath(p string) (string, string) {
	p, _ = strings.CutPrefix(p, "/")
	ps := strings.SplitN(p, "/", 2)
	switch len(ps) {
	case 1:
		return ps[0], ""
	case 2:
		return ps[0], ps[1]
	default:
		panic("splint into at most 2")
	}
}

func init() {
	bringer.Register("smb", SmbBringer)
}
