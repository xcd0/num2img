package main

import (
	"bufio"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"unsafe"
)

func remaped(matrix *[]uint16) *[][]uint16 {
	bit := int(unsafe.Sizeof((*matrix)[0])) * 8
	log.Print("bit : ", bit)
	l := len(*matrix) // l:入力行列数
	n := int(math.Ceil(math.Sqrt((float64(l) / float64(bit)))))
	num := bit * n                  // numには、matrixを正方形にしたときの１辺の長さが入る
	output := make([][]uint16, num) // outputには、matrixを1辺がnum*bit個の正方形に並び替えた結果が入る
	for i, _ := range output {
		output[i] = make([]uint16, num)
	}
	log.Print(n)
	log.Print(num)
	log.Print(num * bit)
	for i, v := range *matrix {
		//log.Printf("%v / %v = %v, %v", l, c, l%c, n/c*bit+bit)
		r := i / num
		c := i % num
		//log.Printf("i:%v,r:%v,c:%v", i, r, c)
		for b := 0; b < bit; b++ {
			output[r*bit+b][c] = (v >> b) & 1
		}
	}
	return &output
}

func GetInputMatrix() (*[]uint16, error) {
	var matrix []uint16 // 16x16のuint16型の行列を表現するためのスライス
	scanner := bufio.NewScanner(os.Stdin)
	floatToUint16 := func(f float32) uint16 {
		x := uint16(f * 65535) // 0から1の範囲を0から65535の範囲に変換
		if x > 65535 {         // オーバーフローを回避
			x = 65535
		}
		return x
	}
	for scanner.Scan() { // 標準入力から1行ずつ読み取る
		// 文字列からfloat32に変換し、float16に変換してからRaw()で16bitのuint16に変換する
		f, err := strconv.ParseFloat(scanner.Text(), 32)
		if err != nil {
			return nil, err
		}
		matrix = append(matrix, floatToUint16(float32(f)))
	}
	if err := scanner.Err(); err != nil { // エラーチェック
		return nil, err
	}
	return &matrix, nil
}

func count(matrix *[]uint16) *[]int {
	var max uint16 = 65535
	acc := make([]int, max)
	// カウントする
	for _, v := range *matrix {
		if v >= max {
			v = max - 1
		}
		acc[v]++
	}
	return &acc
}

func SaveBoolsToImage(acc *[]int, filename string) error {

	maxCount := func(a []int) int {
		sort.Sort(sort.IntSlice(a))
		return a[len(a)-1]
	}(*acc)

	log.Println("maxCount:", maxCount)

	bit := 8
	d := 1 << bit // 256 x 256
	numCols := d
	numRows := d
	log.Printf("numCols:%d,numRows:%d", numCols, numRows)

	img := image.NewGray16(image.Rect(0, 0, numCols, numRows))

	output := make([][]int, numRows) // outputには、matrixを1辺がnum*bit個の正方形に並び替えた結果が入る
	for i, _ := range output {
		output[i] = make([]int, numCols)
	}
	log.Printf("output:%d,%d", len(output[0]), len(output))

	for i, v := range *acc {
		x := i % d
		y := i / d // iが0個だった時の位置
		//log.Printf("i:%d,x:%d,y:%d,v:%d", i, x, y, v)
		u := float64(maxCount-v) / float64(maxCount) * 65535.
		img.SetGray16(x, y, color.Gray16{uint16(u)})
		output[x][y] = int(u)
	}

	/*
		for i, _ := range output {
			for _, v := range output[i] {
				fmt.Printf("%d", v)
			}
			fmt.Printf("\n")
		}
	*/

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, img) // PNG形式で画像を保存 errorがあれば返す
}

func SaveBitsToImage(bits *[][]uint16, filename string) error {
	numRows := len(*bits)
	numCols := len((*bits)[0])
	img := image.NewRGBA(image.Rect(0, 0, numCols, numRows)) // 画像の幅と高さを設定
	for row := 0; row < numRows; row++ {                     // uint16の各要素をビット列に変換し、画像にセット
		for col := 0; col < numCols; col++ {
			if (*bits)[row][col] == 1 {
				img.Set(col, row, color.Black)
			} else {
				img.Set(col, row, color.White)
			}
			/*
				bitIdx := col % 16
				elemIdx := col / 16
				if elemIdx < len((*bits)[row]) {
					bitVal := ((*bits)[row][elemIdx] >> bitIdx) & 1
					if bitVal == 1 {
						img.Set(col, row, color.Black)
					} else {
						img.Set(col, row, color.White)
					}
				}
			*/
		}
	}
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, img) // PNG形式で画像を保存 errorがあれば返す
}

func main() {
	//log.SetPrefix(fmt.Sprintf("[%s] ", cmdName))
	log.SetFlags(log.Ltime | log.Lshortfile)
	filename := "output.png"

	matrix, err := GetInputMatrix()
	if err != nil {
		log.Fatal(err)
	}
	if len(os.Args) > 0 {
		mat := remaped(matrix)
		if err := SaveBitsToImage(mat, filename); err != nil {
			log.Fatal(err)
		}
	} else {
		mat := count(matrix)
		if err := SaveBoolsToImage(mat, filename); err != nil {
			log.Fatal(err)
		}
	}

}
