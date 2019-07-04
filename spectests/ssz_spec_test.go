package autogenerated

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/bazelbuild/rules_go/go/tools/bazel"
	"github.com/ghodss/yaml"
	"github.com/prysmaticlabs/go-ssz"
)

// sszComparisonConfig is used to specify the value to marshal, unmarshal into,
// as well as the expected results from the spec test YAML files.
type sszComparisonConfig struct {
	val                 interface{}
	unmarshalTarget     interface{}
	expected            []byte
	expectedRoot        []byte
	expectedSigningRoot []byte
}

func TestYamlStateRoundTrip(t *testing.T) {
	s := &SszBenchmarkState{}
	populateStructFromYaml(t, "./yaml/ssz_single_state.yaml", s)
	compareSSZEncoding(t, &sszComparisonConfig{
		val:             s.Value,
		unmarshalTarget: new(MainnetBeaconState),
		expected:        s.Serialized,
		expectedRoot:    s.Root,
	})
}

func TestYamlBlockRoundTrip(t *testing.T) {
	s := &SszBenchmarkBlock{}
	populateStructFromYaml(t, "./yaml/ssz_single_block.yaml", s)
	compareSSZEncoding(t, &sszComparisonConfig{
		val:                 s.Value,
		unmarshalTarget:     new(MainnetBlock),
		expected:            s.Serialized,
		expectedRoot:        s.Root,
		expectedSigningRoot: s.SigningRoot,
	})
}

func TestYamlGenericSpecTests(t *testing.T) {
	topPath := "/eth2_spec_tests/tests/ssz_generic/uint/"
	yamlFileNames := []string{
		"uint_bounds.yaml",
		"uint_wrong_length.yaml",
		"uint_random.yaml",
	}
	for _, f := range yamlFileNames {
		fullName := path.Join(topPath, f)
		fPath, err := bazel.Runfile(fullName)
		if err != nil {
			t.Fatal(err)
		}
		yamlFile, err := ioutil.ReadFile(fPath)
		if err != nil {
			t.Fatal(err)
		}
		s := &SszGenericTest{}
		if err := yaml.Unmarshal(yamlFile, s); err != nil {
			t.Fatalf("Failed to unmarshal: %v", err)
		}
		for _, testCase := range s.TestCases {
			switch testCase.Type {
			case "uint8":
				if testCase.Valid {
					num, _ := strconv.ParseUint(testCase.Value, 10, 8)
					encoded, err := ssz.Marshal(uint8(num))
					if err != nil {
						t.Fatal(err)
					}
					if !bytes.Equal(encoded, testCase.Ssz) {
						t.Errorf("Expected %v, received %v", testCase.Ssz, encoded)
					}
				} else {
					if _, err := strconv.ParseUint(testCase.Value, 10, 8); err == nil {
						t.Error("Expected error, received nil")
					}
				}
			case "uint16":
				if testCase.Valid {
					num, _ := strconv.ParseUint(testCase.Value, 10, 16)
					encoded, err := ssz.Marshal(uint16(num))
					if err != nil {
						t.Fatal(err)
					}
					if !bytes.Equal(encoded, testCase.Ssz) {
						t.Errorf("Expected %v, received %v", testCase.Ssz, encoded)
					}
				} else {
					if _, err := strconv.ParseUint(testCase.Value, 10, 16); err == nil {
						t.Error("Expected error, received nil")
					}
				}
			case "uint32":
				if testCase.Valid {
					num, _ := strconv.ParseUint(testCase.Value, 10, 32)
					encoded, err := ssz.Marshal(uint32(num))
					if err != nil {
						t.Fatal(err)
					}
					if !bytes.Equal(encoded, testCase.Ssz) {
						t.Errorf("Expected %v, received %v", testCase.Ssz, encoded)
					}
				} else {
					if _, err := strconv.ParseUint(testCase.Value, 10, 32); err == nil {
						t.Error("Expected error, received nil")
					}
				}
			case "uint64":
				if testCase.Valid {
					num, _ := strconv.ParseUint(testCase.Value, 10, 64)
					encoded, err := ssz.Marshal(num)
					if err != nil {
						t.Fatal(err)
					}
					if !bytes.Equal(encoded, testCase.Ssz) {
						t.Errorf("Expected %v, received %v", testCase.Ssz, encoded)
					}
				} else {
					if _, err := strconv.ParseUint(testCase.Value, 10, 64); err == nil {
						t.Error("Expected error, received nil")
					}
				}
			}
		}
	}
}

