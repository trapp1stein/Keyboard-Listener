//go:build !windows
// +build !windows

package keyboard

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"io/ioutil"
	"os"
	"strings"
)

const (
	ti_header_length = 12
)

var (
	eterm_keys = []string{
		"\x1b[11~", "\x1b[12~", "\x1b[13~", "\x1b[14~", "\x1b[15~", "\x1b[17~", "\x1b[18~", "\x1b[19~", "\x1b[20~", "\x1b[21~", "\x1b[23~", "\x1b[24~", "\x1b[2~", "\x1b[3~", "\x1b[7~", "\x1b[8~", "\x1b[5~", "\x1b[6~", "\x1b[A", "\x1b[B", "\x1b[D", "\x1b[C",
	}
	screen_keys = []string{
		"\x1bOP", "\x1bOQ", "\x1bOR", "\x1bOS", "\x1b[15~", "\x1b[17~", "\x1b[18~", "\x1b[19~", "\x1b[20~", "\x1b[21~", "\x1b[23~", "\x1b[24~", "\x1b[2~", "\x1b[3~", "\x1b[1~", "\x1b[4~", "\x1b[5~", "\x1b[6~", "\x1bOA", "\x1bOB", "\x1bOD", "\x1bOC",
	}
	xterm_keys = []string{
		"\x1bOP", "\x1bOQ", "\x1bOR", "\x1bOS", "\x1b[15~", "\x1b[17~", "\x1b[18~", "\x1b[19~", "\x1b[20~", "\x1b[21~", "\x1b[23~", "\x1b[24~", "\x1b[2~", "\x1b[3~", "\x1b[H", "\x1b[F", "\x1b[5~", "\x1b[6~", "\x1b[A", "\x1b[B", "\x1b[D", "\x1b[C",
	}
	rxvt_keys = []string{
		"\x1b[11~", "\x1b[12~", "\x1b[13~", "\x1b[14~", "\x1b[15~", "\x1b[17~", "\x1b[18~", "\x1b[19~", "\x1b[20~", "\x1b[21~", "\x1b[23~", "\x1b[24~", "\x1b[2~", "\x1b[3~", "\x1b[7~", "\x1b[8~", "\x1b[5~", "\x1b[6~", "\x1b[A", "\x1b[B", "\x1b[D", "\x1b[C",
	}
	linux_keys = []string{
		"\x1b[[A", "\x1b[[B", "\x1b[[C", "\x1b[[D", "\x1b[[E", "\x1b[17~", "\x1b[18~", "\x1b[19~", "\x1b[20~", "\x1b[21~", "\x1b[23~", "\x1b[24~", "\x1b[2~", "\x1b[3~", "\x1b[1~", "\x1b[4~", "\x1b[5~", "\x1b[6~", "\x1b[A", "\x1b[B", "\x1b[D", "\x1b[C",
	}

	terms = []struct {
		name string
		keys []string
	}{
		{"Eterm", eterm_keys},
		{"screen", screen_keys},
		{"xterm", xterm_keys},
		{"xterm-256color", xterm_keys},
		{"rxvt-unicode", rxvt_keys},
		{"rxvt-256color", rxvt_keys},
		{"linux", linux_keys},
	}
)

func load_terminfo() ([]byte, error) {
	var data []byte
	var err error

	term := os.Getenv("TERM")
	if term == "" {
		return nil, errors.New("terminfo: TERM not set")
	}
	// Kontrol edelim terminal içindemi girildi?
	for _, t := range terms {
		if t.name == term {
			return nil, errors.New("use built in!")
		}
	}

	// terminfo(5)

	terminfo := os.Getenv("TERMINFO")
	if terminfo != "" {
		return ti_try_path(terminfo)
	}

	// ~/.terminfo devam..
	home := os.Getenv("HOME")
	if home != "" {
		data, err = ti_try_path(home + "/.terminfo")
		if err == nil {
			return data, nil
		}
	}

	// sonraki, TERMINFO_DIRS
	dirs := os.Getenv("TERMINFO_DIRS")
	if dirs != "" {
		for _, dir := range strings.Split(dirs, ":") {
			if dir == "" {
				// "" -> "/usr/share/terminfo"
				dir = "/usr/share/terminfo"
			}
			data, err = ti_try_path(dir)
			if err == nil {
				return data, nil
			}
		}
	}

	// devamı, /lib/terminfo
	data, err = ti_try_path("/lib/terminfo")
	if err == nil {
		return data, nil
	}

	// fallback gelişi /usr/share/terminfo
	return ti_try_path("/usr/share/terminfo")
}

