//
//Copyright [2016] [SnapRoute Inc]
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//	 Unless required by applicable law or agreed to in writing, software
//	 distributed under the License is distributed on an "AS IS" BASIS,
//	 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//	 See the License for the specific language governing permissions and
//	 limitations under the License.
//
// _______  __       __________   ___      _______.____    __    ____  __  .___________.  ______  __    __
// |   ____||  |     |   ____\  \ /  /     /       |\   \  /  \  /   / |  | |           | /      ||  |  |  |
// |  |__   |  |     |  |__   \  V  /     |   (----` \   \/    \/   /  |  | `---|  |----`|  ,----'|  |__|  |
// |   __|  |  |     |   __|   >   <       \   \      \            /   |  |     |  |     |  |     |   __   |
// |  |     |  `----.|  |____ /  .  \  .----)   |      \    /\    /    |  |     |  |     |  `----.|  |  |  |
// |__|     |_______||_______/__/ \__\ |_______/        \__/  \__/     |__|     |__|      \______||__|  |__|
//

package server

import (
	"fmt"
	"l2/lldp/utils"
	"models"
	"utils/dbutils"
)

func (svr *LLDPServer) InitDB() error {
	var err error
	debug.Logger.Info("Initializing DB")
	svr.lldpDbHdl = dbutils.NewDBUtil(debug.Logger)
	err = svr.lldpDbHdl.Connect()
	if err != nil {
		debug.Logger.Err(fmt.Sprintln("Failed to Create DB Handle", err))
		return err
	}
	debug.Logger.Info(fmt.Sprintln("DB connection is established, error:", err))
	return nil
}

func (svr *LLDPServer) CloseDB() {
	debug.Logger.Info("Closed lldp db")
	svr.lldpDbHdl.Disconnect()
}

func (svr *LLDPServer) ReadDB() error {
	debug.Logger.Info("Reading from Database")
	if svr.lldpDbHdl == nil {
		debug.Logger.Info("Invalid db HDL")
		return nil
	}
	debug.Logger.Info("Getting objects from DB")
	var dbObj models.LLDPIntf
	objList, err := svr.lldpDbHdl.GetAllObjFromDb(dbObj)
	if err != nil {
		debug.Logger.Err(fmt.Sprintln("DB querry faile for LLDPIntf Config", err))
		return nil
	}
	// READ DB is always called before calling asicd get ports..
	debug.Logger.Info(fmt.Sprintln("Objects from db are", objList))
	for _, obj := range objList {
		dbEntry := obj.(models.LLDPIntf)
		gblInfo, _ := svr.lldpGblInfo[dbEntry.IfIndex]
		debug.Logger.Info(fmt.Sprintln("IfIndex", dbEntry.IfIndex, "is set to", dbEntry.Enable))
		switch dbEntry.Enable {
		case true:
			gblInfo.Enable()
		case false:
			gblInfo.Disable()
		}
		svr.lldpGblInfo[dbEntry.IfIndex] = gblInfo
	}

	debug.Logger.Info("Done reading from DB")
	return nil
}
