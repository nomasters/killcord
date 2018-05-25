package local

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	//"gopkg.in/cheggaaa/pb.v1"
	"github.com/nomasters/killcord"
	"github.com/nomasters/killcord/payload"
	"golang.org/x/crypto/nacl/secretbox"
)

const (
	maxChunkSize         = 16000
	defaultEncryptSource = "payload/source/"
	defaultEncryptOutput = "payload/output.kil"
	defaultDecryptOutput = "output.zip"
)

type Options struct {
	Source      string
	Destination string
}

type Payload struct {
	Options Options
	Config  payload.Config
}

func NewPayload(opts Options, config payload.Config) Payload {
	p := new(Payload)
	p.setOptions(opts)
	return p
}

// check to see if Key exists, and if it does not generate a
// new key, return the key.
func (p *Payload) GetKey() string {
	if p.Config.Key == "" {
		p.Config.Key = killcord.GenerateKey()
	}
	return p.Config.Key
}

func (p Payload) GetConfig() payload.Config {
	return p.Config
}

func (p *Payload) SetConfig(config payload.Config) error {
	if config.Status != "" {
		p.Config.Status = config.Status
	}
	if config.Provider != "" {
		p.Config.Provider = config.Provider
	}
	if config.Key != "" {
		p.Config.Key = config.Key
	}
	for x := range config.Settings {
		p.Config.Settings[x] = config.Settings[x]
	}
	// TODO: sanity check the config
	return nil
}

func (p *Payload) setOptions(opts Options) {
	if s := opts.Source; s != "" {
		p.Settings["Source"] = s
	}
	if d := opts.Destination; d != "" {
		p.Settings["Destination"] = d
	}
}

func (p *Payload) Encrypt() error {
	source := defaultEncryptSource
	dest := defaultEncryptOutput

	if x, ok := p.Settings["Source"].(string); ok == true && x != "" {
		source = x
	}

	if x, ok := p.Settings["Destination"].(string); ok == true && x != "" {
		dest = x
	}

	r, w := io.Pipe()

	go func() {
		if err := zip(w, source); err != nil {
			log.Fatal(err)
		}
		if err := w.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	outputFile, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()
	defer r.Close()

	if err := encryptWriter(p.Key, outputFile, r); err != nil {
		log.Fatal(err)
	}
}

func (p *Payload) Decrypt() error {
	source := defaultEncryptOutput
	dest := defaultDecryptOutput

	if x, ok := p.Settings["Destination"].(string); ok == true && x != "" {
		source = x
	}

	r, w := io.Pipe()
	go func() {
		inputFile, err := os.OpenFile(source, os.O_RDONLY, 0666)
		if err != nil {
			log.Fatal(err)
		}
		defer inputFile.Close()
		if err := decryptWriter(p.Key, w, inputFile); err != nil {
			log.Fatal(err)
		}
		if err := w.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	defer r.Close()

	outputFile, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()
	buf := make([]byte, 256)
	if _, err := io.CopyBuffer(outputFile, r, buf); err != nil {
		log.Fatal(err)
	}
}

func decryptWriter(key [32]byte, output io.Writer, input io.Reader) error {
	var totalProcessed int64
	var totalRead int64
	buf := bytes.NewBuffer([]byte{})
	go func() {
		total, err := buf.ReadFrom(input)
		if err != nil {
			log.Fatal(err)
		}
		totalRead = total
	}()

	retry := 0
	maxRetry := 100
	maxSize := maxChunkSize + 40

	for {
		if totalProcessed == 0 {
			time.Sleep(100 * time.Millisecond)
		}
		if totalRead > 0 && totalRead == totalProcessed {
			break
		}
		if buf.Len() >= maxSize {
			chunk := buf.Next(maxSize)
			totalProcessed = totalProcessed + int64(len(chunk))
			output.Write(decrypt(chunk, key))
			continue
		}
		if totalRead > 0 && int64(buf.Len()) == (totalRead-totalProcessed) {
			chunk := buf.Next(maxSize)
			totalProcessed = totalProcessed + int64(len(chunk))
			output.Write(decrypt(chunk, key))
			continue
		} else {
			if retry < maxRetry {
				time.Sleep(50 * time.Millisecond)
				continue
			}
			return errors.New("buffer contention, hit max retries on encryption")
		}
	}
	return nil
}

func encryptWriter(key [32]byte, output io.Writer, input io.Reader) error {
	var totalProcessed int64
	var totalRead int64
	buf := bytes.NewBuffer([]byte{})
	go func() {
		total, err := buf.ReadFrom(input)
		if err != nil {
			log.Fatal(err)
		}
		totalRead = total
	}()

	retry := 0
	maxRetry := 100

	for {
		if totalProcessed == 0 {
			time.Sleep(100 * time.Millisecond)
		}
		if totalRead > 0 && totalRead == totalProcessed {
			break
		}
		if buf.Len() >= maxChunkSize {
			chunk := buf.Next(maxChunkSize)
			totalProcessed = totalProcessed + int64(len(chunk))
			output.Write(encrypt(chunk, key))
			continue
		}
		if totalRead > 0 && int64(buf.Len()) == (totalRead-totalProcessed) {
			chunk := buf.Next(maxChunkSize)
			totalProcessed = totalProcessed + int64(len(chunk))
			output.Write(encrypt(chunk, key))
			continue
		} else {
			if retry < maxRetry {
				time.Sleep(50 * time.Millisecond)
				continue
			}
			return errors.New("buffer contention, hit max retries on encryption")
		}
	}
	return nil
}

func zip(output io.Writer, sourcePath string) error {
	w := zip.NewWriter(output)
	defer w.Close()
	info, err := os.Stat(sourcePath)
	if err != nil {
		return fmt.Errorf("%s: stat: %v", sourcePath, err)
	}
	var base string
	if info.IsDir() {
		base = filepath.Base(sourcePath)
	}

	return filepath.Walk(sourcePath, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("walking to %s: %v", filePath, err)
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return fmt.Errorf("%s: getting header: %v", filePath, err)
		}

		if base != "" {
			name, err := filepath.Rel(sourcePath, filePath)
			if err != nil {
				return err
			}
			header.Name = path.Join(base, filepath.ToSlash(name))
		}

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}
		if strings.HasPrefix(header.Name, sourcePath) {
			header.Name = strings.Replace(header.Name, sourcePath, "/", 1)
		}
		f, err := w.CreateHeader(header)
		if err != nil {
			return fmt.Errorf("%s: making header: %v", filePath, err)
		}
		if info.IsDir() {
			return nil
		}

		if header.Mode().IsRegular() {
			file, err := os.Open(filePath)
			if err != nil {
				return fmt.Errorf("%s: opening: %v", filePath, err)
			}
			defer file.Close()
			r := bufio.NewReader(file)
			buf := make([]byte, 256)
			_, err = io.CopyBuffer(f, r, buf)
			if err != nil && err != io.EOF {
				return fmt.Errorf("%s: copying contents: %v", filePath, err)
			}
		}
		return nil
	})
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
