package chezmoi

import (
	"errors"
	"io/fs"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/twpayne/chezmoi/v2/pkg/archive"
)

func TestArchiveReaderSystemTar(t *testing.T) {
	data, err := archive.NewTar(map[string]interface{}{
		"dir": map[string]interface{}{
			"file": "# contents of dir/file\n",
			"symlink": &archive.Symlink{
				Target: "file",
			},
		},
	})
	require.NoError(t, err)

	archiveReaderSystem, err := NewArchiveReaderSystem("archive.tar", data, ArchiveFormatTar, ArchiveReaderSystemOptions{
		RootAbsPath:     NewAbsPath("/home/user"),
		StripComponents: 1,
	})
	assert.NoError(t, err)

	for _, tc := range []struct {
		absPath      AbsPath
		lstatErr     error
		readlink     string
		readlinkErr  error
		readFileData []byte
		readFileErr  error
	}{
		{
			absPath:      NewAbsPath("/home/user/file"),
			readlinkErr:  fs.ErrInvalid,
			readFileData: []byte("# contents of dir/file\n"),
		},
		{
			absPath:     NewAbsPath("/home/user/notexist"),
			readlinkErr: fs.ErrNotExist,
			lstatErr:    fs.ErrNotExist,
			readFileErr: fs.ErrNotExist,
		},
		{
			absPath:     NewAbsPath("/home/user/symlink"),
			readlink:    "file",
			readFileErr: fs.ErrInvalid,
		},
	} {
		_, err = archiveReaderSystem.Lstat(tc.absPath)
		if tc.lstatErr != nil {
			assert.True(t, errors.Is(err, tc.lstatErr))
		} else {
			assert.NoError(t, err)
		}

		actualLinkname, err := archiveReaderSystem.Readlink(tc.absPath)
		if tc.readlinkErr != nil {
			assert.True(t, errors.Is(err, tc.readlinkErr))
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tc.readlink, actualLinkname)
		}

		actualReadFileData, err := archiveReaderSystem.ReadFile(tc.absPath)
		if tc.readFileErr != nil {
			assert.True(t, errors.Is(err, tc.readFileErr))
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tc.readFileData, actualReadFileData)
		}
	}
}