func TestYamlStaticSpecTests(t *testing.T) {
	topPath := "/eth2_spec_tests/tests/ssz_static/core/"
	yamlFileNames := []string{
		// "ssz_mainnet_random.yaml",
		// "ssz_minimal_lengthy.yaml",
		// "ssz_minimal_max.yaml",
		// "ssz_minimal_nil.yaml",
		"ssz_minimal_one.yaml",
		// "ssz_minimal_random.yaml",
		// "ssz_minimal_random_chaos.yaml",
		// "ssz_minimal_zero.yaml",
	}
	for _, f := range yamlFileNames {
		fullName := path.Join(topPath, f)
		fPath, err := bazel.Runfile(fullName)
		if err != nil {
			t.Fatal(err)
		}
		yamlFile, err := ioutil.ReadFile(fPath)
		if err != nil {
			t.Fatal(err)
		}
		t.Run(f, func(tt *testing.T) {
			if strings.Contains(fullName, "minimal") {
				s := &SszMinimalTest{}
				if err := yaml.Unmarshal(yamlFile, s); err != nil {
					t.Fatalf("Failed to unmarshal: %v", err)
				}
				runMinimalSpecTestCases(tt, s)
			} else if strings.Contains(fullName, "mainnet") {
				s := &SszMainnetTest{}
				if err := yaml.Unmarshal(yamlFile, s); err != nil {
					t.Fatalf("Failed to unmarshal: %v", err)
				}
				runMainnetSpecTestCases(tt, s)
			}
		})
	}
}

