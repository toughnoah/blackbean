package cmd

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/prashantv/gostub"
	"github.com/spf13/afero"
	"path/filepath"
)

var _ = Describe("cat resources test", func() {
	Context("test no filecompletion", func() {
		It("test noCompletions", func() {
			fs := afero.NewOsFs()
			filename := ".blackbean.yaml"
			path := "/tmp"
			file := filepath.Join(path, filename)
			_, createErr := fs.Create(file)
			gostub.Stub(&cfgFile, "/tmp/.blackbean.yaml")
			Expect(createErr).To(BeNil())
			defer func() {
				_ = fs.Remove(file)
			}()
			InitConfig()
		})
	})
})
