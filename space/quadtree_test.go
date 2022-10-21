package space_test

import (
	"github.com/soupstoregames/gamelib/maths"
	"github.com/soupstoregames/gamelib/space"
	"math/rand"
	"testing"
)

func BenchmarkQuadTree_MoveAndProcess(b *testing.B) {
	cases := map[string]struct {
		actors        int
		maxSphereSize float64
	}{
		"?actors=1000": {
			actors: 1000,
		},
		"?actors=2000": {
			actors: 2000,
		},
		"?actors=3000": {
			actors: 3000,
		},
		"?actors=4000": {
			actors: 4000,
		},
		"?actors=5000": {
			actors: 5000,
		},
	}

	type actor struct {
		rect maths.Rectangle
		id   uint64
	}

	center := maths.Vector2{X: 4000, Y: 4000}
	for name, c := range cases {
		b.Run(name, func(b *testing.B) {
			rand.Seed(1)
			qt := space.NewQuadTree(maths.Rectangle{Width: 8000, Height: 8000})

			var actors []actor
			for i := 0; i < c.actors; i++ {
				actor := actor{
					id:   uint64(i),
					rect: maths.Rectangle{X: rand.Float64() * 8000, Y: rand.Float64() * 8000, Width: 1, Height: 1},
				}
				qt.Insert(space.QuadTreeEntry{
					ID:   uint64(i),
					Rect: actor.rect,
				})
				actors = append(actors, actor)
			}

			b.ReportAllocs()
			b.ResetTimer()

			for n := 0; n < b.N; n++ {
				for i := 0; i < c.actors; i++ {
					qt.Remove(space.QuadTreeEntry{
						ID:   actors[i].id,
						Rect: actors[i].rect,
					})

					delta := center.Sub(maths.Vector2{X: actors[i].rect.X, Y: actors[i].rect.Y}).Normalize().Multiply(0.03)
					actors[i].rect.X += delta.X
					actors[i].rect.Y += delta.Y

					qt.Insert(space.QuadTreeEntry{
						ID:   actors[i].id,
						Rect: actors[i].rect,
					})
				}
			}
		})
	}
}

func BenchmarkQuadTree_Scan(b *testing.B) {
	cases := map[string]struct {
		actors int
	}{
		"?actors=1000": {
			actors: 1000,
		},
		"?actors=2000": {
			actors: 2000,
		},
		"?actors=3000": {
			actors: 3000,
		},
		"?actors=4000": {
			actors: 4000,
		},
		"?actors=5000": {
			actors: 5000,
		},
	}

	type actor struct {
		rect maths.Rectangle
		id   uint64
	}

	center := maths.Vector3{X: 4000, Y: 4000}
	for name, c := range cases {
		b.Run(name, func(b *testing.B) {
			rand.Seed(1)
			qt := space.NewQuadTree(maths.Rectangle{Width: 8000, Height: 8000})

			var actors []actor
			for i := 0; i < c.actors; i++ {
				actor := actor{
					id:   uint64(i),
					rect: maths.Rectangle{X: rand.Float64() * 8000, Y: rand.Float64() * 8000, Width: 1, Height: 1},
				}
				qt.Insert(space.QuadTreeEntry{
					ID:   uint64(i),
					Rect: actor.rect,
				})
				actors = append(actors, actor)
			}

			b.ReportAllocs()
			b.ResetTimer()

			var entries []space.QuadTreeEntry
			for n := 0; n < b.N; n++ {
				for i := 0; i < c.actors; i++ {
					qt.Scan(&entries, maths.Rectangle{X: center.X - 50, Y: center.Y - 50, Width: 100, Height: 100})
					entries = entries[:0]
				}
			}
		})
	}
}
