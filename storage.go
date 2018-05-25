package killcord

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/ipfs/go-ipfs-api"
)

const (
	defaultpayloadRPCPath = "https://ipfs.infura.io:5001"
)

var (
	payloadSourcePath    string
	payloadEncryptedPath string
	payloadTempPath      string
	payloadDecryptPath   string
	payloadRPCPath       string = defaultpayloadRPCPath
)

func init() {
	payloadSourcePath = filepath.Join(strings.Split(sourcePrefix, "/")...)
	payloadEncryptedPath = filepath.Join(strings.Split(encryptPrefix, "/")...)
	payloadTempPath = filepath.Join(strings.Split(tempPrefix, "/")...)
	payloadDecryptPath = filepath.Join(strings.Split(decryptPrefix, "/")...)
}

// getShellURL waterfalls settings and returns the RPC url in priority
// order from Options, Config, or Default settings

func (s *Session) setPayloadRPCPath() {
	if s.Options.Payload.RPCURL != "" {
		payloadRPCPath = s.Options.Payload.RPCURL
		return
	}
	if s.Config.Payload.RPCURL != "" {
		payloadRPCPath = s.Config.Payload.RPCURL
		return
	}
	payloadRPCPath = defaultpayloadRPCPath
}

// DeployPayload takes the contents from the /payload/encrypted and adds it to
// the storage endpoint

func (s *Session) DeployPayload() error {
	// check that encrypted payload exists, exit if it doesn't
	if _, err := os.Stat(filepath.Join(payloadEncryptedPath, defaultOutputKilName)); os.IsNotExist(err) {
		return errors.New("encrypted payload does not exist, exiting")
	}
	// check for payload ID in config, exti if it already exists
	if s.Config.Payload.ID != "" {
		return fmt.Errorf("payload %v already deployed, skipping", s.Config.Payload.ID)
	}
	sh := shell.NewShell(payloadRPCPath)
	f, err := os.Open(filepath.Join(payloadEncryptedPath, defaultOutputKilName))
	if err != nil {
		return err
	}
	mhash, err := sh.Add(f)
	if err != nil {
		return err
	}
	s.Config.Payload.ID = mhash
	if err := SetPayloadEndpoint(s.Config.Contract.Owner, s.Config.Contract.ID, s.Config.Payload.ID); err != nil {
		return err
	}
	s.Config.Payload.Status = "deployed"
	return nil
}

// Gets the payload from the storage endpoint and stores it locally
// in the Encrypted payload folder
func (s *Session) GetPayload() error {
	sh := shell.NewShell(payloadRPCPath)
	if err := sh.Get(s.Config.Payload.ID, payloadEncryptedPath); err != nil {
		return err
	}
	s.Config.Payload.Status = "synced"
	if err := os.Rename(filepath.Join(payloadEncryptedPath, s.Config.Payload.ID), filepath.Join(payloadEncryptedPath, defaultOutputKilName)); err != nil {
		return err
	}
	return nil
}

func (s *Session) Encrypt() error {
	var key [32]byte
	os.RemoveAll(payloadTempPath)
	if err := os.Mkdir(payloadTempPath, 0755); err != nil {
		return err
	}
	if err := zipSource(); err != nil {
		return err
	}
	if err := s.setPayloadKey(); err != nil {
		return err
	}
	secret, err := hex.DecodeString(s.Config.Payload.Secret)
	if err != nil {
		return err
	}
	copy(key[:], secret)
	encryptMultiPart(key)
	os.RemoveAll(payloadTempPath)
	s.Config.Payload.Status = "encrypted"

	return nil
}

func (s *Session) Decrypt() error {
	var key [32]byte
	secret, err := hex.DecodeString(s.Config.Payload.Secret)
	if err != nil {
		return err
	}
	copy(key[:], secret)
	os.RemoveAll(payloadTempPath)
	os.Mkdir(payloadTempPath, 0755)
	decryptMultiPart(key)
	if err := unzipSource(); err != nil {
		return err
	}
	os.RemoveAll(payloadDecryptPath)
	if err := os.Rename(filepath.Join(payloadTempPath, "source"), payloadDecryptPath); err != nil {
		return err
	}
	os.RemoveAll(payloadTempPath)
	return nil
}

func (s *Session) setPayloadKey() error {
	if s.Config.Payload.Secret != "" {
		return errors.New("encryption secret already set")
	}
	s.Config.Payload.Secret = generateKey()
	return nil
}

