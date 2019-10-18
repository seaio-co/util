package gen
import (
	"fmt"
	"testing"
)

//output the seriserialized and eserialized test data
func TestGencode(t *testing.T) {
	p:= Block {
		ParentHash: "parentHash",
		Hash: "currentHash",
		Number:100,
		Timestamp:512,
	}
	buf, _ := p.Marshal(nil)
	fmt.Println("serialized data:",buf)
	s1:=Block{}
	s1.Unmarshal(buf)
	fmt.Println("Deserialized data:" ,s1)
}

//Output the Size of the test data
func TestGencodeSize(t *testing.T) {
	p:= Block {
		ParentHash: "parentHash",
		Hash: "currentHash",
		Number:100,
		Timestamp:512,
	}
	buf, _ := p.Marshal(nil)
	fmt.Printf("Gencode encoded size: %v\n", len(buf))
}


//Test the proformance of the test data
func BenchmarkGencodeSerialize(b *testing.B) {
	p:= Block {
		ParentHash: "parentHash",
		Hash: "currentHash",
		Number:100,
		Timestamp:512,
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.Marshal(nil)
	}
}

func BenchmarkGencodeDeserialize(b *testing.B) {
	p:= Block {
		ParentHash: "parentHash",
		Hash: "currentHash",
		Number:100,
		Timestamp:512,
	}
	buf, _ := p.Marshal(nil)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.Unmarshal(buf)
	}
}