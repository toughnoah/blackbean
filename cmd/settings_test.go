package cmd

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/toughnoah/blackbean/pkg/fake"
)

var _ = Describe("put settings test", func() {

	Context("test withFunc.", func() {
		It("test WithAllocationSettings", func() {
			testCases := []struct {
				settings *Settings
				want     *Settings
			}{
				{
					settings: NewSettings(),
					want: &Settings{Persistent: &persistent{
						Cluster: &cluster{
							Routing: &routing{
								Allocation: &allocation{
									ClusterConcurrentRebalanced:    "1",
									NodeConcurrentRecoveries:       "2",
									NodeInitialPrimariesRecoveries: "",
									Enable:                         "",
									Disk:                           nil,
								},
							},
						},
					},
					},
				},
				{
					settings: &Settings{
						Persistent: &persistent{
							Cluster: &cluster{
								Routing: &routing{
									Allocation: &allocation{
										Disk: &disk{
											Watermark: &watermark{
												Low: "87%",
											},
										},
									},
								},
							},
						},
					},
					want: &Settings{Persistent: &persistent{
						Cluster: &cluster{
							Routing: &routing{
								Allocation: &allocation{
									ClusterConcurrentRebalanced:    "1",
									NodeConcurrentRecoveries:       "2",
									NodeInitialPrimariesRecoveries: "",
									Enable:                         "",
									Disk: &disk{
										Watermark: &watermark{
											Low: "87%",
										},
									},
								},
							},
						},
					},
					},
				},
			}
			for _, tc := range testCases {
				settings := tc.settings
				settings.WithAllocationSettings("1", "2", "", "")
				Expect(settings).To(Equal(tc.want))
			}
			settings := NewSettings()
			settings.WithAllocationSettings("", "", "", "")
			Expect(settings).To(Equal(settings))
		})
		It("test WithBreakerTotal", func() {
			testCases := []struct {
				settings *Settings
				want     *Settings
			}{
				{
					settings: NewSettings(),
					want: &Settings{
						Persistent: &persistent{
							Indices: &indices{
								Breaker: &breaker{
									Total: &total{
										Limit: "80%",
									},
								},
							},
						},
					},
				},
				{
					settings: &Settings{
						Persistent: &persistent{
							Indices: &indices{
								Breaker: &breaker{
									Request: &request{
										Limit: "60%",
									},
								},
							},
						},
					},
					want: &Settings{
						Persistent: &persistent{
							Indices: &indices{
								Breaker: &breaker{
									Total: &total{
										Limit: "80%",
									},
									Request: &request{
										Limit: "60%",
									},
								},
							},
						},
					},
				},
			}
			for _, tc := range testCases {
				settings := tc.settings
				settings.WithBreakerTotal("80%")
				Expect(settings).To(Equal(tc.want))
			}
			settings := NewSettings()
			settings.WithBreakerTotal("")
			Expect(settings).To(Equal(settings))
		})
		It("test WithBreakerRequest", func() {
			testCases := []struct {
				settings *Settings
				want     *Settings
			}{
				{
					settings: NewSettings(),
					want: &Settings{
						Persistent: &persistent{
							Indices: &indices{
								Breaker: &breaker{
									Request: &request{
										Limit: "60%",
									},
								},
							},
						},
					},
				},
				{
					settings: &Settings{
						Persistent: &persistent{
							Indices: &indices{
								Breaker: &breaker{
									Total: &total{
										Limit: "80%",
									},
								},
							},
						},
					},
					want: &Settings{
						Persistent: &persistent{
							Indices: &indices{
								Breaker: &breaker{
									Total: &total{
										Limit: "80%",
									},
									Request: &request{
										Limit: "60%",
									},
								},
							},
						},
					},
				},
			}
			for _, tc := range testCases {
				settings := tc.settings
				settings.WithBreakerRequest("60%")
				Expect(settings).To(Equal(tc.want))
			}
			settings := NewSettings()
			settings.WithBreakerRequest("")
			Expect(settings).To(Equal(settings))
		})
		It("test WithBreakerFielddata", func() {
			testCases := []struct {
				settings *Settings
				want     *Settings
			}{
				{
					settings: NewSettings(),
					want: &Settings{
						Persistent: &persistent{
							Indices: &indices{
								Breaker: &breaker{
									Fielddata: &fielddata{
										Limit: "60%",
									},
								},
							},
						},
					},
				},
				{
					settings: &Settings{
						Persistent: &persistent{
							Indices: &indices{
								Breaker: &breaker{
									Total: &total{
										Limit: "80%",
									},
								},
							},
						},
					},
					want: &Settings{
						Persistent: &persistent{
							Indices: &indices{
								Breaker: &breaker{
									Total: &total{
										Limit: "80%",
									},
									Fielddata: &fielddata{
										Limit: "60%",
									},
								},
							},
						},
					},
				},
			}
			for _, tc := range testCases {
				settings := tc.settings
				settings.WithBreakerFielddata("60%")
				Expect(settings).To(Equal(tc.want))
			}
			settings := NewSettings()
			settings.WithBreakerFielddata("")
			Expect(settings).To(Equal(settings))
		})
		It("test WithWatermark", func() {
			testCases := []struct {
				settings *Settings
				want     *Settings
			}{
				{
					settings: NewSettings(),
					want: &Settings{
						Persistent: &persistent{
							Cluster: &cluster{
								Routing: &routing{
									Allocation: &allocation{
										Disk: &disk{
											Watermark: &watermark{
												Low:  "80%",
												High: "85%",
											},
										},
									},
								},
							},
						},
					},
				},
				{
					settings: &Settings{
						Persistent: &persistent{
							Cluster: &cluster{
								Routing: &routing{
									Allocation: &allocation{
										Enable: "null",
									},
								},
							},
						},
					},
					want: &Settings{
						Persistent: &persistent{
							Cluster: &cluster{
								Routing: &routing{
									Allocation: &allocation{
										Disk: &disk{
											Watermark: &watermark{
												Low:  "80%",
												High: "85%",
											},
										},
										Enable: "null",
									},
								},
							},
						},
					},
				},
			}
			for _, tc := range testCases {
				settings := tc.settings
				settings.WithWatermark("85%", "80%")
				Expect(settings).To(Equal(tc.want))
			}
			settings := NewSettings()
			settings.WithWatermark("", "")
			Expect(settings).To(Equal(settings))
		})
		It("test WithRecovery", func() {
			testCases := []struct {
				settings *Settings
				want     *Settings
			}{
				{
					settings: NewSettings(),
					want: &Settings{
						Persistent: &persistent{
							Indices: &indices{
								Recovery: &recovery{
									MaxBytesPerSec: "1000",
								},
							},
						},
					},
				},
				{
					settings: &Settings{
						Persistent: &persistent{
							Indices: &indices{
								Breaker: &breaker{
									Total: &total{
										Limit: "85%",
									},
								},
							},
						},
					},
					want: &Settings{
						Persistent: &persistent{
							Indices: &indices{
								Breaker: &breaker{
									Total: &total{
										Limit: "85%",
									},
								},
								Recovery: &recovery{
									MaxBytesPerSec: "1000",
								},
							},
						},
					},
				},
			}
			for _, tc := range testCases {
				settings := tc.settings
				settings.WithRecovery("1000")
				Expect(settings).To(Equal(tc.want))
			}
			settings := NewSettings()
			settings.WithRecovery("")
			Expect(settings).To(Equal(settings))
		})
		It("test WithMaxShardsPerNode", func() {
			testCases := []struct {
				settings *Settings
				want     *Settings
			}{
				{
					settings: NewSettings(),
					want: &Settings{
						Persistent: &persistent{
							Cluster: &cluster{
								MaxShardsPerNode: "10000",
							},
						},
					},
				},
				{
					settings: &Settings{
						Persistent: &persistent{
							Cluster: &cluster{
								Routing: &routing{
									Allocation: &allocation{
										ClusterConcurrentRebalanced: "1",
									},
								},
							},
						},
					},
					want: &Settings{
						Persistent: &persistent{
							Cluster: &cluster{
								MaxShardsPerNode: "10000",
								Routing: &routing{
									Allocation: &allocation{
										ClusterConcurrentRebalanced: "1",
									},
								},
							},
						},
					},
				},
			}
			for _, tc := range testCases {
				settings := tc.settings
				settings.WithMaxShardsPerNode("10000")
				Expect(settings).To(Equal(tc.want))
			}
			settings := NewSettings()
			settings.WithMaxShardsPerNode("")
			Expect(settings).To(Equal(settings))
		})
		It("test WithMaxCompilationsRate", func() {
			settings := NewSettings()
			settings.WithMaxCompilationsRate("180/1m")
			want := &Settings{Persistent: &persistent{
				Script: &script{
					MaxCompilationsRate: "180/1m",
				},
			}}
			Expect(settings).To(Equal(want))
			settings = NewSettings()
			settings.WithMaxCompilationsRate("")
			Expect(settings).To(Equal(settings))
		})
	})

	Context("test apply settings", func() {
		It("test apply", func() {
			mockTr := &fake.MockEsResponse{
				ResponseString: `{"test":"apply"}`,
			}
			err := executeCommandForTesting("apply settings -i 2", mockTr)
			Expect(err).To(BeNil())
		})
	})
})
