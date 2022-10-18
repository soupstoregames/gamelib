package maths

type Rectangle struct {
	X      float64
	Y      float64
	Width  float64
	Height float64
}

func (r Rectangle) ContainsVec(v Vector2) bool {
	return r.X <= v.X && r.Y <= v.Y && r.X+r.Width > v.X && r.Y+r.Height > v.Y
}

func (r Rectangle) ContainsRect(r2 Rectangle) bool {
	return r.X <= r2.X && r.X+r.Width >= r2.X+r2.Width && r.Y <= r2.Y && r.Y+r.Height >= r2.Y+r2.Height
}

func (r Rectangle) Intersects(r2 Rectangle) bool {
	return r.X < r2.X+r2.Width && r.X+r.Width > r2.X && r.Y < r2.Y+r2.Height && r.Y+r.Height > r2.Y
}

func (r Rectangle) Merge(r2 Rectangle) (Rectangle, bool) {
	if r.ContainsRect(r2) {
		return r, true
	}
	if r2.ContainsRect(r) {
		return r2, true
	}

	// is it vertically aligned?
	if r.X == r2.X && r.Width == r2.Width {
		// if the first rect is above the second
		if r.Y < r2.Y {
			// if the first rect touches the second
			if r.Y+r.Height >= r2.Y {
				r.Height = r2.Y + r2.Height - r.Y
				return r, true
			}
			return r, false
		} else {
			if r2.Y+r2.Height >= r.Y {
				r2.Height = r.Y + r.Height - r2.Y
				return r2, true
			}
			return r, false
		}
	}

	// is it horizontally aligned?
	if r.Y == r2.Y && r.Height == r2.Height {
		// if the first rect is left of the second
		if r.X < r2.X {
			// if the first rect touches the second
			if r.X+r.Width >= r2.X {
				r.Width = r2.X + r2.Width - r.X
				return r, true
			}
			return r, false
		} else {
			if r2.X+r2.Width >= r.X {
				r2.Width = r.X + r.Width - r2.X
				return r2, true
			}
			return r, false
		}
	}

	return r, false
}

type Sphere struct {
	Center Vector3
	Radius float64
}

func NewSphere(center Vector3, radius float64) Sphere {
	return Sphere{
		Center: center,
		Radius: radius,
	}
}

func (s Sphere) IntersectsSphere(s2 Sphere) bool {
	return s.Center.Distance(s2.Center) < s.Radius+s2.Radius
}

func (s Sphere) ContainsSphere(s2 Sphere) bool {
	return s.Radius >= s.Center.Distance(s2.Center)+s2.Radius
}
