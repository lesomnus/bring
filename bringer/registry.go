package bringer

type Factory = func(opts ...Option) Bringer

var bringers = map[string](Factory){
	"file":  FileBringer,
	"http":  HttpBringer,
	"https": HttpBringer,
	"smb":   SmbBringer,
}
