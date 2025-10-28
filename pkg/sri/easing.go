package sri

// EasingFunc is a function that maps linear progress [0,1] to eased progress [0,1].
type EasingFunc func(t float64) float64

// Linear easing (no easing).
func Linear(t float64) float64 {
	return t
}

// EaseIn easing (quadratic).
func EaseIn(t float64) float64 {
	return t * t
}

// EaseOut easing (quadratic).
func EaseOut(t float64) float64 {
	return t * (2.0 - t)
}

// EaseInOut easing (quadratic).
func EaseInOut(t float64) float64 {
	if t < 0.5 {
		return 2.0 * t * t
	}
	return -1.0 + (4.0-2.0*t)*t
}

// EaseInCubic easing (cubic).
func EaseInCubic(t float64) float64 {
	return t * t * t
}

// EaseOutCubic easing (cubic).
func EaseOutCubic(t float64) float64 {
	t--
	return t*t*t + 1.0
}

// EaseInOutCubic easing (cubic).
func EaseInOutCubic(t float64) float64 {
	if t < 0.5 {
		return 4.0 * t * t * t
	}
	t = t*2.0 - 2.0
	return (t*t*t + 2.0) / 2.0
}

// EaseInQuart easing (quartic).
func EaseInQuart(t float64) float64 {
	return t * t * t * t
}

// EaseOutQuart easing (quartic).
func EaseOutQuart(t float64) float64 {
	t--
	return 1.0 - t*t*t*t
}

// EaseInOutQuart easing (quartic).
func EaseInOutQuart(t float64) float64 {
	if t < 0.5 {
		return 8.0 * t * t * t * t
	}
	t--
	return 1.0 - 8.0*t*t*t*t
}

// EaseInQuint easing (quintic).
func EaseInQuint(t float64) float64 {
	return t * t * t * t * t
}

// EaseOutQuint easing (quintic).
func EaseOutQuint(t float64) float64 {
	t--
	return t*t*t*t*t + 1.0
}

// EaseInOutQuint easing (quintic).
func EaseInOutQuint(t float64) float64 {
	if t < 0.5 {
		return 16.0 * t * t * t * t * t
	}
	t = t*2.0 - 2.0
	return (t*t*t*t*t + 2.0) / 2.0
}

// EaseInExpo easing (exponential).
func EaseInExpo(t float64) float64 {
	if t == 0.0 {
		return 0.0
	}
	return pow(2.0, 10.0*(t-1.0))
}

// EaseOutExpo easing (exponential).
func EaseOutExpo(t float64) float64 {
	if t == 1.0 {
		return 1.0
	}
	return 1.0 - pow(2.0, -10.0*t)
}

// EaseInOutExpo easing (exponential).
func EaseInOutExpo(t float64) float64 {
	if t == 0.0 || t == 1.0 {
		return t
	}
	if t < 0.5 {
		return pow(2.0, 20.0*t-10.0) / 2.0
	}
	return (2.0 - pow(2.0, -20.0*t+10.0)) / 2.0
}

// pow approximates power function (simple implementation).
func pow(base, exp float64) float64 {
	if exp == 0.0 {
		return 1.0
	}
	result := 1.0
	for i := 0; i < int(exp); i++ {
		result *= base
	}
	return result
}

// Bounce easing (bouncing effect).
func Bounce(t float64) float64 {
	if t < 1.0/2.75 {
		return 7.5625 * t * t
	} else if t < 2.0/2.75 {
		t -= 1.5 / 2.75
		return 7.5625*t*t + 0.75
	} else if t < 2.5/2.75 {
		t -= 2.25 / 2.75
		return 7.5625*t*t + 0.9375
	}
	t -= 2.625 / 2.75
	return 7.5625*t*t + 0.984375
}

// Elastic easing (elastic effect).
func Elastic(t float64) float64 {
	if t == 0.0 || t == 1.0 {
		return t
	}
	p := 0.3
	s := p / 4.0
	t = t - 1.0
	return -(pow(2.0, 10.0*t) * sinApprox((t-s)*(2.0*3.14159)/p))
}
