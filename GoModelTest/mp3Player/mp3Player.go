/*
#Time      :  2019/7/23 上午9:38 
#Author    :  chuangangshen@deepglint.com
#File      :  mp3Player.go
#Software  :  GoLand
*/
package main

import (
	"encoding/binary"
	"io"
	"github.com/gordonklaus/portaudio"
	"os"
	"os/signal"
	"fmt"
	"github.com/bobertlo/go-mpg123/mpg123"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "usage: mp3dump <infile.mp3> <outfile.raw>")
		return
	}
	mp3File := os.Args[1]
	rawFile := os.Args[2]
	Mp3ToRaw(mp3File, rawFile)
	PortaudioPlayMp3(rawFile)
}

func Mp3ToRaw(mp3File, rawFile string) {
	// create mpg123 decoder instance
	decoder, err := mpg123.NewDecoder("")
	if err != nil {
		panic("could not initialize mpg123")
	}

	// open a file with decoder
	err = decoder.Open(mp3File)
	if err != nil {
		panic("error opening mp3 file")
	}
	defer decoder.Close()

	// get audio format information
	rate, chans, _ := decoder.GetFormat()
	fmt.Fprintln(os.Stderr, "Encoding: Signed 16bit")
	fmt.Fprintln(os.Stderr, "Sample Rate:", rate)
	fmt.Fprintln(os.Stderr, "Channels:", chans)

	// make sure output format does not change
	decoder.FormatNone()
	decoder.Format(rate, chans, mpg123.ENC_SIGNED_16)

	// open output file
	o, err := os.Create(rawFile)
	if err != nil {
		panic("error opening output file")
	}
	defer o.Close()

	// decode mp3 file and dump output
	buf := make([]byte, 2048*16)
	for {
		len, err := decoder.Read(buf)
		o.Write(buf[0:len])
		if err != nil {
			break
		}
	}
	o.Close()
	decoder.Delete()
}

func PortaudioPlayMp3(rawFile string)  {

	fmt.Println("Playing.  Press Ctrl-C to stop.")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)

	fileName := rawFile
	f, err := os.Open(fileName)
	chk(err)
	defer f.Close()

	id, data, err := readChunk(f)
	chk(err)
	if id.String() != "FORM" {
		fmt.Println("bad file format")
		return
	}
	_, err = data.Read(id[:])
	chk(err)
	if id.String() != "AIFF" {
		fmt.Println("bad file format")
		return
	}
	var c commonChunk
	var audio io.Reader
	for {
		id, chunk, err := readChunk(data)
		if err == io.EOF {
			break
		}
		chk(err)
		switch id.String() {
		case "COMM":
			chk(binary.Read(chunk, binary.BigEndian, &c))
		case "SSND":
			chunk.Seek(8, 1) //ignore offset and block
			audio = chunk
		default:
			fmt.Printf("ignoring unknown chunk '%s'\n", id)
		}
	}

	//assume 44100 sample rate, mono, 32 bit

	portaudio.Initialize()
	defer portaudio.Terminate()
	out := make([]int32, 8192)
	stream, err := portaudio.OpenDefaultStream(0, 1, 44100, len(out), &out)
	chk(err)
	defer stream.Close()

	chk(stream.Start())
	defer stream.Stop()
	for remaining := int(c.NumSamples); remaining > 0; remaining -= len(out) {
		if len(out) > remaining {
			out = out[:remaining]
		}
		err := binary.Read(audio, binary.BigEndian, out)
		if err == io.EOF {
			break
		}
		chk(err)
		chk(stream.Write())
		select {
		case <-sig:
			return
		default:
		}
	}
}

func readChunk(r readerAtSeeker) (id ID, data *io.SectionReader, err error) {
	_, err = r.Read(id[:])
	if err != nil {
		return
	}
	var n int32
	err = binary.Read(r, binary.BigEndian, &n)
	if err != nil {
		return
	}
	off, _ := r.Seek(0, 1)
	data = io.NewSectionReader(r, off, int64(n))
	_, err = r.Seek(int64(n), 1)
	return
}

type readerAtSeeker interface {
	io.Reader
	io.ReaderAt
	io.Seeker
}

type ID [4]byte

func (id ID) String() string {
	return string(id[:])
}

type commonChunk struct {
	NumChans      int16
	NumSamples    int32
	BitsPerSample int16
	SampleRate    [10]byte
}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}