func runMinimalSpecTestCases(t *testing.T, s *SszMinimalTest) {
	for _, testCase := range s.TestCases {
		if !isEmpty(testCase.Attestation.Value) {
			compareSSZEncoding(t, &sszComparisonConfig{
				val:                 testCase.Attestation.Value,
				unmarshalTarget:     new(MinimalAttestation),
				expected:            testCase.Attestation.Serialized,
				expectedRoot:        testCase.Attestation.Root,
				expectedSigningRoot: testCase.Attestation.SigningRoot,
			})
		}
		if !isEmpty(testCase.AttestationData.Value) {
			compareSSZEncoding(t, &sszComparisonConfig{
				val:             testCase.AttestationData.Value,
				unmarshalTarget: new(MinimalAttestationData),
				expected:        testCase.AttestationData.Serialized,
				expectedRoot:    testCase.AttestationData.Root,
			})
		}
		if !isEmpty(testCase.AttestationDataAndCustodyBit.Value) {
			compareSSZEncoding(t, &sszComparisonConfig{
				val:             testCase.AttestationDataAndCustodyBit.Value,
				unmarshalTarget: new(MinimalAttestationAndCustodyBit),
				expected:        testCase.AttestationDataAndCustodyBit.Serialized,
				expectedRoot:    testCase.AttestationDataAndCustodyBit.Root,
			})
		}
		// if !isEmpty(testCase.AttesterSlashing.Value) {
		// 	compareSSZEncoding(t, &sszComparisonConfig{
		// 		val:             testCase.AttesterSlashing.Value,
		// 		unmarshalTarget: new(MinimalAttesterSlashing),
		// 		expected:        testCase.AttesterSlashing.Serialized,
		// 		expectedRoot:    testCase.AttesterSlashing.Root,
		// 	})
		// }
		// if !isEmpty(testCase.BeaconBlock.Value) {
		// 	compareSSZEncoding(t, &sszComparisonConfig{
		// 		val:                 testCase.BeaconBlock.Value,
		// 		unmarshalTarget:     new(MinimalBlock),
		// 		expected:            testCase.BeaconBlock.Serialized,
		// 		expectedRoot:        testCase.BeaconBlock.Root,
		// 		expectedSigningRoot: testCase.BeaconBlock.SigningRoot,
		// 	})
		// }
		// if !isEmpty(testCase.BeaconBlockBody.Value) {
		// 	compareSSZEncoding(t, &sszComparisonConfig{
		// 		val:             testCase.BeaconBlockBody.Value,
		// 		unmarshalTarget: new(MinimalBlockBody),
		// 		expected:        testCase.BeaconBlockBody.Serialized,
		// 		expectedRoot:    testCase.BeaconBlockBody.Root,
		// 	})
		// }
		if !isEmpty(testCase.BeaconBlockHeader.Value) {
			compareSSZEncoding(t, &sszComparisonConfig{
				val:                 testCase.BeaconBlockHeader.Value,
				unmarshalTarget:     new(MinimalBlockHeader),
				expected:            testCase.BeaconBlockHeader.Serialized,
				expectedRoot:        testCase.BeaconBlockHeader.Root,
				expectedSigningRoot: testCase.BeaconBlockHeader.SigningRoot,
			})
		}
		// if !isEmpty(testCase.BeaconState.Value) {
		// 	compareSSZEncoding(t, &sszComparisonConfig{
		// 		val:             testCase.BeaconState.Value,
		// 		unmarshalTarget: new(MinimalBeaconState),
		// 		expected:        testCase.BeaconState.Serialized,
		// 		expectedRoot:    testCase.BeaconState.Root,
		// 	})
		// }
		if !isEmpty(testCase.Checkpoint.Value) {
			compareSSZEncoding(t, &sszComparisonConfig{
				val:             testCase.Checkpoint.Value,
				unmarshalTarget: new(MinimalCheckpoint),
				expected:        testCase.Checkpoint.Serialized,
				expectedRoot:    testCase.Checkpoint.Root,
			})
		}
		// if !isEmpty(testCase.CompactCommittee.Value) {
		// 	compareSSZEncoding(t, &sszComparisonConfig{
		// 		val:             testCase.CompactCommittee.Value,
		// 		unmarshalTarget: new(MinimalCompactCommittee),
		// 		expected:        testCase.CompactCommittee.Serialized,
		// 		expectedRoot:    testCase.CompactCommittee.Root,
		// 	})
		// }
		if !isEmpty(testCase.Crosslink.Value) {
			compareSSZEncoding(t, &sszComparisonConfig{
				val:             testCase.Crosslink.Value,
				unmarshalTarget: new(MinimalCrosslink),
				expected:        testCase.Crosslink.Serialized,
				expectedRoot:    testCase.Crosslink.Root,
			})
		}
		if !isEmpty(testCase.Deposit.Value) {
			compareSSZEncoding(t, &sszComparisonConfig{
				val:             testCase.Deposit.Value,
				unmarshalTarget: new(MinimalDeposit),
				expected:        testCase.Deposit.Serialized,
				expectedRoot:    testCase.Deposit.Root,
			})
		}
		if !isEmpty(testCase.DepositData.Value) {
			compareSSZEncoding(t, &sszComparisonConfig{
				val:             testCase.DepositData.Value,
				unmarshalTarget: new(MinimalDepositData),
				expected:        testCase.DepositData.Serialized,
				expectedRoot:    testCase.DepositData.Root,
			})
		}
		if !isEmpty(testCase.Eth1Data.Value) {
			compareSSZEncoding(t, &sszComparisonConfig{
				val:             testCase.Eth1Data.Value,
				unmarshalTarget: new(MinimalEth1Data),
				expected:        testCase.Eth1Data.Serialized,
				expectedRoot:    testCase.Eth1Data.Root,
			})
		}
		if !isEmpty(testCase.Fork.Value) {
			compareSSZEncoding(t, &sszComparisonConfig{
				val:             testCase.Fork.Value,
				unmarshalTarget: new(MinimalFork),
				expected:        testCase.Fork.Serialized,
				expectedRoot:    testCase.Fork.Root,
			})
		}
		if !isEmpty(testCase.HistoricalBatch.Value) {
			compareSSZEncoding(t, &sszComparisonConfig{
				val:             testCase.HistoricalBatch.Value,
				unmarshalTarget: new(MinimalHistoricalBatch),
				expected:        testCase.HistoricalBatch.Serialized,
				expectedRoot:    testCase.HistoricalBatch.Root,
			})
		}
		// if !isEmpty(testCase.IndexedAttestation.Value) {
		// 	compareSSZEncoding(t, &sszComparisonConfig{
		// 		val:                 testCase.IndexedAttestation.Value,
		// 		unmarshalTarget:     new(MinimalIndexedAttestation),
		// 		expected:            testCase.IndexedAttestation.Serialized,
		// 		expectedRoot:        testCase.IndexedAttestation.Root,
		// 		expectedSigningRoot: testCase.IndexedAttestation.SigningRoot,
		// 	})
		// }
		// if !isEmpty(testCase.PendingAttestation.Value) {
		// 	compareSSZEncoding(t, &sszComparisonConfig{
		// 		val:             testCase.PendingAttestation.Value,
		// 		unmarshalTarget: new(MinimalPendingAttestation),
		// 		expected:        testCase.PendingAttestation.Serialized,
		// 		expectedRoot:    testCase.PendingAttestation.Root,
		// 	})
		// }
		if !isEmpty(testCase.ProposerSlashing.Value) {
			compareSSZEncoding(t, &sszComparisonConfig{
				val:             testCase.ProposerSlashing.Value,
				unmarshalTarget: new(MinimalProposerSlashing),
				expected:        testCase.ProposerSlashing.Serialized,
				expectedRoot:    testCase.ProposerSlashing.Root,
			})
		}
		if !isEmpty(testCase.Transfer.Value) {
			compareSSZEncoding(t, &sszComparisonConfig{
				val:                 testCase.Transfer.Value,
				unmarshalTarget:     new(MinimalTransfer),
				expected:            testCase.Transfer.Serialized,
				expectedRoot:        testCase.Transfer.Root,
				expectedSigningRoot: testCase.Transfer.SigningRoot,
			})
		}
		if !isEmpty(testCase.Validator.Value) {
			compareSSZEncoding(t, &sszComparisonConfig{
				val:             testCase.Validator.Value,
				unmarshalTarget: new(MinimalValidator),
				expected:        testCase.Validator.Serialized,
				expectedRoot:    testCase.Validator.Root,
			})
		}
		if !isEmpty(testCase.VoluntaryExit.Value) {
			compareSSZEncoding(t, &sszComparisonConfig{
				val:                 testCase.VoluntaryExit.Value,
				unmarshalTarget:     new(MinimalVoluntaryExit),
				expected:            testCase.VoluntaryExit.Serialized,
				expectedRoot:        testCase.VoluntaryExit.Root,
				expectedSigningRoot: testCase.VoluntaryExit.SigningRoot,
			})
		}
	}
}

