package outputs

import (
	"github.com/fatih/color"
	"github.com/luckywinds/rshell/options"
	"github.com/luckywinds/rshell/outputs/outs"
	"github.com/luckywinds/rshell/types"
	"time"
)

func Output(result chan types.Hostresult, hg types.Hostgroup) {
	simpleOutput(result, hg)
}

func simpleOutput(result chan types.Hostresult, hg types.Hostgroup) {
	cfg := options.GetCfg()
	var taskresult types.Taskresult
	header := false
	for i := 0; i < len(hg.Ips); i++ {
		select {
		case res := <-result:
			taskresult.Results = append(taskresult.Results, res)
			if cfg.Outputintime == true {
				taskresult.Name = res.Actionname
				if !header {
					color.Yellow("TASK [%-30s] *****************************************\n", taskresult.Name + "@" + hg.Groupname)
					header = true
				}
				outFactory(cfg.Outputtype, res).PrintSimple()
			}
		case <-time.After(time.Duration(cfg.Tasktimeout) * time.Second):
			break
		}
	}

	if cfg.Outputintime == false {
		color.Yellow("TASK [%-30s] *****************************************\n", taskresult.Name + "@" + hg.Groupname)
		for _, value := range taskresult.Results {
			outFactory(cfg.Outputtype, value).PrintSimple()
		}
	}

	if len(taskresult.Results) == 0 {
		color.Red("%s\n", "ALL HOSTS TIMEOUT!!!")
	}

	if len(taskresult.Results) != len(hg.Ips) {
		m := make(map[string]types.Hostresult)
		for _, v := range taskresult.Results {
			m[v.Hostaddr] = v
		}

		for _, h := range hg.Ips {
			if _, ok := m[h]; !ok {
				outFactory(cfg.Outputtype, types.Hostresult{
					Actionname: taskresult.Results[0].Actionname,
					Actiontype: taskresult.Results[0].Actiontype,
					Groupname:  hg.Groupname,
					Hostaddr:   h,
					Error:      "TIMEOUT",
					Stdout:     "",
					Stderr:     "",
				}).PrintSimple()
			}
		}
	}
}

func outFactory(t string, r types.Hostresult) outs.OUT {
	switch t {
	case "text":
		return outs.TEXT{
			Actionname: r.Actionname,
			Actiontype: r.Actiontype,
			Groupname: r.Groupname,
			Address: r.Hostaddr,
			Stdout: r.Stdout,
			Stderr: r.Stderr,
			Syserr: r.Error,
		}
	default:
		return outs.TEXT{
			Actionname: r.Actionname,
			Actiontype: r.Actiontype,
			Groupname: r.Groupname,
			Address: r.Hostaddr,
			Stdout: r.Stdout,
			Stderr: r.Stderr,
			Syserr: r.Error,
		}
	}
}