func initFile(file string) *os.File {
	i, err := os.Create(file)
	if err != nil {
		log.Fatal(err)
	}
	i.Close()
	f, err := os.OpenFile(file, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	return f
}

func encryptMultiPart(k [32]byte) {
	s := filepath.Join(payloadTempPath, defaultOutputZipName)
	d := filepath.Join(payloadEncryptedPath, defaultOutputKilName)
	totalSize, _ := getFileSize(d)
	f := initFile(d)
	r := make(chan Chunk)
	w := make(chan Chunk)
	go reader(s, r)
	go encrypter(r, w, k, totalSize)
	writer(w, f)
}

func decryptMultiPart(k [32]byte) {
	s := filepath.Join(payloadEncryptedPath, defaultOutputKilName)
	d := filepath.Join(payloadTempPath, defaultOutputZipName)
	totalSize, _ := getFileSize(s)
	f := initFile(d)
	r := make(chan Chunk)
	w := make(chan Chunk)
	go reader(s, r)
	go decrypter(r, w, k, totalSize)
	writer(w, f)
}

func encrypter(r, w chan Chunk, k [32]byte, totalSize int64) {
	var block []byte
	count := int((totalSize / maxChunkSize) + 1)
	fmt.Println("-- encrypting payload --")
	bar := pb.StartNew(count)
	for {
		x, ok := <-r
		if len(r) == 0 {
			if !ok {
				break
			}
		}
		for _, b := range x {
			if len(block) < maxChunkSize {
				block = append(block, b)
			}
			if len(block) == maxChunkSize {
				e := encrypt(block, k)
				bar.Increment()
				w <- e
				block = []byte{}
			}
		}
	}
	e := encrypt(block, k)
	w <- e
	bar.Set(count)
	bar.Finish()
	close(w)
}

func decrypter(r, w chan Chunk, k [32]byte, totalSize int64) {
	var block []byte
	var chunkSize = maxChunkSize + 40
	count := int((totalSize / maxChunkSize) + 1)
	fmt.Println("-- decrypting payload --")
	bar := pb.StartNew(count)
	for {
		x, ok := <-r
		if len(r) == 0 {
			if !ok {
				break
			}
		}
		for _, b := range x {
			if len(block) < chunkSize {
				block = append(block, b)
			}
			if len(block) == chunkSize {
				e := decrypt(block, k)
				bar.Increment()
				w <- e
				block = []byte{}
			}
		}
	}
	e := decrypt(block, k)
	w <- e
	bar.Set(count)
	bar.Finish()
	close(w)
}

func reader(file string, r chan Chunk) {
	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	fi, err := f.Stat()
	if err != nil {
		log.Fatal(err)
	}
	remaining := fi.Size()
	var offset int64

	for {
		if len(r) == 0 {
			if remaining <= maxReadSize {
				rc := make([]byte, remaining)
				f.ReadAt(rc, offset)
				r <- rc
				f.Close()
				close(r)
				return
			} else {
				rc := make([]byte, maxReadSize)
				f.ReadAt(rc, offset)
				r <- rc
				remaining = remaining - maxReadSize
				offset = offset + maxReadSize
			}
		}
	}

}

func writer(w chan Chunk, f *os.File) {
	var writeCache []byte
	for {
		x, ok := <-w
		if len(w) == 0 {
			if !ok {
				writeToFile(writeCache, f)
				f.Close()
				return
			}
		}
		for _, b := range x {
			if len(writeCache) < maxChunksCache {
				writeCache = append(writeCache, b)
			}
			if len(writeCache) == maxChunksCache {
				writeToFile(writeCache, f)
				writeCache = []byte{}
			}
		}
	}
}

func writeToFile(data []byte, f *os.File) {
	if _, err := f.Write(data); err != nil {
		log.Fatal(err)
	}
}

func encrypt(data []byte, key [32]byte) []byte {
	var nonce [24]byte
	if _, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
		panic(err)
	}
	return secretbox.Seal(nonce[:], data, &nonce, &key)
}

func decrypt(data []byte, key [32]byte) []byte {
	var nonce [24]byte
	copy(nonce[:], data[:24])
	d, ok := secretbox.Open(nil, data[24:], &nonce, &key)
	if !ok {
		panic("decryption error")
	}
	return d
}
