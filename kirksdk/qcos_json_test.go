package kirksdk

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestServiceSpecMarshal(t *testing.T) {
	spec := ServiceSpec{
		AutoRestart: "always",
		Command:     []string{},
		EntryPoint:  nil,
		Envs:        []string{"a=a1", "b=b1"},
		// Host:	nil
		Image:         "Image",
		LogCollectors: []LogCollectorSpec{},
		StopGraceSec:  1,
		WorkDir:       "~/",
		UnitType:      "unittype",
		GpuUUIDs:      []string{},
	}

	testMarshal(
		t,
		spec,
		`{"command":[],"envs":["a=a1","b=b1"],"logCollectors":[],"gpuUUIDs":[],"autoRestart":"always","image":"Image","stopGraceSec":1,"workDir":"~/","unitType":"unittype"}`,
	)
}

func TestJobTaskSpecMarshal(t *testing.T) {
	spec := JobTaskSpec{
		Image:         "testimage",
		Command:       nil,
		EntryPoint:    []string{},
		Envs:          []string{},
		Hosts:         []string{},
		LogCollectors: []LogCollectorSpec{},
		WorkDir:       "testworkdir",
		InstanceNum:   2,
		UnitType:      "testunittype",
	}

	testMarshal(
		t,
		spec,
		`{"entryPoint":[],"envs":[],"hosts":[],"logCollectors":[],"image":"testimage","workDir":"testworkdir","instanceNum":2,"unitType":"testunittype"}`,
	)
}

func TestUpdateJobArgsMarshal(t *testing.T) {
	args := UpdateJobArgs{
		Spec:     map[string]JobTaskSpec{},
		Metadata: []string{},
		RunAt:    "testrunat",
		Timeout:  3000,
		Mode:     "testmode",
	}

	testMarshal(
		t,
		args,
		`{"metadata":[],"runAt":"testrunat","timeout":3000,"mode":"testmode"}`,
	)
}

func TestJobTaskSpecExMarshal(t *testing.T) {
	specEx := JobTaskSpecEx{
		WorkDir: "",
		LogCollectors: []LogCollectorSpec{
			LogCollectorSpec{
				Directory: "testdir1",
				Patterns:  []string{"p1", "p2"},
			},
		},
		Command:     []string{},
		EntryPoint:  []string{},
		Envs:        []string{},
		Hosts:       []string{},
		UnitType:    "testunittype",
		Deps:        []string{},
		InstanceNum: 3,
	}

	testMarshal(
		t,
		&specEx,
		`{"logCollectors":[{"directory":"testdir1","patterns":["p1","p2"]}],"command":[],"entryPoint":[],"envs":[],"hosts":[],"deps":[],"unitType":"testunittype","instanceNum":3}`,
	)
}

func testMarshal(t *testing.T, input interface{}, expected string) {
	bytes, err := json.Marshal(input)
	assert.Nil(t, err)

	actual := string(bytes)

	fmt.Println(actual)

	assert.Equal(t, expected, actual)
}
