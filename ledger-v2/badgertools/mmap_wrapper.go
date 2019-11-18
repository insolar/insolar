//
//    Copyright 2019 Insolar Technologies
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.
//
//
//
//   Includes elements from Badger, licensed under Apache License, Version 2.0
//   Copyright 2017 Dgraph Labs, Inc. and Contributors
//

package badgertools

import (
	"fmt"
	"github.com/pkg/errors"
	"math"
	"os"
	"sync"
)

// FileLoadingMode specifies how data in LSM table files and value log files should
// be loaded.
type FileLoadingMode int

const (
	// FileIO indicates that files must be loaded using standard I/O
	FileIO FileLoadingMode = iota
	// LoadToRAM indicates that file must be loaded into RAM
	LoadToRAM
	// MemoryMap indicates that that the file must be memory-mapped
	MemoryMap
)

type logFile struct {
	path string
	// This is a lock on the log file. It guards the fd’s value, the file’s
	// existence and the file’s memory map.
	//
	// Use shared ownership when reading/writing the file or memory map, use
	// exclusive ownership to open/close the descriptor, unmap or remove the file.
	lock sync.RWMutex
	fd   *os.File
	//fid         uint32
	fmap        []byte
	size        uint32
	loadingMode FileLoadingMode
}

// openReadOnly assumes that we have a write lock on logFile.
func (lf *logFile) openReadOnly() error {
	var err error
	lf.fd, err = os.OpenFile(lf.path, os.O_RDONLY, 0666)
	if err != nil {
		return errors.Wrapf(err, "Unable to open %q as RDONLY.", lf.path)
	}

	fi, err := lf.fd.Stat()
	if err != nil {
		return errors.Wrapf(err, "Unable to check stat for %q", lf.path)
	}
	AssertTrue(fi.Size() <= math.MaxUint32)
	lf.size = uint32(fi.Size())

	if err = lf.mmap(fi.Size()); err != nil {
		_ = lf.fd.Close()
		return Wrapf(err, "Unable to map file: %q", fi.Name())
	}

	return nil
}

func (lf *logFile) mmap(size int64) (err error) {
	if lf.loadingMode != MemoryMap {
		// Nothing to do
		return nil
	}
	lf.fmap, err = Mmap(lf.fd, false, size)
	if err == nil {
		err = Madvise(lf.fmap, false) // Disable readahead
	}
	return err
}

func (lf *logFile) munmap() (err error) {
	if lf.loadingMode != MemoryMap {
		// Nothing to do
		return nil
	}
	if err := Munmap(lf.fmap); err != nil {
		return errors.Wrapf(err, "Unable to munmap value log: %q", lf.path)
	}
	return nil
}

// Acquire lock on mmap/file if you are calling this
func (lf *logFile) read(p valuePointer) (buf []byte, err error) {
	var nbr int64
	offset := p.Offset
	if lf.loadingMode == FileIO {
		buf = s.Resize(int(p.Len))
		var n int
		n, err = lf.fd.ReadAt(buf, int64(offset))
		nbr = int64(n)
	} else {
		// Do not convert size to uint32, because the lf.fmap can be of size
		// 4GB, which overflows the uint32 during conversion to make the size 0,
		// causing the read to fail with ErrEOF. See issue #585.
		size := int64(len(lf.fmap))
		valsz := p.Len
		if int64(offset) >= size || int64(offset+valsz) > size {
			err = ErrEOF
		} else {
			buf = lf.fmap[offset : offset+valsz]
			nbr = int64(valsz)
		}
	}
	//y.NumReads.Add(1)
	//y.NumBytesRead.Add(nbr)
	return buf, err
}

func (lf *logFile) doneWriting(offset uint32) error {
	// Sync before acquiring lock.  (We call this from write() and thus know we have shared access
	// to the fd.)
	if err := FileSync(lf.fd); err != nil {
		return errors.Wrapf(err, "Unable to sync value log: %q", lf.path)
	}
	// Close and reopen the file read-only.  Acquire lock because fd will become invalid for a bit.
	// Acquiring the lock is bad because, while we don't hold the lock for a long time, it forces
	// one batch of readers wait for the preceding batch of readers to finish.
	//
	// If there's a benefit to reopening the file read-only, it might be on Windows.  I don't know
	// what the benefit is.  Consider keeping the file read-write, or use fcntl to change
	// permissions.
	lf.lock.Lock()
	defer lf.lock.Unlock()
	if err := lf.munmap(); err != nil {
		return err
	}
	// TODO: Confirm if we need to run a file sync after truncation.
	// Truncation must run after unmapping, otherwise Windows would crap itself.
	if err := lf.fd.Truncate(int64(offset)); err != nil {
		return errors.Wrapf(err, "Unable to truncate file: %q", lf.path)
	}
	if err := lf.fd.Close(); err != nil {
		return errors.Wrapf(err, "Unable to close value log: %q", lf.path)
	}

	return lf.openReadOnly()
}

// You must hold lf.lock to sync()
func (lf *logFile) sync() error {
	return FileSync(lf.fd)
}

// AssertTrue asserts that b is true. Otherwise, it would log fatal.
func AssertTrue(b bool) {
	if !b {
		panic(fmt.Sprintf("%+v", errors.Errorf("Assert failed")))
	}
	//if !b {
	//	log.Fatalf("%+v", errors.Errorf("Assert failed"))
	//}
}

// Wrap wraps errors from external lib.
func Wrap(err error) error {
	//if !debugMode {
	//	return err
	//}
	return errors.Wrap(err, "")
}

// Wrapf is Wrap with extra info.
func Wrapf(err error, format string, args ...interface{}) error {
	//if !debugMode {
	//	if err == nil {
	//		return nil
	//	}
	//	return fmt.Errorf(format+" error: %+v", append(args, err)...)
	//}
	return errors.Wrapf(err, format, args...)
}
