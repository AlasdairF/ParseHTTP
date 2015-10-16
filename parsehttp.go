package parsehttp

import (
 "bytes"
 "mime"
 "net/http"
 "io/ioutil"
 "errors"
)

func EasyGET(r *http.Request, fields []string) ([][]byte, error) {
	res := make([][]byte, len(fields))
	var i int
	var c, kstring string
	fn := func(k, v []byte) {
		kstring = string(k)
		for i, c = range fields {
			if kstring == c {
				res[i] = v
				break
			}
		}
	}
	err := GET(r, fn)
	return res, err
}

func EasyPOST(r *http.Request, fields []string) ([][]byte, error) {
	res := make([][]byte, len(fields))
	var i int
	var c, kstring string
	fn := func(k, v []byte) {
		kstring = string(k)
		for i, c = range fields {
			if kstring == c {
				res[i] = v
				break
			}
		}
	}
	err := POST(r, fn)
	return res, err
}

func POST(r *http.Request, fn func([]byte, []byte)) error {
	
	// Check that a request body even exists to parse
	if r.Body == nil {
		return nil
	}
	
	// Ensure content-type is correct
	var err error
	ct := r.Header.Get(`Content-Type`)
	if ct == `` {
		ct = `application/octet-stream`
	}
	ct, _, err = mime.ParseMediaType(ct)
	if err != nil {
		return err
	}
	if ct != `application/x-www-form-urlencoded` {
		return errors.New(`Form type is not application/x-www-form-urlencoded`)
	}
	
	// Read in the request body
	query, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	
	// Parse the request body
	var i int
	var key []byte
	for len(query) > 0 {
		key = query
		if i = bytes.IndexAny(key, "&;"); i >= 0 {
			key, query = key[:i], key[i+1:]
		} else {
			query = []byte{}
		}
		if len(key) == 0 {
			continue
		}
		var value []byte
		if i = bytes.IndexByte(key, '='); i >= 0 {
			key, value = key[:i], key[i+1:]
		}
		key, err1 := unescape(key)
		if err1 != nil {
			if err == nil {
				err = err1
			}
			continue
		}
		value, err1 = unescape(value)
		if err1 != nil {
			if err == nil {
				err = err1
			}
			continue
		}
		fn(key, value)
	}
	return err
}

func GET(r *http.Request, fn func([]byte, []byte)) error {
	
	var err error
	query := []byte(r.URL.RawQuery)
	
	// Parse the request body
	var i int
	var key []byte
	for len(query) > 0 {
		key = query
		if i = bytes.IndexAny(key, "&;"); i >= 0 {
			key, query = key[:i], key[i+1:]
		} else {
			query = []byte{}
		}
		if len(key) == 0 {
			continue
		}
		var value []byte
		if i = bytes.IndexByte(key, '='); i >= 0 {
			key, value = key[:i], key[i+1:]
		}
		key, err1 := unescape(key)
		if err1 != nil {
			if err == nil {
				err = err1
			}
			continue
		}
		value, err1 = unescape(value)
		if err1 != nil {
			if err == nil {
				err = err1
			}
			continue
		}
		fn(key, value)
	}
	return err
}

func unescape(s []byte) ([]byte, error) {
	// Count %, check that they're well-formed.
	var n int
	hasPlus := false
	for i := 0; i < len(s); {
		switch s[i] {
		case '%':
			n++
			if i+2 >= len(s) || !ishex(s[i+1]) || !ishex(s[i+2]) {
				s = s[i:]
				if len(s) > 3 {
					s = s[0:3]
				}
				return nil, errors.New(`Escape error`)
			}
			i += 3
		case '+':
			hasPlus = true
			i++
		default:
			i++
		}
	}

	if n == 0 && !hasPlus {
		return s, nil
	}

	t := make([]byte, len(s)-2*n)
	j := 0
	for i := 0; i < len(s); {
		switch s[i] {
		case '%':
			t[j] = unhex(s[i+1])<<4 | unhex(s[i+2])
			j++
			i += 3
		case '+':
			t[j] = ' '
			j++
			i++
		default:
			t[j] = s[i]
			j++
			i++
		}
	}
	return t, nil
}

func ishex(c byte) bool {
	switch {
	case '0' <= c && c <= '9':
		return true
	case 'a' <= c && c <= 'f':
		return true
	case 'A' <= c && c <= 'F':
		return true
	}
	return false
}

func unhex(c byte) byte {
	switch {
	case '0' <= c && c <= '9':
		return c - '0'
	case 'a' <= c && c <= 'f':
		return c - 'a' + 10
	case 'A' <= c && c <= 'F':
		return c - 'A' + 10
	}
	return 0
}
