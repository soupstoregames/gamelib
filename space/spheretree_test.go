package space_test

import (
	"github.com/soupstoregames/gamelib/maths"
	"github.com/soupstoregames/gamelib/space"
	"math/rand"
	"testing"
)

func BenchmarkNewSphereTree_Move(b *testing.B) {
	cases := map[string]struct {
		actors        int
		maxSphereSize float64
	}{
		"?actors=2000&maxSphereSize=100": {
			actors:        2000,
			maxSphereSize: 100,
		},
		"?actors=2000&maxSphereSize=250": {
			actors:        2000,
			maxSphereSize: 250,
		},
		"?actors=2000&maxSphereSize=500": {
			actors:        2000,
			maxSphereSize: 500,
		},
		"?actors=2000&maxSphereSize=750": {
			actors:        2000,
			maxSphereSize: 750,
		},
		"?actors=2000&maxSphereSize=1000": {
			actors:        2000,
			maxSphereSize: 1000,
		},
	}

	type actor struct {
		sphere  maths.Sphere
		entryID int
	}

	center := maths.Vector3{X: 4000, Y: 4000}
	for name, c := range cases {
		b.Run(name, func(b *testing.B) {
			st := space.NewSphereTree(maths.Sphere{Center: center, Radius: 6000}, c.maxSphereSize, c.maxSphereSize/10)

			var actors []actor
			for i := 0; i < c.actors; i++ {
				actor := actor{
					sphere: maths.Sphere{Center: maths.Vector3{X: rand.Float64() * 8000, Y: rand.Float64() * 8000}, Radius: 1},
				}
				actor.entryID = st.Insert(uint64(i), actor.sphere)
				actors = append(actors, actor)
			}

			// static actors
			for i := 0; i < c.actors; i++ {
				st.Insert(uint64(i+c.actors), maths.Sphere{Center: maths.Vector3{X: rand.Float64() * 8000, Y: rand.Float64() * 8000}, Radius: 1})
			}

			b.ReportAllocs()
			b.ResetTimer()

			for n := 0; n < b.N; n++ {
				for i := 0; i < c.actors; i++ {
					delta := center.Sub(actors[i].sphere.Center).Normalize()
					actors[i].sphere.Center = actors[i].sphere.Center.Add(delta)
					st.Move(actors[i].entryID, actors[i].sphere)
				}
				st.Integrate()
				st.Recompute()
			}
		})
	}
}

func BenchmarkNewSphereTree_Scan(b *testing.B) {
	cases := map[string]struct {
		actors int
	}{
		"?actors=10": {
			actors: 10,
		},
		"?actors=100": {
			actors: 100,
		},
		"?actors=1000": {
			actors: 1000,
		},
		"?actors=10000": {
			actors: 10000,
		},
	}

	type actor struct {
		sphere  maths.Sphere
		entryID int
	}

	center := maths.Vector3{X: 4000, Y: 4000}
	for name, c := range cases {
		b.Run(name, func(b *testing.B) {
			st := space.NewSphereTree(maths.Sphere{Center: center, Radius: 6000}, 1000, 100)

			for i := 0; i < c.actors; i++ {
				actor := actor{
					sphere: maths.Sphere{Center: maths.Vector3{X: rand.Float64() * 8000, Y: rand.Float64() * 8000}, Radius: 1},
				}
				actor.entryID = st.Insert(uint64(i), actor.sphere)
			}

			// static actors
			for i := 0; i < c.actors; i++ {
				st.Insert(uint64(i+c.actors), maths.Sphere{Center: maths.Vector3{X: rand.Float64() * 8000, Y: rand.Float64() * 8000}, Radius: 1})
			}

			b.ReportAllocs()
			b.ResetTimer()

			for n := 0; n < b.N; n++ {
				for i := 0; i < c.actors; i++ {
					st.Scan(maths.Sphere{Center: center, Radius: 100})
				}
			}
		})
	}

}
