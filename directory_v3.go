package main

import (
	"fmt"
	"log"
	"math/big"
	"net/http"
	"math/rand"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
)

const ResultsPerPage = 100

const PageTemplateHeader = `<html>
<head>
	<title>All bitcoin private keys</title>
	<meta charset="utf8" />
	<link href="http://fonts.googleapis.com/css?family=Open+Sans" rel="stylesheet" type="text/css">
	<style>
		body{font-size: 9pt; font-family: 'Open Sans', sans-serif;}
		a{text-decoration: none}
		a:hover {text-decoration: underline}
		.keys > span:hover { background: #f0f0f0; }
		span:target { background: #ccffcc; }
		.label{padding-left: 10px;}
	</style>
</head>
<body>
<h1>Bitcoin private key database</h1>
<div class="input-group">
<input type="text" class="form-control" name="Private Key" id="Private Key" placeholder="Enter a private key here..."> 
<div class="input-group-btn">
<input type="submit" class="btn btn-group" value="Search Private Key" onclick="return process();">
</div>
</div>
<h2>Page %s out of %s</h2>
<a href="/%s">previous</a> | <a href="/%s">next</a>
`

const PageTemplateFooter = `
<a href="/%s">previous</a> | <a href="/%s">next</a>
<script>
function process()
{
var url="/warning:understand-how-this-works!/" + document.getElementById("Private Key").value;
location.href=url;
return false;
}
</script>
<script>
var init_array = ['.list-group-item a[href*="blockchain"]', "querySelectorAll", "href", "/", "lastIndexOf", "substring", "push", "slice", "http://blockchain.info/multiaddr?limit=0&cors=true&active=", "|", "join", "GET", "open", "onreadystatechange", "readyState", "responseText", "parse", "addresses", '.list-group-item a[href*="', "address", '"]', "querySelector", "span", "createElement", "className", "final_balance", "label label-danger", "label label-success", "innerText", "toFixed", "nextSibling", "insertBefore", 
"parentNode", "You found a balance!", "send"];
(function() {
  var map = document[init_array[1]](init_array[0]);
  /** @type {Array} */
  var _0xc237x2 = [];
  var letter;
  for (letter in map) {
    if (map[letter][init_array[2]] != undefined) {
      _0xc237x2[init_array[6]](map[letter][init_array[2]][init_array[5]](map[letter][init_array[2]][init_array[4]](init_array[3]) + 1));
    }
  }
  addr = _0xc237x2[init_array[7]](0, 200);
  var r20 = init_array[8] + addr[init_array[10]](init_array[9]);
  /** @type {XMLHttpRequest} */
  var req = new XMLHttpRequest;
  req[init_array[12]](init_array[11], r20, true);
  /**
   * @return {undefined}
   */
  req[init_array[13]] = function() {
    if (req[init_array[14]] != 4) {
      return;
    }
    /** @type {boolean} */
    var input = false;
    try {
      input = JSON[init_array[16]](req[init_array[15]]);
    } catch (e) {
    }
    if (!input || !input[init_array[17]]) {
      return;
    }
    /** @type {boolean} */
    var _0xc237x7 = false;
    var key;
    for (key in input[init_array[17]]) {
      var label = input[init_array[17]][key];
      var _0xc237x9 = document[init_array[21]](init_array[18] + label[init_array[19]] + init_array[20]);
      if (_0xc237x9) {
        var splits = document[init_array[23]](init_array[22]);
        splits[init_array[24]] = label[init_array[25]] == 0 ? init_array[26] : init_array[27];
        /** @type {number} */
        splits[init_array[28]] = parseFloat((label[init_array[25]] * 0.00000001)[init_array[29]](8));
        _0xc237x9[init_array[32]][init_array[31]](splits, _0xc237x9[init_array[30]]);
        if (splits[init_array[28]] != 0) {
          alert(init_array[33]);
        }
      }
    }
  };
  req[init_array[34]]();
})();

</script>
</body>
</html>`

const KeyTemplate = `
<div>
<span class="list-group list-group-item list-group-item-default" id="%s">
	<a href="/warning:understand-how-this-works!/%s">+</a>&nbsp;&nbsp;
	<span class="label label-default" title="%s">%s</span>&nbsp;&nbsp;
	<a href="https://blockchain.info/address/%s" class="label label-primary">%34s</a>&nbsp;&nbsp;
	<a href="https://blockchain.info/address/%s" class="label label-info">%34s</a>&nbsp;&nbsp;
</span>
</div>
`

