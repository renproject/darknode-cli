package darknode

import (
	"errors"
	"math/big"

	"github.com/btcsuite/btcd/btcec"
	"github.com/ethereum/go-ethereum/common"
	"github.com/renproject/rzl/asm/value"
	"github.com/republicprotocol/darknode-cli/darknode/addr"
)

// The Network type defines the different RenVM networks that exist.
type Network string

const (
	// Devnet is used as an internal staging environment to run large-scale
	// chaos tests.
	Devnet = Network("devnet")

	// Testnet is used by third-party developers to build their application
	// without incurring real-world costs.
	Testnet = Network("testnet")

	// Chaosnet is the pre-production version of mainnet for production testing
	// with real-world incentives and punishments.
	Chaosnet = Network("chaosnet")

	// Mainnet is production network.
	Mainnet = Network("mainnet")
)

// NewNetwork parses the string to a specific Network.
func NewNetwork(network string) (Network, error) {
	switch network {
	case "devnet":
		return Devnet, nil
	case "testnet":
		return Testnet, nil
	case "chaosnet":
		return Chaosnet, nil
	case "mainnet":
		return Mainnet, nil
	default:
		return "", errors.New("unknown network")
	}
}

// PublicKey returns the public key of .
func (network Network) PublicKey() value.ECDSAPubKey {
	var x *big.Int
	var xOk bool
	switch network {
	case Mainnet:
		panic("unimplemented")
	case Chaosnet:
		x, xOk = big.NewInt(0).SetString("54769503130895894163949174470748707835675520766218565814337221309492303621497", 10)
	case Testnet:
		x, xOk = big.NewInt(0).SetString("6258831358146983420781042002047732738577946776960027585197438124940321371484", 10)
	case Devnet:
		x, xOk = big.NewInt(0).SetString("15988544014143623672381260113528425219902628559557338422452853193412094205021", 10)
	}
	if !xOk {
		panic("invalid x for shared public key")
	}

	var y *big.Int
	var yOk bool
	switch network {
	case Mainnet:
		panic("unimplemented")
	case Chaosnet:
		y, yOk = big.NewInt(0).SetString("87745967375764291795837331450616094559320177780884666147029497390497322495569", 10)
	case Testnet:
		y, yOk = big.NewInt(0).SetString("22471449852503869623778529670369476102885501157580830978857986155713794677963", 10)
	case Devnet:
		y, yOk = big.NewInt(0).SetString("18856215896348556820657579775929067367606870869895908607019381042047463166252", 10)
	}
	if !yOk {
		panic("invalid y for shared public key")
	}

	return value.ECDSAPubKey{
		Curve: btcec.S256(),
		X:     x,
		Y:     y,
	}
}

