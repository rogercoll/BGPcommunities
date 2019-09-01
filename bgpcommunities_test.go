package BGPcommunities


import(
	"testing"
)

func TestHello(t *testing.T){
	actualResult := Hello()
	expectedResult := "Hello"
	if actualResult != expectedResult {
				t.Fatalf("Expected %s but got %s", expectedResult, actualResult)
	}
}

