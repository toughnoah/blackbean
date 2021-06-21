package cmd

import (
	. "github.com/onsi/ginkgo"
)

var _ = Describe("cat resources test", func() {
	Context("test no FileCompletion", func() {
		//It("test noCompletions", func() {
		//	home, err := homedir.Dir()
		//	stubs := gostub.New()
		//	Expect(err).To(BeNil())
		//	testCase := []struct {
		//		path       string
		//		stubsOrNot bool
		//	}{
		//		{
		//			path: home,
		//		},
		//		{
		//			path:       "/tmp",
		//			stubsOrNot: true,
		//		},
		//	}
		//	for _, tc := range testCase {
		//		fs := afero.NewOsFs()
		//		filename := ".blackbean.yaml"
		//		file := filepath.Join(tc.path, filename)
		//		_, createErr := fs.Create(file)
		//		if tc.stubsOrNot {
		//			stubs.Stub(&cfgFile, file)
		//		}
		//		Expect(createErr).To(BeNil())
		//		InitConfig()
		//		deleteErr := fs.Remove(file)
		//		Expect(deleteErr).To(BeNil())
		//		stubs.Reset()
		//	}
		//})
	})
})