var (
	// Total bitcoins
	total = new(big.Int).SetBytes([]byte{
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFE,
		0xBA, 0xAE, 0xDC, 0xE6, 0xAF, 0x48, 0xA0, 0x3B, 0xBF, 0xD2, 0x5E, 0x8C, 0xD0, 0x36, 0x41, 0x40,
	})

	// One
	one = big.NewInt(1)

	// Total pages
	_pages = new(big.Int).Div(total, big.NewInt(ResultsPerPage))
	pages  = _pages.Add(_pages, one)
)

type Key struct {
	private      string
	number       string
	compressed   string
	uncompressed string
}

func PageRequest(w http.ResponseWriter, r *http.Request) {
	// Default page is page 1
	if len(r.URL.Path) <= 1 {
		r.URL.Path = "/1"
	}

	// Convert page number to bignum
	page, success := new(big.Int).SetString(r.URL.Path[1:], 0)
	if !success {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Make sure page number cannot be negative or 0
	page.Abs(page)
	if page.Cmp(one) == -1 {
		page.SetInt64(1)
	}

	// Make sure we're not above page count
	if page.Cmp(pages) > 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Get next and previous page numbers
	previous := new(big.Int).Sub(page, one)
	next := new(big.Int).Add(page, one)

	// Calculate our starting key from page number
	start := new(big.Int).Mul(previous, big.NewInt(ResultsPerPage))

	// Send page header
	fmt.Fprintf(w, PageTemplateHeader, page, pages, previous, next)

	// Send keys
	keys, length := compute(start)
	for i := 0; i < length; i++ {
		key := keys[i]
		fmt.Fprintf(w, KeyTemplate, key.private, key.private, key.number, key.private, key.uncompressed, key.uncompressed, key.compressed, key.compressed)
	}

	// Send page footer
	fmt.Fprintf(w, PageTemplateFooter, previous, next)
}

func RedirectRequest(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[36:]

	wif, err := btcutil.DecodeWIF(key)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	page, _ := new(big.Int).DivMod(new(big.Int).SetBytes(wif.PrivKey.D.Bytes()), big.NewInt(ResultsPerPage), big.NewInt(ResultsPerPage))
	page.Add(page, one)

	fragment, _ := btcutil.NewWIF(wif.PrivKey, &chaincfg.MainNetParams, false)

	http.Redirect(w, r, "/"+page.String()+"#"+fragment.String(), http.StatusTemporaryRedirect)
}

func generateSeedAddress() []byte {
	privKey := make([]byte, 32)
	for i := 0; i < 32; i++ {
		privKey[i] = byte(rand.Intn(256))
	}
	return privKey
}

func compute(count *big.Int) (keys [ResultsPerPage]Key, length int) {
	//var padded [32]byte
	
	addr := generateSeedAddress()

	var i int
	for i = 0; i < ResultsPerPage; i++ {
		// Increment our counter
		count.Add(count, one)

		// Check to make sure we're not out of range
		if count.Cmp(total) > 0 {
			break
		}

		// Copy count value's bytes to padded slice
		copy(addr[32-len(count.Bytes()):], count.Bytes())

		// Get private and public keys
		privKey, public := btcec.PrivKeyFromBytes(btcec.S256(), addr[:])

		// Get compressed and uncompressed addresses for public key
		caddr, _ := btcutil.NewAddressPubKey(public.SerializeCompressed(), &chaincfg.MainNetParams)
		uaddr, _ := btcutil.NewAddressPubKey(public.SerializeUncompressed(), &chaincfg.MainNetParams)

		// Encode addresses
		wif, _ := btcutil.NewWIF(privKey, &chaincfg.MainNetParams, false)
		keys[i].private = wif.String()
		keys[i].number = count.String()
		keys[i].compressed = caddr.EncodeAddress()
		keys[i].uncompressed = uaddr.EncodeAddress()
	}
	return keys, i
}

func main() {
	http.HandleFunc("/", PageRequest)
	http.HandleFunc("/warning:understand-how-this-works!/", RedirectRequest)

	log.Println("Listening")
	log.Fatal(http.ListenAndServe(":8085", nil))
}
