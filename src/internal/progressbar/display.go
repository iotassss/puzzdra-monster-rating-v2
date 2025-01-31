package progressbar

import "fmt"

// プログレスバーを表示する関数
func Display(current, total int) {
	barLength := 100 // プログレスバーの長さ
	filled := (current * barLength) / total
	empty := barLength - filled

	// 最後は改行して終わる
	if current == total {
		fmt.Printf("\r%s%s(%d/%d)\n", repeat("■", filled), repeat("□", empty), current, total)
		return
	}

	// プログレスバーの表示
	fmt.Printf("\r%s%s(%d/%d)", repeat("■", filled), repeat("□", empty), current, total)
}

// 繰り返し文字列を生成する関数
func repeat(char string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += char
	}
	return result
}
