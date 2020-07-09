/*
#Time      :  2019/12/2 11:02 AM 
#Author    :  chuangangshen@deepglint.com
#File      :  polygonTest.go
#Software  :  GoLand
*/
package main

import (
	"fmt"
	"github.com/deepglint/flowservice/datamodel"
)

func main() {
	// polygonInOrOut()
	//polygonIn()
	pathTest()
}

func pathTest() {
	pathSlice := []int{
		-551, 2864, 1032, -545, 2834, 1029, -537, 2798, 1029, -531, 2765, 1032,
		-526, 2740, 1040, -523, 2722, 1050, -521, 2710, 1060, -520, 2706, 1068,
		-521, 2709, 1069, -523, 2719, 1059, -527, 2757, 1042, -532, 2872, 1030,
		-549, 3036, 1050, -565, 3028, 1104, -522, 2665, 1113}
	for i := 1; i < len(pathSlice)/3; i++ {
		fmt.Println(datamodel.Trajectory(pathSlice[i*3:(i+1)*3]).Distance2D(pathSlice[(i-1)*3:i*3]) )
	}
}

/*func polygonIn() {
	var pg polygon.Polygon
	pg.Append(polygon.Point{0, 0})
	pg.Append(polygon.Point{0, 4})
	pg.Append(polygon.Point{4, 0})
	pg.Append(polygon.Point{4, 4})
	pIn := polygon.Point{X: 5, Y: 2}
	if pg.ContainsPoint(pIn, polygon.OddEvenFill) {
		fmt.Println("点pIn在多边形区域内")
	} else {
		fmt.Println("点pIn不在多边形区域内")
	}
}

func polygonInOrOut() {
	var pg polygon.Polygon
	pg.Append(polygon.Point{-8234.09, 3247.55})
	pg.Append(polygon.Point{-9207.86, -9216.69})
	pg.Append(polygon.Point{179.271, -14085.5})
	pg.Append(polygon.Point{10657, -9995.71})
	pg.Append(polygon.Point{11903.5, 2234.83})
	pg.Append(polygon.Point{2360.52, 7376.33})
	pg.Append(polygon.Point{-8234.09, 3247.55})
	// 此点在多边形内
	p1In := polygon.Point{-7212, -7941}
	// 此点不在多边形内
	p1Out := polygon.Point{-7455, -12956}
	if pg.ContainsPoint(p1In, polygon.OddEvenFill) {
		fmt.Println("点p1In在多边形区域内")
	} else {
		fmt.Println("点p1In不在多边形区域内")
	}
	if pg.ContainsPoint(p1Out, polygon.OddEvenFill) {
		fmt.Println("点p1Out在多边形区域内")
	} else {
		fmt.Println("点p1Out不在多边形区域内")
	}
}*/
