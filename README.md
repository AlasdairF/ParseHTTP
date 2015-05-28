##ParsePost

This package replaces http.ParseForm. It's a slightly simpler and more optimized implementation when all you want is to parse a simple POST form (which is often the case), excluding GET, PUT & multipart forms.

Rather than returning a map[string][]string, parsepost.Form takes the *http.Request and a function var func([]byte, []byte) as an input. You can build a map with this or whatever you want.

###Usage

     fn := func(key, value []byte) {
        // Do whatever you want with the key and value here
        // Note that they are slices and may need to be copied (unleaked) if they are to be stored
     }
     err := parsepost.Form(r, fn)
