module github.com/toughnoah/blackbean

go 1.16

require (
	github.com/bouk/monkey v1.0.2
	github.com/deckarep/golang-set v1.7.1
	github.com/elastic/go-elasticsearch/v7 v7.13.1
	github.com/mattn/go-shellwords v1.0.12
	github.com/mitchellh/go-homedir v1.1.0
	github.com/onsi/ginkgo v1.16.4
	github.com/onsi/gomega v1.13.0
	github.com/pkg/errors v0.9.1
	github.com/prashantv/gostub v1.0.0
	github.com/spf13/afero v1.6.0
	github.com/spf13/cobra v1.1.3
	github.com/spf13/viper v1.8.0
	github.com/stretchr/testify v1.7.0
	golang.org/x/term v0.0.0-20201126162022-7de9c90e9dd1
	golang.org/x/text v0.3.6
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/apimachinery v0.21.2
	sigs.k8s.io/yaml v1.2.0
)

replace github.com/bouk/monkey v1.0.2 => bou.ke/monkey v1.0.0
