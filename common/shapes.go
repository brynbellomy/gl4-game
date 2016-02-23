package common

func Rect(size Size) []float32 {
	halfWidth := size.Width() / 2
	halfHeight := size.Height() / 2
	return []float32{
		//  X, Y, Z, U, V

		// Front
		-halfWidth, -halfHeight, 0.0, 1.0, 0.0,
		halfWidth, -halfHeight, 0.0, 0.0, 0.0,
		-halfWidth, halfHeight, 0.0, 1.0, 1.0,

		halfWidth, -halfHeight, 0.0, 0.0, 0.0,
		halfWidth, halfHeight, 0.0, 0.0, 1.0,
		-halfWidth, halfHeight, 0.0, 1.0, 1.0,
	}
}
