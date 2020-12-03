package fetchraw

import (
	"fmt"
	"os"
	"os/exec"

	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/JAORMX/oc-compliance/internal/common"
)

type FetchRawOptions struct {
	common.CommandContext

	OutputPath string
	HTML       bool
}

func NewFetchRawOptions(streams genericclioptions.IOStreams) *FetchRawOptions {
	return &FetchRawOptions{
		CommandContext: common.CommandContext{
			ConfigFlags: genericclioptions.NewConfigFlags(true),
			IOStreams:   streams,
		},
	}
}

// Validate ensures that all required arguments and flag values are provided
func (o *FetchRawOptions) Validate() error {
	finfo, err := os.Stat(o.OutputPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("The directory at path '%s' doesn't exist", o.OutputPath)
	}

	if !finfo.IsDir() {
		return fmt.Errorf("The output path must be a directory")
	}

	objref, err := common.ValidateObjectArgs(o.Args)
	if err != nil {
		return err
	}

	if o.HTML {
		_, err := exec.LookPath("oscap")
		if err != nil {
			return fmt.Errorf("The oscap command is needed for getting the HTML output")
		}
	}

	switch objref.Type {
	case common.ScanSettingBinding:
		o.Helper = NewScanSettingBindingHelper(o.Kuser, objref.Name, o.OutputPath, o.HTML, o.IOStreams)
	case common.ComplianceSuite:
		o.Helper = NewComplianceSuiteHelper(o.Kuser, objref.Name, o.OutputPath, o.HTML, o.IOStreams)
	case common.ComplianceScan:
		o.Helper = NewComplianceScanHelper(o.Kuser, objref.Name, o.OutputPath, o.HTML, o.IOStreams)
	default:
		return fmt.Errorf("Invalid object type for this command")
	}

	return nil
}

func (o *FetchRawOptions) Run() error {
	return o.Helper.Handle()
}
