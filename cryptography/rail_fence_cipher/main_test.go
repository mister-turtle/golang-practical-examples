package main

import (
	"testing"
)

func TestDecrypt(t *testing.T) {
	type args struct {
		ciphertext string
		rails      int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"3 chars, 3 rails", args{"abc", 3}, "abc"},
		{"4 chars, 3 rails", args{"abdc", 3}, "abcd"},
		{"5 chars, 3 rails", args{"aebdc", 3}, "abcde"},
		{"6 chars, 3 rails", args{"aebdfc", 3}, "abcdef"},
		{"7 chars, 3 rails", args{"aebdfcg", 3}, "abcdefg"},
		{"8 chars, 3 rails", args{"aebdfhcg", 3}, "abcdefgh"},
		{"4 chars, 4 rails", args{"abcd", 4}, "abcd"},
		{"5 chars, 4 rails", args{"abced", 4}, "abcde"},
		{"6 chars, 4 rails", args{"abfced", 4}, "abcdef"},
		{"7 chars, 4 rails", args{"agbfced", 4}, "abcdefg"},
		{"8 chars, 4 rails", args{"agbfhced", 4}, "abcdefgh"},
		{"9 chars, 4 rails", args{"agbfhceid", 4}, "abcdefghi"},
		{"10 chars, 4 rails", args{"agbfhceidj", 4}, "abcdefghij"},
		{"20 chars, 4 rails", args{"agmsbfhlnrtceikoqdjp", 4}, "abcdefghijklmnopqrst"},
		{"20 chars, 13 rails", args{"abcdeftgshriqjpkolnm", 13}, "abcdefghijklmnopqrst"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Decrypt(tt.args.ciphertext, tt.args.rails); got != tt.want {
				t.Errorf("QuickDecode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNaiveDecrypt(t *testing.T) {
	type args struct {
		ciphertext string
		rails      int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"3 chars, 3 rails", args{"abc", 3}, "abc"},
		{"4 chars, 3 rails", args{"abdc", 3}, "abcd"},
		{"5 chars, 3 rails", args{"aebdc", 3}, "abcde"},
		{"6 chars, 3 rails", args{"aebdfc", 3}, "abcdef"},
		{"7 chars, 3 rails", args{"aebdfcg", 3}, "abcdefg"},
		{"8 chars, 3 rails", args{"aebdfhcg", 3}, "abcdefgh"},
		{"4 chars, 4 rails", args{"abcd", 4}, "abcd"},
		{"5 chars, 4 rails", args{"abced", 4}, "abcde"},
		{"6 chars, 4 rails", args{"abfced", 4}, "abcdef"},
		{"7 chars, 4 rails", args{"agbfced", 4}, "abcdefg"},
		{"8 chars, 4 rails", args{"agbfhced", 4}, "abcdefgh"},
		{"9 chars, 4 rails", args{"agbfhceid", 4}, "abcdefghi"},
		{"10 chars, 4 rails", args{"agbfhceidj", 4}, "abcdefghij"},
		{"20 chars, 4 rails", args{"agmsbfhlnrtceikoqdjp", 4}, "abcdefghijklmnopqrst"},
		{"20 chars, 13 rails", args{"abcdeftgshriqjpkolnm", 13}, "abcdefghijklmnopqrst"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NaiveDecrypt(tt.args.ciphertext, tt.args.rails); got != tt.want {
				t.Errorf("QuickDecode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncrypt(t *testing.T) {
	type args struct {
		plaintext string
		rails     int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"0 chars", args{"", 3}, ""},
		{"3 chars, 3 rails", args{"abc", 3}, "abc"},
		{"4 chars, 3 rails", args{"abcd", 3}, "abdc"},
		{"5 chars, 3 rails", args{"abcde", 3}, "aebdc"},
		{"6 chars, 3 rails", args{"abcdef", 3}, "aebdfc"},
		{"7 chars, 3 rails", args{"abcdefe", 3}, "aebdfce"},
		{"4 chars, 4 rails", args{"abcd", 4}, "abcd"},
		{"5 chars, 4 rails", args{"abcde", 4}, "abced"},
		{"6 chars, 4 rails", args{"abcdef", 4}, "abfced"},
		{"7 chars, 4 rails", args{"abcdefg", 4}, "agbfced"},
		{"8 chars, 4 rails", args{"abcdefgh", 4}, "agbfhced"},
		{"9 chars, 4 rails", args{"abcdefghi", 4}, "agbfhceid"},
		{"10 chars, 4 rails", args{"abcdefghij", 4}, "agbfhceidj"},
		{"20 chars, 4 rails", args{"abcdefghijklmnopqrst", 4}, "agmsbfhlnrtceikoqdjp"},
		{"20 chars, 13 rails", args{"abcdefghijklmnopqrst", 13}, "abcdeftgshriqjpkolnm"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Encrypt(tt.args.plaintext, tt.args.rails); got != tt.want {
				t.Errorf("Encrypt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func benchmarkNaiveDecrypt(plain string, rails int, b *testing.B) {
	for n := 0; n < b.N; n++ {
		NaiveDecrypt(plain, rails)
	}
}
func benchmarkDecrypt(plain string, rails int, b *testing.B) {
	for n := 0; n < b.N; n++ {
		Decrypt(plain, rails)
	}
}

func BenchmarkNaiveDecrypt3Rails(b *testing.B)  { benchmarkNaiveDecrypt("agmsbfhlnrtceikoqdjp", 3, b) }
func BenchmarkNaiveDecrypt4Rails(b *testing.B)  { benchmarkNaiveDecrypt("agmsbfhlnrtceikoqdjp", 4, b) }
func BenchmarkNaiveDecrypt8Rails(b *testing.B)  { benchmarkNaiveDecrypt("aobnpcmqdlreksfjtgih", 8, b) }
func BenchmarkNaiveDecrypt16Rails(b *testing.B) { benchmarkNaiveDecrypt("abcdefghijkltmsnroqp", 16, b) }
func BenchmarkNaiveDecrypt32Rails(b *testing.B) {
	benchmarkNaiveDecrypt("this is a longer sentences rwaihtch  93", 32, b)
}
func BenchmarkNaiveDecrypt64Rails(b *testing.B) {
	benchmarkNaiveDecrypt("t acornchhe eat,ye hsi  uhlne etneta a vnmr hrces ab vn7ism gse tseocarmee5", 64, b)
}

func BenchmarkDecrypt3Rails(b *testing.B)  { benchmarkDecrypt("agmsbfhlnrtceikoqdjp", 3, b) }
func BenchmarkDecrypt4Rails(b *testing.B)  { benchmarkDecrypt("agmsbfhlnrtceikoqdjp", 4, b) }
func BenchmarkDecrypt8Rails(b *testing.B)  { benchmarkDecrypt("aobnpcmqdlreksfjtgih", 8, b) }
func BenchmarkDecrypt16Rails(b *testing.B) { benchmarkDecrypt("abcdefghijkltmsnroqp", 16, b) }
func BenchmarkDecrypt32Rails(b *testing.B) {
	benchmarkDecrypt("this is a longer sentences rwaihtch  93", 32, b)
}
func BenchmarkDecrypt64Rails(b *testing.B) {
	benchmarkDecrypt("t acornchhe eat,ye hsi  uhlne etneta a vnmr hrces ab vn7ism gse tseocarmee5", 64, b)
}
