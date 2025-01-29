package internal

import (
	"crypto/rand"
	"crypto/sha1"
	"fmt"
	"math/big"
)

func GenIntStd(min int, max int) int {
	rng := max - min
	nBig, err := rand.Int(rand.Reader, big.NewInt((int64)(rng)))
	if err != nil {
		panic(err)
	}
	n := nBig.Int64()
	return (int)(n)
}

func Get3DRandFloat(seed, d1, d2, d3 int64) float64 {
	var temp uint64 = Get3DRand(seed, d1, d2, d3)
	return float64(temp%1000001) / 1000000.0
}

func Get3DRand(seed, d1, d2, d3 int64) uint64 {
	var mass []byte = make([]byte, 32)

	for j := 0; j < 4; j++ {
		var temp uint64
		var val int64
		switch j {
		case 0:
			val = seed
		case 1:
			val = d1
		case 2:
			val = d2
		case 3:
			val = d3
		default:
			panic("что-то сильно пошло не так...")
		}
		if val < 0 {
			temp = uint64(val * -1)
			temp = temp << 1
			temp += 1
		} else {
			temp = uint64(val)
		}
		for i := 0; i < 4; i++ {
			mass[i+j*4] = byte((temp << (i * 2)) & 0xff)
		}
	}
	ShaMass := sha1.Sum(mass)
	return ChecksumToUint(ShaMass[:])
}

func ChecksumToUint(mass []byte) uint64 {
	var result uint64
	j := 0
	for i := 0; i < 20; i++ {
		result |= uint64(mass[i]) << (8 * j)
		j++
		if j == 8 {
			j = 0
		}
	}
	return result
}

func Str_To_Seed(text string) int64 {
	var g uint64
	fmt.Println()
	var mass [20]byte
	fmt.Println()
	mass = sha1.Sum([]byte(text))
	g = ChecksumToUint(mass[:])
	g >>= 1
	return int64(g)
}
