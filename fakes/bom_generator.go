package fakes

import "sync"

type BillOfMaterialGenerator struct {
	InstallAndRunCall struct {
		sync.Mutex
		CallCount int
		Receives  struct {
			WorkingDir string
		}
		Returns struct {
			String string
			Error  error
		}
		Stub func(string) (string, error)
	}
}

func (f *BillOfMaterialGenerator) InstallAndRun(param1 string) (string, error) {
	f.InstallAndRunCall.Lock()
	defer f.InstallAndRunCall.Unlock()
	f.InstallAndRunCall.CallCount++
	f.InstallAndRunCall.Receives.WorkingDir = param1
	if f.InstallAndRunCall.Stub != nil {
		return f.InstallAndRunCall.Stub(param1)
	}
	return f.InstallAndRunCall.Returns.String, f.InstallAndRunCall.Returns.Error
}
