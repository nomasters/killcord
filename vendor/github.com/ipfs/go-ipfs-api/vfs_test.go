package shell

import (
	"crypto/md5"
	"fmt"
	"io"
	"testing"

	"github.com/cheekybits/is"
	"github.com/shurcooL/go/vfs/httpfs/vfsutil"
)

func TestVFS_Open(t *testing.T) {
	is := is.New(t)
	s := NewShell(shellUrl)

	vfs := NewVFS("/ipfs/"+examplesHash, s)

	f, err := vfs.Open("readme")
	is.Nil(err)

	st, err := f.Stat()
	is.Nil(err)
	is.False(st.IsDir())

	md5 := md5.New()
	_, err = io.Copy(md5, f)
	is.Nil(err)
	is.Equal(fmt.Sprintf("%x", md5.Sum(nil)), "3fdcaad186e79983a6920b4c7eeda949")
	is.Nil(f.Close())
}

func TestVFS_OpenDir(t *testing.T) {
	is := is.New(t)
	s := NewShell(shellUrl)

	vfs := NewVFS("/ipfs/"+examplesHash, s)

	f, err := vfs.Open(".")
	is.Nil(err)

	st, err := f.Stat()
	is.Nil(err)
	is.True(st.IsDir())
	is.Nil(f.Close())

	files, err := vfsutil.ReadDir(vfs, ".")
	is.Nil(err)

	is.Equal(len(files), 6)

	expected := map[string]UnixLsLink{
		"about":          {Type: "File", Hash: "QmZTR5bcpQD7cFgTorqxZDYaew1Wqgfbd2ud9QqGPAkK2V", Name: "about", Size: 1677},
		"contact":        {Type: "File", Hash: "QmYCvbfNbCwFR45HiNP45rwJgvatpiW38D961L5qAhUM5Y", Name: "contact", Size: 189},
		"help":           {Type: "File", Hash: "QmY5heUM5qgRubMDD1og9fhCPA6QdkMp3QCwd4s7gJsyE7", Name: "help", Size: 311},
		"quick-start":    {Type: "File", Hash: "QmUzLxaXnM8RYCPEqLDX5foToi5aNZHqfYr285w2BKhkft", Name: "quick-start", Size: 1686},
		"readme":         {Type: "File", Hash: "QmPZ9gcCEpqKTo6aq61g2nXGUhM4iCL3ewB6LDXZCtioEB", Name: "readme", Size: 1091},
		"security-notes": {Type: "File", Hash: "QmTumTjvcYCAvRRwQ8sDRxh8ezmrcr88YFU7iYNroGGTBZ", Name: "security-notes", Size: 1016},
	}
	for _, f := range files {
		el, ok := expected[f.Name()]
		is.True(ok)
		is.NotNil(el)
		is.Equal(f.Size(), el.Size)
		is.Equal(f.IsDir(), false)
	}
}
