/* binsanity.go - auto-generated; edit at your own peril!

More info: https://github.com/biztos/binsanity

Generated: 

*/

package binsanity

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"errors"
	"io"
	"sort"
)

// Asset returns the byte content of the asset for the given name, or an error
// if no such asset is available.
func Asset(name string) ([]byte, error) {

	_, found := binsanity_cache[name]
	if !found {
		i := sort.SearchStrings(binsanity_names, name)
		if i == len(binsanity_names) || binsanity_names[i] != name {
			return nil, errors.New("Asset not found.")
		}

		// We ignore errors because we controlled the data from the begining.
		// It's not perfect but it seems better than having additional funcs
		// hanging around that might confuse the user: tried that already, not
		// nicer.
		decoded, _ := base64.StdEncoding.DecodeString(binsanity_data[i])
		buf := bytes.NewReader(decoded)
		gzr, _ := gzip.NewReader(buf)
		defer gzr.Close()
		data, _ := io.ReadAll(gzr)

		// Not cached, so decode and cache it.
		binsanity_cache[name] = data

	}
	return binsanity_cache[name], nil

}

// MustAsset returns the byte content of the asset for the given name, or
// panics if no such asset is available.
func MustAsset(name string) []byte {
	b, err := Asset(name)
	if err != nil {
		panic(err.Error())
	}
	return b
}

// MustAssetString returns the string content of the asset for the given name,
// or panics if no such asset is available.  This is a convenience function
// for string(MustAsset(name)).
func MustAssetString(name string) string {
	return string(MustAsset(name))
}

// AssetNames returns the sorted names of the assets.
func AssetNames() []string {
	return binsanity_names
}

// this must remain sorted or everything breaks!
var binsanity_names = []string{
	"code.tmpl",
	"tests.tmpl",
}

// only decode once per asset.
var binsanity_cache = map[string][]byte{}

