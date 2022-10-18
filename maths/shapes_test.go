package maths_test

import (
	"github.com/soupstoregames/gamelib/maths"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRectangle_Intersects(t *testing.T) {
	a := maths.Rectangle{
		X:      0,
		Y:      0,
		Width:  2,
		Height: 2,
	}
	b := maths.Rectangle{
		X:      1,
		Y:      0,
		Width:  2,
		Height: 2,
	}
	c := maths.Rectangle{
		X:      0,
		Y:      1,
		Width:  2,
		Height: 2,
	}
	d := maths.Rectangle{
		X:      1,
		Y:      1,
		Width:  2,
		Height: 2,
	}
	e := maths.Rectangle{
		X:      2,
		Y:      1,
		Width:  2,
		Height: 2,
	}

	assert.True(t, a.Intersects(b))
	assert.True(t, a.Intersects(c))
	assert.True(t, a.Intersects(d))
	assert.True(t, b.Intersects(a))
	assert.True(t, b.Intersects(c))
	assert.True(t, b.Intersects(d))
	assert.True(t, c.Intersects(a))
	assert.True(t, c.Intersects(b))
	assert.True(t, c.Intersects(d))
	assert.True(t, d.Intersects(a))
	assert.True(t, d.Intersects(b))
	assert.True(t, d.Intersects(c))

	assert.False(t, a.Intersects(e))
}

func TestRectangle_Merge(t *testing.T) {
	baseRect := maths.Rectangle{X: 1, Y: 1, Width: 1, Height: 1}

	cases := []struct {
		name     string
		rect     maths.Rectangle
		canMerge bool
		newRect  maths.Rectangle
	}{
		{
			name:     "horizontally distant",
			rect:     maths.Rectangle{X: 3, Y: 1, Width: 1, Height: 1},
			canMerge: false,
			newRect:  baseRect,
		},
		{
			name:     "horizontally adjacent left",
			rect:     maths.Rectangle{X: 0, Y: 1, Width: 1, Height: 1},
			canMerge: true,
			newRect:  maths.Rectangle{X: 0, Y: 1, Width: 2, Height: 1},
		},
		{
			name:     "horizontally adjacent right",
			rect:     maths.Rectangle{X: 2, Y: 1, Width: 1, Height: 1},
			canMerge: true,
			newRect:  maths.Rectangle{X: 1, Y: 1, Width: 2, Height: 1},
		},
		{
			name:     "vertically distant",
			rect:     maths.Rectangle{X: 1, Y: 3, Width: 1, Height: 1},
			canMerge: false,
			newRect:  baseRect,
		},
		{
			name:     "vertically adjacent up",
			rect:     maths.Rectangle{X: 1, Y: 0, Width: 1, Height: 1},
			canMerge: true,
			newRect:  maths.Rectangle{X: 1, Y: 0, Width: 1, Height: 2},
		},
		{
			name:     "vertically adjacent down",
			rect:     maths.Rectangle{X: 1, Y: 2, Width: 1, Height: 1},
			canMerge: true,
			newRect:  maths.Rectangle{X: 1, Y: 1, Width: 1, Height: 2},
		},
		{
			name:     "contains",
			rect:     maths.Rectangle{X: 1, Y: 1, Width: 1, Height: 1},
			canMerge: true,
			newRect:  maths.Rectangle{X: 1, Y: 1, Width: 1, Height: 1},
		},
		{
			name:     "contains with half margin",
			rect:     maths.Rectangle{X: 0, Y: 0, Width: 2, Height: 2},
			canMerge: true,
			newRect:  maths.Rectangle{X: 0, Y: 0, Width: 2, Height: 2},
		},
		{
			name:     "contains with margin",
			rect:     maths.Rectangle{X: 0, Y: 0, Width: 3, Height: 3},
			canMerge: true,
			newRect:  maths.Rectangle{X: 0, Y: 0, Width: 3, Height: 3},
		},
		{
			name:     "diagonally adjacent",
			rect:     maths.Rectangle{X: 0, Y: 0, Width: 1, Height: 1},
			canMerge: false,
			newRect:  maths.Rectangle{X: 1, Y: 1, Width: 1, Height: 1},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			newRect, ok := baseRect.Merge(c.rect)
			if assert.Equal(t, c.canMerge, ok) {
				assert.Equal(t, c.newRect, newRect)
			}
		})
	}
}