func runMainnetSpecTestCases(t *testing.T, s *SszMainnetTest) {
	for _, testCase := range s.TestCases {
		// if !isEmpty(testCase.Attestation.Value) {
		// 	compareSSZEncoding(t, &sszComparisonConfig{
		// 		val:                 testCase.Attestation.Value,
		// 		unmarshalTarget:     new(MainnetAttestation),
		// 		expected:            testCase.Attestation.Serialized,
		// 		expectedRoot:        testCase.Attestation.Root,
		// 		expectedSigningRoot: testCase.Attestation.SigningRoot,
		// 	})
		// }
		if !isEmpty(testCase.AttestationData.Value) {
			compareSSZEncoding(t, &sszComparisonConfig{
				val:             testCase.AttestationData.Value,
				unmarshalTarget: new(MainnetAttestationData),
				expected:        testCase.AttestationData.Serialized,
				expectedRoot:    testCase.AttestationData.Root,
			})
		}
		if !isEmpty(testCase.AttestationDataAndCustodyBit.Value) {
			compareSSZEncoding(t, &sszComparisonConfig{
				val:             testCase.AttestationDataAndCustodyBit.Value,
				unmarshalTarget: new(MainnetAttestationAndCustodyBit),
				expected:        testCase.AttestationDataAndCustodyBit.Serialized,
				expectedRoot:    testCase.AttestationDataAndCustodyBit.Root,
			})
		}
		// if !isEmpty(testCase.AttesterSlashing.Value) {
		// 	compareSSZEncoding(t, &sszComparisonConfig{
		// 		val:             testCase.AttesterSlashing.Value,
		// 		unmarshalTarget: new(MainnetAttesterSlashing),
		// 		expected:        testCase.AttesterSlashing.Serialized,
		// 		expectedRoot:    testCase.AttesterSlashing.Root,
		// 	})
		// }
		// if !isEmpty(testCase.BeaconBlock.Value) {
		// 	compareSSZEncoding(t, &sszComparisonConfig{
		// 		val:                 testCase.BeaconBlock.Value,
		// 		unmarshalTarget:     new(MainnetBlock),
		// 		expected:            testCase.BeaconBlock.Serialized,
		// 		expectedRoot:        testCase.BeaconBlock.Root,
		// 		expectedSigningRoot: testCase.BeaconBlock.SigningRoot,
		// 	})
		// }
		// if !isEmpty(testCase.BeaconBlockBody.Value) {
		// 	compareSSZEncoding(t, &sszComparisonConfig{
		// 		val:             testCase.BeaconBlockBody.Value,
		// 		unmarshalTarget: new(MainnetBlockBody),
		// 		expected:        testCase.BeaconBlockBody.Serialized,
		// 		expectedRoot:    testCase.BeaconBlockBody.Root,
		// 	})
		// }
		if !isEmpty(testCase.BeaconBlockHeader.Value) {
			compareSSZEncoding(t, &sszComparisonConfig{
				val:                 testCase.BeaconBlockHeader.Value,
				unmarshalTarget:     new(MainnetBlockHeader),
				expected:            testCase.BeaconBlockHeader.Serialized,
				expectedRoot:        testCase.BeaconBlockHeader.Root,
				expectedSigningRoot: testCase.BeaconBlockHeader.SigningRoot,
			})
		}
		// if !isEmpty(testCase.BeaconState.Value) {
		// 	compareSSZEncoding(t, &sszComparisonConfig{
		// 		val:             testCase.BeaconState.Value,
		// 		unmarshalTarget: new(MainnetBeaconState),
		// 		expected:        testCase.BeaconState.Serialized,
		// 		expectedRoot:    testCase.BeaconState.Root,
		// 	})
		// }
		if !isEmpty(testCase.Checkpoint.Value) {
			compareSSZEncoding(t, &sszComparisonConfig{
				val:             testCase.Checkpoint.Value,
				unmarshalTarget: new(MinimalCheckpoint),
				expected:        testCase.Checkpoint.Serialized,
				expectedRoot:    testCase.Checkpoint.Root,
			})
		}
		// if !isEmpty(testCase.CompactCommittee.Value) {
		// 	compareSSZEncoding(t, &sszComparisonConfig{
		// 		val:             testCase.CompactCommittee.Value,
		// 		unmarshalTarget: new(MinimalCompactCommittee),
		// 		expected:        testCase.CompactCommittee.Serialized,
		// 		expectedRoot:    testCase.CompactCommittee.Root,
		// 	})
		// }
		if !isEmpty(testCase.Crosslink.Value) {
			compareSSZEncoding(t, &sszComparisonConfig{
				val:             testCase.Crosslink.Value,
				unmarshalTarget: new(MainnetCrosslink),
				expected:        testCase.Crosslink.Serialized,
				expectedRoot:    testCase.Crosslink.Root,
			})
		}
		if !isEmpty(testCase.Deposit.Value) {
			compareSSZEncoding(t, &sszComparisonConfig{
				val:             testCase.Deposit.Value,
				unmarshalTarget: new(MainnetDeposit),
				expected:        testCase.Deposit.Serialized,
				expectedRoot:    testCase.Deposit.Root,
			})
		}
		if !isEmpty(testCase.DepositData.Value) {
			compareSSZEncoding(t, &sszComparisonConfig{
				val:             testCase.DepositData.Value,
				unmarshalTarget: new(MainnetDepositData),
				expected:        testCase.DepositData.Serialized,
				expectedRoot:    testCase.DepositData.Root,
			})
		}
		if !isEmpty(testCase.Eth1Data.Value) {
			compareSSZEncoding(t, &sszComparisonConfig{
				val:             testCase.Eth1Data.Value,
				unmarshalTarget: new(MainnetEth1Data),
				expected:        testCase.Eth1Data.Serialized,
				expectedRoot:    testCase.Eth1Data.Root,
			})
		}
		if !isEmpty(testCase.Fork.Value) {
			compareSSZEncoding(t, &sszComparisonConfig{
				val:             testCase.Fork.Value,
				unmarshalTarget: new(MainnetFork),
				expected:        testCase.Fork.Serialized,
				expectedRoot:    testCase.Fork.Root,
			})
		}
		if !isEmpty(testCase.HistoricalBatch.Value) {
			compareSSZEncoding(t, &sszComparisonConfig{
				val:             testCase.HistoricalBatch.Value,
				unmarshalTarget: new(MainnetHistoricalBatch),
				expected:        testCase.HistoricalBatch.Serialized,
				expectedRoot:    testCase.HistoricalBatch.Root,
			})
		}
		// if !isEmpty(testCase.IndexedAttestation.Value) {
		// 	compareSSZEncoding(t, &sszComparisonConfig{
		// 		val:                 testCase.IndexedAttestation.Value,
		// 		unmarshalTarget:     new(MainnetIndexedAttestation),
		// 		expected:            testCase.IndexedAttestation.Serialized,
		// 		expectedRoot:        testCase.IndexedAttestation.Root,
		// 		expectedSigningRoot: testCase.IndexedAttestation.SigningRoot,
		// 	})
		// }
		// if !isEmpty(testCase.PendingAttestation.Value) {
		// 	compareSSZEncoding(t, &sszComparisonConfig{
		// 		val:             testCase.PendingAttestation.Value,
		// 		unmarshalTarget: new(MainnetPendingAttestation),
		// 		expected:        testCase.PendingAttestation.Serialized,
		// 		expectedRoot:    testCase.PendingAttestation.Root,
		// 	})
		// }
		if !isEmpty(testCase.ProposerSlashing.Value) {
			compareSSZEncoding(t, &sszComparisonConfig{
				val:             testCase.ProposerSlashing.Value,
				unmarshalTarget: new(MainnetProposerSlashing),
				expected:        testCase.ProposerSlashing.Serialized,
				expectedRoot:    testCase.ProposerSlashing.Root,
			})
		}
		if !isEmpty(testCase.Transfer.Value) {
			compareSSZEncoding(t, &sszComparisonConfig{
				val:                 testCase.Transfer.Value,
				unmarshalTarget:     new(MainnetTransfer),
				expected:            testCase.Transfer.Serialized,
				expectedRoot:        testCase.Transfer.Root,
				expectedSigningRoot: testCase.Transfer.SigningRoot,
			})
		}
		if !isEmpty(testCase.Validator.Value) {
			compareSSZEncoding(t, &sszComparisonConfig{
				val:             testCase.Validator.Value,
				unmarshalTarget: new(MainnetValidator),
				expected:        testCase.Validator.Serialized,
				expectedRoot:    testCase.Validator.Root,
			})
		}
		if !isEmpty(testCase.VoluntaryExit.Value) {
			compareSSZEncoding(t, &sszComparisonConfig{
				val:                 testCase.VoluntaryExit.Value,
				unmarshalTarget:     new(MainnetVoluntaryExit),
				expected:            testCase.VoluntaryExit.Serialized,
				expectedRoot:        testCase.VoluntaryExit.Root,
				expectedSigningRoot: testCase.VoluntaryExit.SigningRoot,
			})
		}
	}
}

