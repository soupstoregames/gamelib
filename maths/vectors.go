package maths

import "math"

type Vector2 struct {
	X float64
	Y float64
}

func (v Vector2) Add(v2 Vector2) Vector2 {
	v.X += v2.X
	v.Y += v2.Y
	return v
}

func (v Vector2) Sub(v2 Vector2) Vector2 {
	v.X -= v2.X
	v.Y -= v2.Y
	return v
}

func (v Vector2) Multiply(scalar float64) Vector2 {
	v.X *= scalar
	v.Y *= scalar
	return v
}

func (v Vector2) Magnitude() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func (v Vector2) Normalize() Vector2 {
	div := v.Magnitude()
	v.X /= div
	v.Y /= div
	return v
}

func (v Vector2) Distance(v2 Vector2) float64 {
	return v.Sub(v2).Magnitude()
}

func (v Vector2) Distance2(v2 Vector2) float64 {
	v2 = v.Sub(v2)
	return v.X*v.X + v.Y*v.Y
}

type Vector3 struct {
	X float64
	Y float64
	Z float64
}

func (v Vector3) Add(v2 Vector3) Vector3 {
	v.X += v2.X
	v.Y += v2.Y
	v.Z += v2.Z
	return v
}

func (v Vector3) Sub(v2 Vector3) Vector3 {
	v.X -= v2.X
	v.Y -= v2.Y
	v.Z -= v2.Z
	return v
}

func (v Vector3) Multiply(scalar float64) Vector3 {
	v.X *= scalar
	v.Y *= scalar
	v.Z *= scalar
	return v
}

func (v Vector3) Magnitude() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

func (v Vector3) Normalize() Vector3 {
	div := v.Magnitude()
	v.X /= div
	v.Y /= div
	v.Z /= div
	return v
}

func (v Vector3) Distance(v2 Vector3) float64 {
	return v.Sub(v2).Magnitude()
}

func (v Vector3) Distance2(v2 Vector3) float64 {
	v2 = v.Sub(v2)
	return v.X*v.X + v.Y*v.Y + v.Z*v.Z
}
