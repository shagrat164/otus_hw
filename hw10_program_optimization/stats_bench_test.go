package hw10programoptimization

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"
)

func generateLargeData(n int) string {
	var builder strings.Builder
	for i := 0; i < n; i++ {
		builder.WriteString(`{"Id":1,"Name":"Howard Mendoza","Username":"0Oliver","Email":"user`)
		builder.WriteString(fmt.Sprintf("%d", i))
		builder.WriteString(`@domain.com","Phone":"6-866-899-36-79","Password":"InAQJvsq","Address":"Blackbird Place 25"}\n`)
	}
	return builder.String()
}

// $ go test -benchmem -count 10 -bench=.
func BenchmarkGetDomainStat(b *testing.B) {
	domain := "domain.com"
	reader := bytes.NewReader([]byte(generateLargeData(100)))

	for i := 0; i < b.N; i++ {
		reader.Seek(0, io.SeekStart)

		_, err := GetDomainStat(reader, domain)
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
	}
}