func isEmpty(item interface{}) bool {
	val := reflect.ValueOf(item)
	for i := 0; i < val.NumField(); i++ {
		if !reflect.DeepEqual(val.Field(i).Interface(), reflect.Zero(val.Field(i).Type()).Interface()) {
			return false
		}
	}
	return true
}

func compareSSZEncoding(t *testing.T, cfg *sszComparisonConfig) {
	encoded, err := ssz.Marshal(cfg.val)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(encoded, cfg.expected) {
		fmt.Println(encoded)
		fmt.Println(cfg.expected)
		t.Fatal("Failed to encode")
	}
	if err := ssz.Unmarshal(encoded, cfg.unmarshalTarget); err != nil {
		t.Fatal(err)
	}
	concreteValue := reflect.ValueOf(cfg.unmarshalTarget).Elem().Interface()
	if !ssz.DeepEqual(concreteValue, cfg.val) {
		t.Error("Unmarshaled encoding did not match original value")
	}
	root, err := ssz.HashTreeRoot(cfg.val)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(root[:], cfg.expectedRoot) {
		t.Fatalf("Expected hash tree root %#x, received %#x", cfg.expectedRoot, root[:])
	}
	if cfg.expectedSigningRoot != nil {
		signingRoot, err := ssz.SigningRoot(cfg.val)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(signingRoot[:], cfg.expectedSigningRoot) {
			t.Errorf("Expected signing root %#x, received %#x", cfg.expectedSigningRoot, signingRoot)
		}
	}
}
