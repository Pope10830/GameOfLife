package main

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test(t *testing.T) {
	type args struct {
		p             golParams
		expectedAlive []cell
	}
	tests := []struct {
		name string
		args args
	}{
		{"16x16x2-0", args{
			p: golParams{
				turns:       0,
				threads:     2,
				imageWidth:  16,
				imageHeight: 16,
			},
			expectedAlive: []cell{
				{x: 4, y: 5},
				{x: 5, y: 6},
				{x: 3, y: 7},
				{x: 4, y: 7},
				{x: 5, y: 7},
			},
		}},

		{"16x16x4-0", args{
			p: golParams{
				turns:       0,
				threads:     4,
				imageWidth:  16,
				imageHeight: 16,
			},
			expectedAlive: []cell{
				{x: 4, y: 5},
				{x: 5, y: 6},
				{x: 3, y: 7},
				{x: 4, y: 7},
				{x: 5, y: 7},
			},
		}},

		{"16x16x8-0", args{
			p: golParams{
				turns:       0,
				threads:     8,
				imageWidth:  16,
				imageHeight: 16,
			},
			expectedAlive: []cell{
				{x: 4, y: 5},
				{x: 5, y: 6},
				{x: 3, y: 7},
				{x: 4, y: 7},
				{x: 5, y: 7},
			},
		}},

		{"16x16x2-1", args{
			p: golParams{
				turns:       1,
				threads:     2,
				imageWidth:  16,
				imageHeight: 16,
			},
			expectedAlive: []cell{
				{x: 3, y: 6},
				{x: 5, y: 6},
				{x: 4, y: 7},
				{x: 5, y: 7},
				{x: 4, y: 8},
			},
		}},

		{"16x16x4-1", args{
			p: golParams{
				turns:       1,
				threads:     4,
				imageWidth:  16,
				imageHeight: 16,
			},
			expectedAlive: []cell{
				{x: 3, y: 6},
				{x: 5, y: 6},
				{x: 4, y: 7},
				{x: 5, y: 7},
				{x: 4, y: 8},
			},
		}},

		{"16x16x8-1", args{
			p: golParams{
				turns:       1,
				threads:     8,
				imageWidth:  16,
				imageHeight: 16,
			},
			expectedAlive: []cell{
				{x: 3, y: 6},
				{x: 5, y: 6},
				{x: 4, y: 7},
				{x: 5, y: 7},
				{x: 4, y: 8},
			},
		}},

		{"16x16x2-100", args{
			p: golParams{
				turns:       100,
				threads:     2,
				imageWidth:  16,
				imageHeight: 16,
			},
			expectedAlive: []cell{
				{x: 12, y: 0},
				{x: 13, y: 0},
				{x: 14, y: 0},
				{x: 13, y: 14},
				{x: 14, y: 15},
			},
		}},

		{"16x16x4-100", args{
			p: golParams{
				turns:       100,
				threads:     4,
				imageWidth:  16,
				imageHeight: 16,
			},
			expectedAlive: []cell{
				{x: 12, y: 0},
				{x: 13, y: 0},
				{x: 14, y: 0},
				{x: 13, y: 14},
				{x: 14, y: 15},
			},
		}},

		{"16x16x8-100", args{
			p: golParams{
				turns:       100,
				threads:     8,
				imageWidth:  16,
				imageHeight: 16,
			},
			expectedAlive: []cell{
				{x: 12, y: 0},
				{x: 13, y: 0},
				{x: 14, y: 0},
				{x: 13, y: 14},
				{x: 14, y: 15},
			},
		}},

		// Special test to be used to generate traces - not a real test
		//{"trace", args{
		//	p: golParams{
		//		turns:       10,
		//		threads:     4,
		//		imageWidth:  64,
		//		imageHeight: 64,
		//	},
		//}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			alive := gameOfLife(test.args.p, nil)
			//fmt.Println("Ran test:", test.name)
			if test.name != "trace" {
				assert.ElementsMatch(t, alive, test.args.expectedAlive)
			}
		})
	}
}

const benchLength = 1000

func Benchmark(b *testing.B) {
	benchmarks := []struct {
		name string
		p    golParams
	}{
		{
			"16x16x2", golParams{
			turns:       benchLength,
			threads:     2,
			imageWidth:  16,
			imageHeight: 16,
		}},

		{
			"16x16x4", golParams{
			turns:       benchLength,
			threads:     4,
			imageWidth:  16,
			imageHeight: 16,
		}},

		{
			"16x16x8", golParams{
			turns:       benchLength,
			threads:     8,
			imageWidth:  16,
			imageHeight: 16,
		}},

		{
			"64x64x2", golParams{
			turns:       benchLength,
			threads:     2,
			imageWidth:  64,
			imageHeight: 64,
		}},

		{
			"64x64x4", golParams{
			turns:       benchLength,
			threads:     4,
			imageWidth:  64,
			imageHeight: 64,
		}},

		{
			"64x64x8", golParams{
			turns:       benchLength,
			threads:     8,
			imageWidth:  64,
			imageHeight: 64,
		}},

		{
			"128x128x2", golParams{
			turns:       benchLength,
			threads:     2,
			imageWidth:  128,
			imageHeight: 128,
		}},

		{
			"128x128x4", golParams{
			turns:       benchLength,
			threads:     4,
			imageWidth:  128,
			imageHeight: 128,
		}},

		{
			"128x128x8", golParams{
			turns:       benchLength,
			threads:     8,
			imageWidth:  128,
			imageHeight: 128,
		}},

		{
			"256x256x2", golParams{
			turns:       benchLength,
			threads:     2,
			imageWidth:  256,
			imageHeight: 256,
		}},

		{
			"256x256x4", golParams{
			turns:       benchLength,
			threads:     4,
			imageWidth:  256,
			imageHeight: 256,
		}},

		{
			"256x256x8", golParams{
			turns:       benchLength,
			threads:     8,
			imageWidth:  256,
			imageHeight: 256,
		}},

		{
			"512x512x2", golParams{
			turns:       benchLength,
			threads:     2,
			imageWidth:  512,
			imageHeight: 512,
		}},

		{
			"512x512x4", golParams{
			turns:       benchLength,
			threads:     4,
			imageWidth:  512,
			imageHeight: 512,
		}},

		{
			"512x512x8", golParams{
			turns:       benchLength,
			threads:     8,
			imageWidth:  512,
			imageHeight: 512,
		}},
	}
	for _, bm := range benchmarks {
		os.Stdout = nil // Disable all program output apart from benchmark results
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				gameOfLife(bm.p, nil)
				//fmt.Println("Ran bench:", bm.name)
			}
		})
	}
}
