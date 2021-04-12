package str_rune

import (
	"testing"
)

//summary: 	纯ascii 用builder
//			unicode 3字节的 []rune转

const (
	ExampleStr1 = "1234567890"
	ExampleStr2 = "12312easfsdasdnlruewrhio3012938uhr345 43brj12ie7xzfcgbnmgu3812367e8vuihlkjnmhu87098fupoijhjbkhjvy2t7684709rfusijklnjkbndhjvfeuyt8769870938upijekhljkbhvjfuyr6t8769870938uiejkdnjbkhvjfyut78987098-0iorpjklnjkbjhvfuydt8ew79873098ueiorjpiklnjfkbvjgsdhfutewr67t8769870498-ritojkgnjbkvhjzsGFEUWR68T7Q9809JKGLNBXZVJHDGSFUTEW68T73948709RTXNFJBDHGSEWQYU38749"
	ExampleStr3 = "一二三四五六七八九十龘"
	ExampleStr4 = "12312eani 你打算动物企鹅 43brj12ie7xzfcgbnmgu3812367e8vuihlkjnmhu87098fupoijhjbkhjvy2t7684709rfusijklnjkbndhj的撒都强迫我没钱饿了；企鹅们去玩bhvjfuyr6t8769870938uiejkdnjbkhvjfyut78987098-0ior额外企鹅请问马卡罗夫是你的考拉网俄去r67t8769870498-ritojkgnjbkvhjzsGFEUWR68T7Q9809JKGLNBXZVJHDGSFUTEW68T7的撒昆都伦区难为情我WQYU38749"
)

func BenchmarkNormalTruncate(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		normalTruncate(ExampleStr1, 100)
		normalTruncate(ExampleStr2, 100)
		normalTruncate(ExampleStr3, 100)
		normalTruncate(ExampleStr4, 100)
	}
}

func BenchmarkBuildTruncate(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		buildTruncate(ExampleStr1, 100)
		buildTruncate(ExampleStr2, 100)
		buildTruncate(ExampleStr3, 100)
		buildTruncate(ExampleStr4, 100)
	}
}

func BenchmarkDecodeTruncate(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		decodeTruncate(ExampleStr1, 100)
		decodeTruncate(ExampleStr2, 100)
		decodeTruncate(ExampleStr3, 100)
		decodeTruncate(ExampleStr4, 100)
	}
}

func BenchmarkTruncateByWords(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		TruncateByWords(ExampleStr1, 100)
		TruncateByWords(ExampleStr2, 100)
		TruncateByWords(ExampleStr3, 100)
		TruncateByWords(ExampleStr4, 100)
	}
}

//func Test_TestLongLive(t *testing.T) {
//	t.Run("Diff", func(t *testing.T) {
//		fmt.Println(normalTruncate(ExampleStr1) == buildTruncate(ExampleStr1))
//		fmt.Println(normalTruncate(ExampleStr2) == buildTruncate(ExampleStr2))
//		fmt.Println(normalTruncate(ExampleStr3) == buildTruncate(ExampleStr3))
//		fmt.Println(normalTruncate(ExampleStr4) == buildTruncate(ExampleStr4))
//		fmt.Println(normalTruncate(ExampleStr4))
//		fmt.Println(buildTruncate(ExampleStr4))
//	})
//}
