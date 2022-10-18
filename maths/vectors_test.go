package maths

import "testing"

func BenchmarkVector2_Magnitude(b *testing.B) {
	v := Vector2{
		X: 3,
		Y: 4,
	}

	for n := 0; n < b.N; n++ {
		v.Magnitude()
	}
}