func (network Network) BootstrapNodes() []addr.MultiAddress {
	bootstraps := make([]addr.MultiAddress, 0, 16)
	switch network {
	case Mainnet:
		panic("unimplemented")
	case Chaosnet:
		b1, _ := addr.NewSignedMultiAddressFromString("/ip4/3.115.117.251/tcp/18514/ren/8MGrkr3CCG5gxnipWD5RUc8BMQnU1s", "3w8PoELIeSh0sqcb6qONy1FNgSIgP9hELVh44D/IE0saY518C9vWvBYSQn4xUmYRb7Y+nYNPY54NoH1y0zMnXAE=")
		b2, _ := addr.NewSignedMultiAddressFromString("/ip4/18.182.28.215/tcp/18514/ren/8MGjmhtNxsqT4NphYt3usvJBXqVTeS", "YntncxzVBMHA+QNwgkdAQc95gSoWXrdf7r1T38+rtYkuCud5EkV7tWy0GDLeSKCvEuOpVtdPFlaVHXQyYNPwVAA=")
		b3, _ := addr.NewSignedMultiAddressFromString("/ip4/35.180.66.220/tcp/18514/ren/8MHEtUrZBQuRRtxAgBTTM6Zov3imfP", "Smw7e2DZ7nyPOr94oVuWmezAIDz1+uAkwkmQnH5nI0d+w6bECUV96wlJgvylILIL0cVITn5M1Mfm3fQeMWwVHQE=")
		b4, _ := addr.NewSignedMultiAddressFromString("/ip4/15.188.15.210/tcp/18514/ren/8MGob6LJcneeFSiQStU9FvP83W3xMA", "gcePwG3m0JonXpVZD8xJtWr7RnBsBuNViOkZElNAkaEDgYXmOdTBtXd1HMTjbwAilIG+hVlpIRpyC97n2F2z6gE=")
		b5, _ := addr.NewSignedMultiAddressFromString("/ip4/18.138.225.107/tcp/18514/ren/8MKUZzR3oM4ALnQ5vjQti1X41DwkEW", "iiiRnd/bkQlVVxr5XKrKi0uHhk4tAs5ct3rUjrYf8NAzhn49D/CAuVYFlxlmM31mCqOFOa1xO7HHZbV4zQJyggE=")
		b6, _ := addr.NewSignedMultiAddressFromString("/ip4/3.9.164.193/tcp/18514/ren/8MJSu4N1FgyT4ZYRH9faB9G6oMUUiF", "mmVxayVx6vyOaGR/r4hWezXhszGf9MC3OQxgBuc7kL8+OyvtUl6TwmzmXIXUoZqPRN07IwmGY11BUH1W43lsAgE=")
		b7, _ := addr.NewSignedMultiAddressFromString("/ip4/13.209.5.177/tcp/18514/ren/8MG7JhRuoj6SSQuzCWeWCdXRXd6Mn3", "CeHap47KDuH3oLcSUQOvtF/nK8HhjyejOqF+93FqhR0B2yRdOqHoHVUeY6Zsy2kOplo+2YNTpkzL9tDH2CpI2gE=")
		bootstraps = append(bootstraps, b1, b2, b3, b4, b5, b6, b7)
	case Testnet:
		b1, _ := addr.NewSignedMultiAddressFromString("/ip4/165.22.58.69/tcp/18514/ren/8MHjCu8ZiFaPShXx7SfJ93hpHRMLwv", "5QQTiQuEIWoVotCLDPd4BdjyJPKr+YDT6WsMYxps/q9a+ZwBZMRQbQscfnBQgUwEC2rPVVbTfZITjPKcbVIRJAA=")
		b2, _ := addr.NewSignedMultiAddressFromString("/ip4/165.22.193.227/tcp/18514/ren/8MJWSxiNmY3ghCYYo14yB1VPq7Su5h", "DEEByr1HSuV0gwyrZSh20Ym9RAP4J0iC+ErUJSpiojwG5TG1y9CI+umHD/gPXkjGdNfkB1wiplFYM6VSeRFidgE=")
		b3, _ := addr.NewSignedMultiAddressFromString("/ip4/68.183.0.112/tcp/18514/ren/8MGWsPMRNhgbCePqGU6Rk8SpUWrwLt", "XR4a6lVTnol221RgCsQa0Lj7aJ743iYtGd4wQvRBwuga1Oh6UEagsYplI1FZGshaLlwFsr68l21Hfb3DbR4JWgE=")
		b4, _ := addr.NewSignedMultiAddressFromString("/ip4/167.71.92.168/tcp/18514/ren/8MGPSZmCYeUk5iSu23DUAVDvyhdL7v", "9N8S+/Ar6zm3EqjisViGf3jW15K/1/exunjH1AG7GcpANiJdxUNOsYMADzZmjL6oKnVmZmU8iOtZdcN+cehrkwA=")
		b5, _ := addr.NewSignedMultiAddressFromString("/ip4/134.209.251.49/tcp/18514/ren/8MHMe1RP1sBQ5DUqhu1AGSt5Sc9ULQ", "R15UlnJazg6XH+VsdQ++wJJTqanpP9PXf0dzsFDbGJ5gNL+4EN7it53gbY1Aw2KWxZS+tcrSzhg32DF/PI+ABgA=")
		b6, _ := addr.NewSignedMultiAddressFromString("/ip4/167.99.115.41/tcp/18514/ren/8MJZ7FbRpvvnntf6yAC9Y6MacRnjjC", "IDOixOf/cgDTjPo7NlAln9O8I6psg182F0qttcGzFiE+LlK058n4sFtQQKp5Qb21gDCS2Q6O6P0la1pKrHrhbQE=")
		b7, _ := addr.NewSignedMultiAddressFromString("/ip4/15.188.51.176/tcp/18514/ren/8MGMnxzVuTESp8nMQGtXVjqLX7c54e", "1Px6FxTTIMTjLGfVMKo2YZA7zADfHz46DlyVZ03NBW4DuPawlhRvtUa6g5aOdUbRvvckwWzpDLcmvtAqlZTvxwA=")
		b8, _ := addr.NewSignedMultiAddressFromString("/ip4/34.207.81.121/tcp/18514/ren/8MKSGeXG3YRgc4VCdkKfGbbuxq9M4z", "0XT2oStrWmtVj6Mbp9rJP3L8kMB9yjO8MRr9uv3K7ckjC0VJ060GN4vLt16Ja6NIE31MD+hfDXlIM4fxW8RXUgA=")
		b9, _ := addr.NewSignedMultiAddressFromString("/ip4/35.183.106.97/tcp/18514/ren/8MGvowp18gG3qZsvDNEWBaPgSkMB8g", "ALKRx3nfDZ2e5qIs/pyfqDs9jqfi0ahkuwjYofQehyEMZO1M4uUltyVKQCJcWDPKNt7KwwIfLX9V0TDo8BKG3QE=")
		b10, _ := addr.NewSignedMultiAddressFromString("/ip4/54.206.68.198/tcp/18514/ren/8MKHfAVy5UY8E5DG9CwQdxcSTYx71L", "yRyiSNex2MfpsZqIcenGQFqASsiG4gG7GHlndM2aIUh7qDPjrK7zexnulf4lkO1tHJ7fNs0/fvyg7VLq6Ju90gE=")
		b11, _ := addr.NewSignedMultiAddressFromString("/ip4/52.47.59.114/tcp/18514/ren/8MK1Uw3YSiK6qPXwKuLYNiFR1f2ByX", "PvPnmEsqLZZ1MBR5q6OV7kId5HhtreND3sbsNAcm/AMZ/WnLphpVodqMoWPPZrM2HgD9Oyib6OFIzavj3YO0JQE=")
		b12, _ := addr.NewSignedMultiAddressFromString("/ip4/54.206.71.153/tcp/18514/ren/8MGbx8WcfkTaeHGX13VtDX4o2764R6", "zIyk9Q4bVyOOiJ9rT9ML1LK7vrxDYlgp/dGAkp0xxoMlV9oX0PLcw8le5d3iBHeYzXoFqdPV72TXSljH6cBSHQE=")
		b13, _ := addr.NewSignedMultiAddressFromString("/ip4/35.180.181.6/tcp/18514/ren/8MHGc7XQSFJqaKTaxiDnTMW496qBGL", "TosgLLMSj04rKrrnIpbkJoDN0E1UqfnLORuYs90GuMU8oUvBaI7dhNRLTfv0i7pyLCa2LSgsYNyWZ4bNOHSJ0QA=")
		bootstraps = append(bootstraps, b1, b2, b3, b4, b5, b6, b7, b8, b9, b10, b11, b12, b13)
	case Devnet:
		b1, _ := addr.NewSignedMultiAddressFromString("/ip4/165.22.219.22/tcp/18514/ren/8MJjieFyX3bXB7kqPdQJiuFroKMjCa", "oqvHbUhGhsnEwZ0iBI9NxTglMTirFehiD0r9D5scYXwCAiheWHaMrr433lhk8ZqsoDXpgKBFYi7hO0EykCM+JwE=")
		b2, _ := addr.NewSignedMultiAddressFromString("/ip4/159.203.177.223/tcp/18514/ren/8MGCtJ26XBZqLuYqshh14NHf1otmwc", "sfsJ9QGPA5anlU0uf7yG04musYqZgewEtak1uio4GRVtEhQ5kZUqoWW4Xj9Nmdypq3sE7+99XoojXmdS/XULewE=")
		b3, _ := addr.NewSignedMultiAddressFromString("/ip4/165.22.233.100/tcp/18514/ren/8MGMwMxVzjh54S5BoTMGkEnyEyDFj5", "XGGvArLgA83G3wgeS6ZAOioPOf/sOJ+5omLCuHUwhb0VOQtg1VgWb597BEVhxh+5XD99DBo2zxjCwG3XGrSoXQA=")
		b4, _ := addr.NewSignedMultiAddressFromString("/ip4/128.199.194.237/tcp/18514/ren/8MK9QTLA9XzptS9tdx5Je7BR5z9aen", "9oAvjlB1FJlskNM81I5ws9FsgrvqY+TWFy8wRoj19O1gHaxceR+PnYAN+Di/yEeHJ4D+qydX2h+a2N5E0CkhlgE=")
		b5, _ := addr.NewSignedMultiAddressFromString("/ip4/159.65.218.182/tcp/18514/ren/8MJ3r6FwD5inuoCfHNJdBKzsBLJsDC", "1hwlAT7B6AS2QNzqf2/nOwkVNEY4ocg1Ede10KQe5rRd4hXGCuvKzp64IGgXzkLqwTf2HA3sUyeHDmB8RByFhwA=")
		b6, _ := addr.NewSignedMultiAddressFromString("/ip4/142.93.172.163/tcp/18514/ren/8MHSgJ5eqebETYN7kd69gKakMZE9Ht", "U5A5bIze2cgO+SzqoMkGkihTOpRHoOZ6o7/yT30NsQJmTYCZfMBax05HPmSiDMz9kdzb9Ngfg2eFJqDEFttraQE=")
		b7, _ := addr.NewSignedMultiAddressFromString("/ip4/54.233.186.124/tcp/18514/ren/8MJVSvCFo2zj7UNsr4ZMh6h3Bi4iB9", "Rto4k5x4TxPplY8y0uO2ksnQBpJHcrVHYw0wH3H699tMV0nt2qBPaWC4Rx6sHF/oGp1P/zxodD34uBKhmjuwmgE=")
		b8, _ := addr.NewSignedMultiAddressFromString("/ip4/13.250.39.106/tcp/18514/ren/8MKVi5oAw1ULe9qpzMR8yLmGbSVeUb", "HEhzef0qKQgbhkJKeI51wXaJ/2OieifLc9czo2ANgow0ywRnIg0D5Up4zOQfLFGqKsJCP4UusPKMj8eqXObOwAE=")
		b9, _ := addr.NewSignedMultiAddressFromString("/ip4/35.183.132.51/tcp/18514/ren/8MJjhHRyk5QjZsSDvcCJAzjat98BVu", "LjvRDJZWxGr/F41qmwG+GEAgX19uS4TrnxGiwmDiWrtjpSh7jJbyrjyhUF7u7QjPGTI/gDL8QtLM0frVJtLk4gE=")
		b10, _ := addr.NewSignedMultiAddressFromString("/ip4/34.242.102.128/tcp/18514/ren/8MKbsQfDSNtnaH4JE7s2ewpCVgrSMJ", "d3xDdKFNeTX6cC297klU0euBLuUPbMyQ4MvIvMfZvrlkxXbalAU7zKucN4NzVhx+NpNUEwFffc+Vmih86XEingA=")
		b11, _ := addr.NewSignedMultiAddressFromString("/ip4/52.77.226.237/tcp/18514/ren/8MHDjCSqpD1Kc8qPpty3ApiJHmzMbR", "lkNpaqFbtlUlKWJsM5AZm4tiAERb4Gt8iZRzzXsNvUtd9WnthBXvIEYTW+3F5Ac83Nga/HnikosVimdiGYXu3AE=")
		b12, _ := addr.NewSignedMultiAddressFromString("/ip4/52.63.198.228/tcp/18514/ren/8MGx5J9M2qpUDvaGEqi4ABAb4ACKyx", "64thzydG6WjNeSOd11G5CoNp179vicGrzsjLR/tpCIxVP0YSO2RjdnyqwLNSaW9x1beOdRX3ISTxy28W2drcuwA=")
		b13, _ := addr.NewSignedMultiAddressFromString("/ip4/18.185.81.195/tcp/18514/ren/8MKT3RBR1ZHTMU6ha8pdDDF4TgyDCq", "O3ILS5uTz3TopMaGZpx9fvNBlTX3zx3CNo/WFzTQggV65N1evYeZkFflF6zu+rztEwLf3x9jwzdxVccs8HqfZQA=")
		bootstraps = append(bootstraps, b1, b2, b3, b4, b5, b6, b7, b8, b9, b10, b11, b12, b13)
	default:
		panic("unknown network")
	}
	return bootstraps
}

func (network Network) DnrAddress() common.Address {
	switch network {
	case Mainnet:
		panic("unimplemented")
	case Chaosnet:
		return common.HexToAddress("0xe0a0fc6ecaba77638b4f516c11b68f3837fe949a")
	case Testnet:
		return common.HexToAddress("0x39238ddCc6b301068f6D5B3DCa3A981EB20ba9A3")
	case Devnet:
		return common.HexToAddress("0x6E1a6b85f05bfec5c24C7a26E302cB28e639651c")
	default:
		panic("unknown network")
	}
}

func (network Network) ShiftRegistryAddress() common.Address {
	switch network {
	case Mainnet:
		panic("unimplemented")
	case Chaosnet:
		return common.HexToAddress("0x1d4713b74e79a3696722aebe05de976979876cac")
	case Testnet:
		return common.HexToAddress("0xbA563a8510d86dE95F5a50007E180d6d4966ad12")
	case Devnet:
		return common.HexToAddress("0xc7B310c18D78f175812CFfD8896d3cC959aC28d6")
	default:
		panic("unknown network")
	}
}
