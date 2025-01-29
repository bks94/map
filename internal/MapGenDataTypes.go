package internal

import (
	"image/color"
)

type MapPoint struct {
	height float64
	temp   float64
	rain   float64

	C color.RGBA
}

func (p *MapPoint) SetHeight(x, y int64, genparam MapGenParam, seed int64) {
	var tempf float64
	tempf = 0
	p.height = 0
	for s := 0; s < len(genparam.Layers); s++ {
		p.height += GetPoint(int64(x), int64(y), genparam.Layers[s], seed) * genparam.Layers[s].Weight
		tempf += genparam.Layers[s].Weight
	}
	p.height /= tempf
}

func (p *MapPoint) SetTemp(x, y int64, genparam MapGenParam, seed int64) {
	var tempf float64
	tempf = 0
	p.temp = 0
	for s := 0; s < len(genparam.Layers); s++ {
		p.temp += GetPoint(int64(x), int64(y), genparam.Layers[s], seed) * genparam.Layers[s].Weight
		tempf += genparam.Layers[s].Weight
	}
	p.temp /= tempf
}

func (p *MapPoint) SetRain(x, y int64, genparam MapGenParam, seed int64) {
	var tempf float64
	tempf = 0
	p.rain = 0
	for s := 0; s < len(genparam.Layers); s++ {
		p.rain += GetPoint(int64(x), int64(y), genparam.Layers[s], seed) * genparam.Layers[s].Weight
		tempf += genparam.Layers[s].Weight
	}
	p.rain /= tempf
}

func (p MapPoint) GetColor() color.RGBA {
	//Temp := byte(point.height * 255.0)
	//return color.RGBA{Temp, Temp, Temp, 255}

	if p.height < 0.2 {
		return color.RGBA{32, 32, 196, 255} //глубокая вода
	} else if p.height < 0.43 {
		return color.RGBA{0, 0, 255, 255} //мелкая вода
	} else if p.height < 0.45 {
		return color.RGBA{255, 255, 0, 255} //песок
	} else if p.height < 0.72 {
		if p.temp < 0.3 {
			return color.RGBA{255, 255, 255, 255} //снег
		}
		if p.rain > 0.70 {
			return color.RGBA{92, 160, 80, 255} //болото
		}
		if p.rain+(p.temp*0.8)+0.2 < 1.0 {
			return color.RGBA{206, 184, 39, 255} //пустыня
		}
		return color.RGBA{0, 255, 0, 255} //лес
	} else if p.height < 0.78 {
		return color.RGBA{96, 96, 96, 255} //горы
	} else {
		return color.RGBA{255, 255, 255, 255} //снежные горы
	}
}

type MapGenParam struct {
	Layers       []MapGenLayer
	IsIni        bool
	TotalWeights float64
	LayerId      int64
}

func ini(param *MapGenParam) {
	if param.IsIni {
		return
	}
	param.TotalWeights = 0
	param.IsIni = true
	for i := 0; i < len(param.Layers); i++ {
		param.TotalWeights += param.Layers[i].Weight
	}
}

func SetDefaultHeight(param *MapGenParam, LID int64) {
	param.Layers = make([]MapGenLayer, 6)
	param.LayerId = LID
	var grid int64
	var weight float64
	grid = 4
	weight = 1.0

	for i := 0; i < len(param.Layers); i++ {
		param.Layers[i] = MapGenLayer{grid, weight, param.LayerId}
		grid *= 2
		weight *= 2
	}
	ini(param)
}

func SetDefaultTemp(param *MapGenParam, LID int64) {
	param.Layers = make([]MapGenLayer, 4)
	param.LayerId = LID
	var grid int64
	var weight float64
	grid = 64
	weight = 1.0

	for i := 0; i < len(param.Layers); i++ {
		param.Layers[i] = MapGenLayer{grid, weight, param.LayerId}
		grid *= 2
		weight *= 2
	}
	ini(param)
}

func SetDefaultRain(param *MapGenParam, LID int64) {
	param.Layers = make([]MapGenLayer, 6)
	param.LayerId = LID
	var grid int64
	var weight float64
	grid = 16
	weight = 1.0

	for i := 0; i < len(param.Layers); i++ {
		param.Layers[i] = MapGenLayer{grid, weight, param.LayerId}
		grid *= 2
		weight *= 2
	}
	ini(param)
}

func GetPoint(x, y int64, params MapGenLayer, seed int64) float64 {
	dx := x % params.Grid
	var xmod, ymod float64
	var xmin, xmax int64
	if dx < 0 {
		xmin = x + dx
		xmod = 1.0 - (float64(dx) / float64(params.Grid))
	} else {
		xmin = x - dx
		xmod = (float64(dx) / float64(params.Grid))
	}
	xmin /= params.Grid
	xmax = xmin + 1
	dy := y % params.Grid
	var ymin, ymax int64
	if dy < 0 {
		ymin = y + dy
		ymod = 1.0 - (float64(dy) / float64(params.Grid))
	} else {
		ymin = y - dy
		ymod = (float64(dy) / float64(params.Grid))
	}
	ymin /= params.Grid
	ymax = ymin + 1

	//fmt.Printf("x=%d ,xmin=%d ,xmax=%d,xmod=%f y=%d, ymin=%d, ymax=%d ymod=%f", x, xmin, xmax, xmod, y, ymin, ymax, ymod)
	//fmt.Println()

	var grid [2][2]float64
	grid[0][0] = Get3DRandFloat(seed, params.LayerId, xmin, ymin)
	grid[1][0] = Get3DRandFloat(seed, params.LayerId, xmax, ymin)
	grid[0][1] = Get3DRandFloat(seed, params.LayerId, xmin, ymax)
	grid[1][1] = Get3DRandFloat(seed, params.LayerId, xmax, ymax)
	var result float64

	var closeness [2][2]float64
	closeness[0][0] = (1.0 - xmod) * (1.0 - ymod)
	closeness[1][0] = xmod * (1.0 - ymod)
	closeness[0][1] = (1.0 - xmod) * ymod
	closeness[1][1] = xmod * ymod
	//result += grid[0][0] * (xmod + ymod) / 2.0
	//result += grid[0][1] * (xmod + (1.0 - ymod)) / 2.0
	//result += grid[1][0] * ((1.0 - xmod) + ymod) / 2.0
	//result += grid[1][1] * ((1.0 - xmod) + (1.0 - ymod)) / 2.0
	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			result += grid[i][j] * closeness[i][j]
		}
	}
	return result
}

func Lerp(v0, v1, x float64) float64 {
	D := v0 - v1
	return v0 + D*x
}

type MapGenLayer struct {
	Grid    int64
	Weight  float64
	LayerId int64
}
