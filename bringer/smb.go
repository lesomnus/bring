package bringer

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"strings"

	"github.com/hirochachacha/go-smb2"
	"github.com/lesomnus/bring/log"
	"github.com/lesomnus/bring/thing"
)

type smbBringer struct {
	password string
}

func SmbBringer(opts ...Option) Bringer {
	b := &smbBringer{}
	for _, opt := range opts {
		switch o := opt.(type) {
		case (*pwOpt):
			b.password = o.v
		}
	}

	return b
}

type smbFile struct {
	*smb2.File
	share   *smb2.Share
	session *smb2.Session
	conn    net.Conn
}

func (f *smbFile) Close() error {
	return errors.Join(
		f.share.Umount(),
		f.session.Logoff(),
		f.conn.Close(),
	)
}

func (b *smbBringer) bring(ctx context.Context, t thing.Thing) (v *smbFile, err error) {
	l := log.From(ctx).With(name("smb"))
	// TODO: connection pool? session pool?

	v = &smbFile{}

	host := t.Url.Host
	if !strings.Contains(host, ":") {
		// Add default port number
		host += ":445"
		l.Debug("use default por number")
	}

	l.Info("dial TCP", slog.String("host", host))
	v.conn, err = net.Dial("tcp", host)
	if err != nil {
		return v, fmt.Errorf("dial TCP: %w", err)
	}

	username := t.Url.User.Username()
	password, _ := t.Url.User.Password()
	share, p := b.splitPath(t.Url.Path)

	d := &smb2.Dialer{
		Initiator: &smb2.NTLMInitiator{
			User:     username,
			Password: password,
		},
	}

	l.Info("dial SMB",
		slog.String("username", username),
		slog.Bool("password", password != ""),
	)
	v.session, err = d.Dial(v.conn)
	if err != nil {
		return v, fmt.Errorf("dial SMB: %w", err)
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
		return nil, fmt.Errorf("open: %w", err)
	}

	return v, nil
}

func (b *smbBringer) Bring(ctx context.Context, t thing.Thing) (io.ReadCloser, error) {
	f, err := b.bring(ctx, t)
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
func (b *smbBringer) splitPath(p string) (string, string) {
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