func BenchmarkRectangle_Merge(b *testing.B) {
	baseRect := maths.Rectangle{X: 1, Y: 1, Width: 1, Height: 1}

	cases := []struct {
		name     string
		rect     maths.Rectangle
		canMerge bool
		newRect  maths.Rectangle
	}{
		{
			name:     "horizontally distant",
			rect:     maths.Rectangle{X: 3, Y: 1, Width: 1, Height: 1},
			canMerge: false,
			newRect:  baseRect,
		},
		{
			name:     "horizontally adjacent left",
			rect:     maths.Rectangle{X: 0, Y: 1, Width: 1, Height: 1},
			canMerge: true,
			newRect:  maths.Rectangle{X: 0, Y: 1, Width: 2, Height: 1},
		},
		{
			name:     "horizontally adjacent right",
			rect:     maths.Rectangle{X: 2, Y: 1, Width: 1, Height: 1},
			canMerge: true,
			newRect:  maths.Rectangle{X: 1, Y: 1, Width: 2, Height: 1},
		},
		{
			name:     "vertically distant",
			rect:     maths.Rectangle{X: 1, Y: 3, Width: 1, Height: 1},
			canMerge: false,
			newRect:  baseRect,
		},
		{
			name:     "vertically adjacent up",
			rect:     maths.Rectangle{X: 1, Y: 0, Width: 1, Height: 1},
			canMerge: true,
			newRect:  maths.Rectangle{X: 1, Y: 0, Width: 1, Height: 2},
		},
		{
			name:     "vertically adjacent down",
			rect:     maths.Rectangle{X: 1, Y: 2, Width: 1, Height: 1},
			canMerge: true,
			newRect:  maths.Rectangle{X: 1, Y: 1, Width: 1, Height: 2},
		},
		{
			name:     "contains",
			rect:     maths.Rectangle{X: 1, Y: 1, Width: 1, Height: 1},
			canMerge: true,
			newRect:  maths.Rectangle{X: 1, Y: 1, Width: 1, Height: 1},
		},
		{
			name:     "contains with half margin",
			rect:     maths.Rectangle{X: 0, Y: 0, Width: 2, Height: 2},
			canMerge: true,
			newRect:  maths.Rectangle{X: 0, Y: 0, Width: 2, Height: 2},
		},
		{
			name:     "contains with margin",
			rect:     maths.Rectangle{X: 0, Y: 0, Width: 3, Height: 3},
			canMerge: true,
			newRect:  maths.Rectangle{X: 0, Y: 0, Width: 3, Height: 3},
		},
		{
			name:     "diagonally adjacent",
			rect:     maths.Rectangle{X: 0, Y: 0, Width: 1, Height: 1},
			canMerge: false,
			newRect:  maths.Rectangle{X: 1, Y: 1, Width: 1, Height: 1},
		},
	}
	b.ReportAllocs()
	for _, c := range cases {
		b.Run(c.name, func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				baseRect.Merge(c.rect)
			}
		})
	}
}

func TestSphere_IntersectsSphere(t *testing.T) {
	cases := map[string]struct {
		s1       maths.Sphere
		s2       maths.Sphere
		expected bool
	}{
		"out of square distance": {
			s1:       maths.NewSphere(maths.Vector3{}, 3),
			s2:       maths.NewSphere(maths.Vector3{X: 5}, 2),
			expected: false,
		},
		"intersecting": {
			s1:       maths.NewSphere(maths.Vector3{}, 3),
			s2:       maths.NewSphere(maths.Vector3{X: 5}, 2.1),
			expected: true,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, c.expected, c.s1.IntersectsSphere(c.s2))
		})
	}
}

func BenchmarkSphere_IntersectsSphere(b *testing.B) {
	cases := map[string]struct {
		s1 maths.Sphere
		s2 maths.Sphere
	}{
		"out of square distance": {
			s1: maths.NewSphere(maths.Vector3{}, 3),
			s2: maths.NewSphere(maths.Vector3{X: 5}, 2),
		},
		"intersecting": {
			s1: maths.NewSphere(maths.Vector3{}, 3),
			s2: maths.NewSphere(maths.Vector3{X: 5}, 2.1),
		},
	}

	for name, c := range cases {
		b.Run(name, func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				_ = c.s1.IntersectsSphere(c.s2)
			}
		})
	}
}
