package main

import (
	"fmt"
	"mapgen/internal"
	"sync"
	"time"
)

func main() {
	var text string

	timestamp := time.Now()
	fmt.Println("generation map started")
	mass := generate(true)
	fmt.Printf("generation map ended, gen time=%v", time.Now().Sub(timestamp))
	timestamp = time.Now()
	fmt.Println()
	fmt.Println("image saving start")
	internal.GenerateImg(mass, "out.png")
	fmt.Printf("image saving end, gen time=%v", time.Now().Sub(timestamp))
	fmt.Println()
	fmt.Println("ok done")
	fmt.Println()

	fmt.Scan(&text)
}

func generate(log bool) [][]internal.MapPoint {
	fieldx, fieldy := internal.GetField()
	var mass [][]internal.MapPoint = ini_empty_mass(fieldx, fieldy)

	var map_gen_param_height internal.MapGenParam
	internal.SetDefaultHeight(&map_gen_param_height, 1)
	var map_gen_param_temp internal.MapGenParam
	internal.SetDefaultTemp(&map_gen_param_temp, 2)
	var map_gen_param_rain internal.MapGenParam
	internal.SetDefaultRain(&map_gen_param_rain, 3)

	seedstr := fmt.Sprint(internal.GenIntStd(10, 1000000))
	seed := internal.Str_To_Seed(seedstr)
	fmt.Printf("seed=%d", seed)
	fmt.Println()

	wg := new(sync.WaitGroup)

	for i := 0; i < len(mass); i++ {
		wg.Add(1)
		go calc(wg, map_gen_param_height, map_gen_param_temp, map_gen_param_rain, mass, seed, int64(i), log)

		//for j := 0; j < len(mass[0]); j++ {
		//	mass[i][j].SetHeight(int64(i), int64(j), map_gen_param, seed)
		//	mass[i][j].c = mass[i][j].GetColor()
		//}
		//if log && ((i % 32) == 0) {
		//	fmt.Printf("%d/%d (%f)", i+1, len(mass), 100.0*float64(i+1)/float64(len(mass)))
		//	fmt.Println()
		//}
	}
	wg.Wait()
	return mass
}

func calc(wg *sync.WaitGroup, map_gen_param_height, map_gen_param_temp, map_gen_param_rain internal.MapGenParam, mass [][]internal.MapPoint, seed, i int64, log bool) {
	defer wg.Done()
	for j := 0; j < len(mass[0]); j++ {
		mass[i][j].SetHeight(int64(i), int64(j), map_gen_param_height, seed)
		mass[i][j].SetTemp(int64(i), int64(j), map_gen_param_temp, seed)
		mass[i][j].SetRain(int64(i), int64(j), map_gen_param_rain, seed)
		mass[i][j].C = mass[i][j].GetColor()
	}
	if log { //&& ((i % 32) == 0) {
		//fmt.Printf("%d/%d (%f)", i+1, len(mass), 100.0*float64(i+1)/float64(len(mass)))
		fmt.Println("id=", i, " done")
	}
}

func ini_empty_mass(fieldx, fieldy int) [][]internal.MapPoint {
	var mass [][]internal.MapPoint = make([][]internal.MapPoint, fieldx)
	for i := 0; i < len(mass); i++ {
		mass[i] = make([]internal.MapPoint, fieldy)
	}
	return mass
}
