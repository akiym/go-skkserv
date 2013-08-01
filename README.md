# go-skkserv

Lightweight skkserv implementation for golang

## Example

### Google IME SKK

    package main

    import (
    	"code.google.com/p/mahonia"
    	"encoding/json"
    	"github.com/akiym/go-skkserv"
    	"net/http"
    	"net/url"
    )

    type GoogleIMESKK struct{}

    func (s *GoogleIMESKK) Request(text string) ([]string, error) {
    	words, err := Transliterate(text)
    	if err != nil {
    		return nil, err
    	}
    	return words, nil
    }

    var enc = mahonia.NewEncoder("euc-jp")

    func Transliterate(text string) (words []string, err error) {
    	text = enc.ConvertString(text)
    	v := url.Values{"langpair": {"ja-Hira|ja"}, "text": {text + ","}}
    	resp, err := http.Get("http://www.google.com/transliterate?" + v.Encode())
    	if err != nil {
    		return nil, err
    	}
    	defer resp.Body.Close()
    	dec := json.NewDecoder(resp.Body)
    	var w [][]interface{}
    	if err := dec.Decode(&w); err != nil {
    		return nil, err
    	}
    	for _, v := range w[0][1].([]interface{}) {
    		word := v.(string)
    		result, ok := enc.ConvertStringOK(word)
    		if ok {
    			words = append(words, result)
    		}
    	}
    	return words, nil
    }

    func main() {
    	var server = skkserv.NewServer(":55100", &GoogleIMESKK{})
    	server.Run()
    }
