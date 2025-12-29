package util

import (
	"testing"
)

func BenchmarkRandomAlphanumeric(b *testing.B) {
	lengths := []int{8, 12, 16, 32, 64}

	for _, length := range lengths {
		b.Run(string(rune(length)), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, err := RandomAlphanumeric(length)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkRandomAlphanumeric_Parallel(b *testing.B) {
	lengths := []int{8, 12, 16, 32}

	for _, length := range lengths {
		b.Run(string(rune(length)), func(b *testing.B) {
			b.ReportAllocs()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					_, err := RandomAlphanumeric(length)
					if err != nil {
						b.Fatal(err)
					}
				}
			})
		})
	}
}

func TestRandomAlphanumeric(t *testing.T) {
	tests := []struct {
		name   string
		length int
	}{
		{"short", 8},
		{"medium", 12},
		{"long", 32},
		{"very_long", 64},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := RandomAlphanumeric(tt.length)
			if err != nil {
				t.Fatalf("RandomAlphanumeric() error = %v", err)
			}
			if len(result) != tt.length {
				t.Errorf("RandomAlphanumeric() length = %v, want %v", len(result), tt.length)
			}

			// Verify all characters are alphanumeric
			for _, char := range result {
				if !((char >= 'a' && char <= 'z') ||
					(char >= 'A' && char <= 'Z') ||
					(char >= '0' && char <= '9')) {
					t.Errorf("RandomAlphanumeric() contains invalid character: %c", char)
				}
			}
		})
	}
}

func TestRandomAlphanumeric_Uniqueness(t *testing.T) {
	const iterations = 1000
	const length = 12

	seen := make(map[string]bool)
	collisions := 0

	for range iterations {
		result, err := RandomAlphanumeric(length)
		if err != nil {
			t.Fatalf("RandomAlphanumeric() error = %v", err)
		}

		if seen[result] {
			collisions++
		}
		seen[result] = true
	}

	// With 12 chars and 62 possible characters, collision rate should be near zero
	collisionRate := float64(collisions) / float64(iterations)
	if collisionRate > 0.01 { // Allow 1% collision rate maximum
		t.Errorf("Too many collisions: %v out of %v (%.2f%%)", collisions, iterations, collisionRate*100)
	}
}
