package rules

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/uber-go/tally"
	"github.com/uber/aresdb/client"
	memCom "github.com/uber/aresdb/memstore/common"
	"github.com/uber/aresdb/subscriber/common/tools"
	"github.com/uber/aresdb/subscriber/config"
	"github.com/uber/aresdb/utils"
	"go.uber.org/zap"
)

var _ = Describe("job_config", func() {
	serviceConfig := config.ServiceConfig{
		Environment: utils.EnvironmentContext{
			Deployment:         "test",
			RuntimeEnvironment: "test",
			Zone:               "local",
		},
		Logger: zap.NewNop(),
		Scope:  tally.NoopScope,
	}
	serviceConfig.ActiveJobs = []string{"job1"}
	serviceConfig.ActiveAresClusters = map[string]client.ConnectorConfig{
		"dev01": client.ConnectorConfig{Address: "localhost:8888"},
	}
	p := Params{
		ServiceConfig: serviceConfig,
	}

	It("NewJobConfigs", func() {
		serviceConfig.ActiveJobs = []string{"dispatch"}
		rootPath := tools.GetModulePath("")
		os.Chdir(rootPath)
		rst, err := NewJobConfigs(p)
		Ω(rst).ShouldNot(BeNil())
		Ω(err).Should(BeNil())
		Ω(rst.JobConfigs["job1"]).ShouldNot(BeNil())
		Ω(rst.JobConfigs["job1"]["dev01"]).ShouldNot(BeNil())

		dst := rst.JobConfigs["job1"]["dev01"].GetDestinations()
		transformation := rst.JobConfigs["job1"]["dev01"].GetTranformations()
		primaryKey := rst.JobConfigs["job1"]["dev01"].GetPrimaryKeys()
		Ω(dst).ShouldNot(BeNil())
		Ω(transformation).ShouldNot(BeNil())
		Ω(primaryKey).ShouldNot(BeNil())
	})
	It("parseUpdateMode", func() {
		mode := parseUpdateMode("overwrite_notnull")
		Ω(mode).Should(Equal(memCom.UpdateOverwriteNotNull))

		mode = parseUpdateMode("overwrite_force")
		Ω(mode).Should(Equal(memCom.UpdateForceOverwrite))

		mode = parseUpdateMode("min")
		Ω(mode).Should(Equal(memCom.UpdateWithMin))

		mode = parseUpdateMode("max")
		Ω(mode).Should(Equal(memCom.UpdateWithMax))

		mode = parseUpdateMode("")
		Ω(mode).Should(Equal(memCom.UpdateOverwriteNotNull))
	})
})