func ti_try_path(path string) (data []byte, err error) {
	term := os.Getenv("TERM")

	// *nix path ilk deneme
	terminfo := path + "/" + term[0:1] + "/" + term
	data, err = ioutil.ReadFile(terminfo)
	if err == nil {
		return
	}

	// fallback darwin dirs yapısı
	terminfo = path + "/" + hex.EncodeToString([]byte(term[:1])) + "/" + term
	data, err = ioutil.ReadFile(terminfo)
	return
}

func setup_term_builtin() error {
	name := os.Getenv("TERM")
	if name == "" {
		return errors.New("terminfo: TERM environment variable not set")
	}

	for _, t := range terms {
		if t.name == name {
			keys = t.keys
			return nil
		}
	}

	compat_table := []struct {
		partial string
		keys    []string
	}{
		{"xterm", xterm_keys},
		{"xterm-256color", xterm_keys},
		{"rxvt", rxvt_keys},
		{"rxvt-unicode", rxvt_keys},
		{"rxvt-256color", rxvt_keys},
		{"linux", linux_keys},
		{"Eterm", eterm_keys},
		{"screen", screen_keys},
		{"cygwin", xterm_keys},		// 'cygwin' + xterm uyumlu?
		{"st", xterm_keys},
	}

	// birleşebilir varyasyonları deneyelim.
	for _, it := range compat_table {
		if strings.Contains(name, it.partial) {
			keys = it.keys
			return nil
		}
	}

	return errors.New("termbox: unsupported terminal")
}

func setup_term() (err error) {
	var data []byte
	var header [6]int16
	var str_offset, table_offset int16

	data, err = load_terminfo()
	if err != nil {
		return setup_term_builtin()
	}

	rd := bytes.NewReader(data)
	// 0: magic numaramız, 1: isim section boyutu, 2: boolean section boyutu, 3:
	// number section boyutu (integer'larda), 4: strings section boyutu (integer'larda), 5: string tablo boyutu

	err = binary.Read(rd, binary.LittleEndian, header[:])
	if err != nil {
		return
	}

	if header[0] != 542 && header[0] != 282 {
		return setup_term_builtin()
	}

	number_sec_len := int16(2)
	if header[0] == 542 { // burası octal 0542 olma ihtimali olabilir, terminfo dosyaları 542 varsayımlar.
		number_sec_len = 4
	}

	if (header[1]+header[2])%2 != 0 {
		header[2] += 1 // eski boundary yapılarda 1 bayrak eklenmesi yapılır.

	}
	str_offset = ti_header_length + header[1] + header[2] + number_sec_len*header[3]
	table_offset = str_offset + 2*header[4]

	keys = make([]string, 0xFFFF-key_min)
	for i := range keys {
		keys[i], err = ti_read_string(rd, str_offset+2*ti_keys[i], table_offset)
		if err != nil {
			return
		}
	}
	return nil
}

func ti_read_string(rd *bytes.Reader, str_off, table int16) (string, error) {
	var off int16

	_, err := rd.Seek(int64(str_off), 0)
	if err != nil {
		return "", err
	}
	err = binary.Read(rd, binary.LittleEndian, &off)
	if err != nil {
		return "", err
	}
	_, err = rd.Seek(int64(table+off), 0)
	if err != nil {
		return "", err
	}
	var bs []byte
	for {
		b, err := rd.ReadByte()
		if err != nil {
			return "", err
		}
		if b == byte(0x00) {
			break
		}
		bs = append(bs, b)
	}
	return string(bs), nil
}

// "Maps" özel olarak hazırlanır.
var ti_keys = []int16{
	66, 68, 69, 70,
	71, 72, 73, 74, 75, 67, 216, 217, 77, 59, 76, 164, 82, 81, 87, 61, 79, 83,
}