// assets are gzipped and base64 encoded
var binsanity_data = []string{
	"H4sIAAAAAAAA/6RVTWvkRhA9q39FeSBEWmalS8jBYQ7LrgM5rAlxIAdjTEtdkorVdCvdpTFjrf57qG75YwYHDHtZ1qWq96peveqpPsA8l5+dwd9pwGWBj6Andh87tOg1o/kN0BCDZji6yYN7sDCip+FCqa/OI5Bt3SX0zGO4rKqOuJ/qsnH7qqZHdqGqyQZtiY9KfaiUGnXzTXcopH+m/y6LUrQfnWfIVbapj4xho7JN4/ajxxCq7pFGCaBtnCHbVbUO+OsvMeS98zGbnPwbnOeNKpSqKvgUAjJ45MnbANwjCDQ0zjJaBtfGmI5ZrfPxr44OaMHqPW7BedAWIoPAUQvWQZiafq2hAPqgadD1gKVqJ9skylzKIbAn2xWQ394J7TYBFTArld1voXWTNXC5g2d57hvd9HgrxXcqoxYuUs6ssowkU2Yrb1D7pr+J4CF/KZaysI2dF1LQAsFuBwPa86QCvn+Hs9gt3cHFLlZHviypBpaGtfFQXuNDvkmaWsdpgHIjZItSWVZV8A8CdVY8kUqgxkZPAeEhqe7dMKCJOhvNGlrv9mkv2JEl25UJ5w/+OUSOEX2LDUM9MRBDQNwLKDPKtrSFXh/IdqCNISZn9QCyhpBgem27+NVHHbnXDHvqepZmWulLuKeA/hLYE64pevCozXErHSQgSw166c1g4wyaLdzHzUUXljdsrlZjll9iQtrOK91l2lu6E63qqY21YnKR9C/UBn2+IktG9+hXAvH9q5x6aovYRIseukdffh5cwDzGNOu1iFwp+Z+GIe8efbGu5toxRH+ZLQQHiQ60NSkKxDLfm16EXdyWUtminnzxZuJW7KLUEo/v6xT4xw9QkEZtqQnvOr9n0tMTTBcoxq6jm0Wml7QiHpuE5QJoiAcQSXP0vrwSK+dFcTL++ZBp4yejJvJ3Dytwzr9vWIC/ewoxJvgHtIS2weh9OQPBEorUQn6qSlGci7Xa9USytfv5eeT/wVqFiKFreUlONXCe0cQBw4kC4fV7GetyWdMz6zxTC2X8Gq72Ix+X5amRp6x5mWccAr58OXvT5hmtWZa1QxbF9lMQP+412afenAc8oD9yL8y1R/0tXKiD9udwsHuhVvPste0Qytj7smTzPHqy3MLmp383UC7LVq38id7Z4fh0dU52NaJPUpRnZOked7DX422iu0v+nRNS0g+0x/g+jGjiGafHCOJPJJozzPjavtn/F816/TFZlmwzz+WybF71/l8AAAD//2m7bSweCAAA",
	"H4sIAAAAAAAA/9RXXW/bOhJ9Fn/FREAAuVBltED74MJYtLvOog9Ji7WLIsgGAS1REhGJ1JIjN15B//1iSMXxh5K2L/fivtiiSB7OnDOcGU1fQdclK2HxQlai7+E18Bb160IoYTiK7AOITCJwhK1uDegfChphZHXG2EoDCouApYC0FOm9bWsLuTbAqwpSrVAojMEKv0SojTRa1UIhbLiRfF0J9unz1fLj1efV9d1qsVzd/fPL1WpxtQLUoJUAnc/gOr5eLONVvPrPt0X8BiKCWpkWyy0sS22wkhYnCWOX2giQKtczKBEbO5tOC4llu05SXU/X8v+o7XQtleVK4paxV1PGGp7e80IQBV/9Y9/fkU+MybrRBiFiQZiabYN6akv+9t37kAVhXiP9aUu/Fo1UhXuknVIVIWNB2HXJpc5a4jRkE8ZSrSzCp8fzP1or8FJaK1UBc+i6xkiFOYTn/wshGSbcoitei74f3f/VCEtsnuxfPEhnyS8CLNv6JxjLtu57xjbcHCEQuIU53Nx6HjrWdTKHxE3aRd3gtu+D6RQEPbKuE5Ula7rOcFUISBxA3wdHp/d9TItV1vfD3+jxS4q4w9MH3H9x5DT7InTPWN6qFCj+n9yJEF4NUiarCXSMBcq5OZsfREqyt2XCAplDJVTklk7gbO5GI2wRYhBgcsGRV3kUfjdaFaDaei0M6BwcwOy/CkA8NCJFkc3gPKMxT7HlFY3CmAVB8NwB8Z4hExb0jJECSZLUmm6fhR+lUHSfwQheVVuopbWOBJlvkySBdYtw9QUy0fj7jKVwELusALmshD1jAc3K7CEGRex47j1Z5KTMQRETI0beyOzh1i3ao+JS2ppjWlK6Oc+OObAHHFjPQeBOf+6AGNSEBUHvSDhVW+OFblU2IvhdDMKYccWj0VvsI4A2zeegZLWvchReaZrSxtFZD/ee0/YkdBINm5MFrYpc/IQOHpRGyMnMJDzE9IHzMuyJz885vP51h4es8eTw2YjDi51ZzhzAkiPYUrdV5jxai0d7HwmwbU3H5zUmS39ho/D8IYzBJ95k2dZv372P1hN/MC0/Ca29hDZClQdyO58MyzjyEbIu258HiXhoyOJTnVjQcCXT+61zqFVpNIHukNgd/tJlrmdiCnoW0AuDXwnQfpdYRhjDAB/T9Ygh3IFBtLNiEk7YuEfPhsCJ+rstz0bAXy/auI+e1b+Vdh7slxT0S59z7bRQvWzwb6t5c7veoojs5E9SdVfqT5wNqCX4wRW6rs3CWuvKVapVKS1ICxwqiVgJWEsEvRHmXlaVq2+N0E0loOQb+llLtGBkUeI/WEAo0pZERc2bG99b3NJb8iS8DmcAAGhaEbvxYklvduPV0Tx1ruHsafzmYL5nQV7xgg4bWslkpb81jTCRtsm/BQq1icLxPjkkAfbcn8Ng+g1B3jpxzvbmvRDLe9lQwxIYga1R3oRdHee1eCrlY80egbxYLQhi4qv/UXXYK/Xndgbnm9Af6NCGQv14M08bPVfRWfBbOWcvPAn3yIjRWKQodF8evuE4bkBGWpDBBbrOtq13DUfP2HQKFxTeUNPXSWtF3lawEcZKrajXQwpTK8RL3ywu4At/G07zyd6FiCEfEpYzZginGGpbDM8+P/gEJDJHIq+sYEGh0WXDkAWPKY8FQSZyYWDvBbHpVDcipcsUTT7AocBP4HMX4O4doR/rZcMYnORO897FY06/9OgDdwd12C9ftCpF4i+TvpVw67wc5CthAH1mgMyBjt5pfwDjSHT7uT0R+eei17YYJC80PjXZ0WcFhf8udtzbNk2FdanIykooTFxO/yMAAP//CGDlWXcPAAA=",
}
