package tests

import(
	"github.com/stretchr/testify/mock"
	"slices"
	"aitring/tests/mock_repos"

)

 var MockAudioRepo  = new(mock_repos.MockAudioRepo)
  var MockPieRepo  = new(mock_repos.MockPipeRepo)


// Mock for Repos

func InitializeMockFunctions(functions []func() *mock.Call) []*mock.Call {
	var mockCalls = make([]*mock.Call, 0)
	for _, function := range functions {
		mockCall := function()
		mockCalls = append(mockCalls, mockCall)
	}
	return mockCalls
}

func UnsetMockFunctionCalls(mockCalls []*mock.Call, skipMockCallUnsetMethods []string) {
	mockCallsList := make([]string, 0)
	for _, mockCall := range mockCalls {
		if !slices.Contains(skipMockCallUnsetMethods, mockCall.Method) && !slices.Contains(mockCallsList, mockCall.Method) {
			mockCallsList = append(mockCallsList, mockCall.Method)
			mockCall.Unset()
		}
	}
}

func MockCurrentMillis() int {
	// Return a fixed value for testing purposes.
	return 1624158000000 // Replace this value with the desired timestamp
}
func MockCuidGenerator() string {
	return "cliy8jvtw00233b6n1ly531au"
}

func MockMD5HashGenerator(plainText string) string {
	plainText = "ed076287532e86365e841e92bfc50d8c"
	return plainText
}
