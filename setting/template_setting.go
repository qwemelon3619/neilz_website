package setting

import (
	"math"
	"text/template"

	"github.com/gin-gonic/gin"
)

func AddTemplateFunction(r *gin.Engine) {
	funcMap := template.FuncMap{
		"mod": func(a, b int) int {
			return a % b
		},
		"ge": func(a, b int) bool {
			return a >= b
		},
		"le": func(a, b int) bool {
			return a <= b
		},
		"seq": func(start, end int) []int {
			var seq []int
			for i := start; i <= end; i++ {
				seq = append(seq, i)
			}
			return seq
		},
		"add": func(a, b int) int {
			return a + b
		},
		"sub": func(a, b int) int {
			return a - b
		},
		"max": func(a, b int) int {
			return int(math.Max(float64(a), float64(b)))
		},
	}
	r.SetFuncMap(funcMap)
}
