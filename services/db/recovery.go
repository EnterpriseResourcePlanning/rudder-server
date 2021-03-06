package db

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"time"

	"github.com/rudderlabs/rudder-server/rruntime"
	"github.com/rudderlabs/rudder-server/services/alert"
	"github.com/rudderlabs/rudder-server/services/stats"

	"github.com/rudderlabs/rudder-server/config"
	"github.com/rudderlabs/rudder-server/utils/logger"
)

const (
	normalMode    = "normal"
	degradedMode  = "degraded"
	migrationMode = "migration"
)

type RecoveryHandler interface {
	RecordAppStart(int64)
	HasThresholdReached() bool
	Handle()
}

var CurrentMode string = normalMode // default mode

// RecoveryDataT : DS to store the recovery process data
type RecoveryDataT struct {
	StartTimes                      []int64
	ReadableStartTimes              []string
	DegradedModeStartTimes          []int64
	ReadableDegradedModeStartTimes  []string
	MigrationModeStartTimes         []int64
	ReadableMigrationModeStartTimes []string
	Mode                            string
}

func getRecoveryData() RecoveryDataT {
	storagePath := config.GetString("recovery.storagePath", "/tmp/recovery_data.json")
	data, err := ioutil.ReadFile(storagePath)
	if os.IsNotExist(err) {
		defaultRecoveryJSON := "{\"mode\":\"" + normalMode + "\"}"
		data = []byte(defaultRecoveryJSON)
	} else {
		if err != nil {
			panic(err)
		}
	}

	var recoveryData RecoveryDataT
	err = json.Unmarshal(data, &recoveryData)
	if err != nil {
		panic(err)
	}

	return recoveryData
}

func saveRecoveryData(recoveryData RecoveryDataT) {
	recoveryDataJSON, err := json.MarshalIndent(&recoveryData, "", " ")
	storagePath := config.GetString("recovery.storagePath", "/tmp/recovery_data.json")
	err = ioutil.WriteFile(storagePath, recoveryDataJSON, 0644)
	if err != nil {
		panic(err)
	}
}

// IsNormalMode checks if the current mode is normal
func IsNormalMode() bool {
	return CurrentMode == normalMode
}

/*
CheckOccurences : check if this occurred numTimes times in numSecs seconds
*/
func CheckOccurences(occurences []int64, numTimes int, numSecs int) (occurred bool) {

	sort.Slice(occurences, func(i, j int) bool {
		return occurences[i] < occurences[j]
	})

	recentOccurences := 0
	checkPointTime := time.Now().Unix() - int64(numSecs)

	for i := len(occurences) - 1; i >= 0; i-- {
		if occurences[i] < checkPointTime {
			break
		}
		recentOccurences++
	}
	if recentOccurences >= numTimes {
		occurred = true
	}
	return
}

func getForceRecoveryMode(forceNormal bool, forceDegraded bool) string {
	switch {
	case forceNormal:
		return normalMode
	case forceDegraded:
		return degradedMode
	}
	return ""

}

func getNextMode(currentMode string) string {
	switch currentMode {
	case normalMode:
		return degradedMode
	case degradedMode:
		return ""
	case migrationMode: //Staying in the migrationMode forever on repeated restarts.
		return migrationMode
	}
	return ""
}

func NewRecoveryHandler(recoveryData *RecoveryDataT) RecoveryHandler {
	var recoveryHandler RecoveryHandler
	switch recoveryData.Mode {
	case normalMode:
		recoveryHandler = &NormalModeHandler{recoveryData: recoveryData}
	case degradedMode:
		recoveryHandler = &DegradedModeHandler{recoveryData: recoveryData}
	case migrationMode:
		recoveryHandler = &MigrationModeHandler{recoveryData: recoveryData}
	default:
		panic("Invalid Recovery Mode " + recoveryData.Mode)
	}
	return recoveryHandler
}

func alertOps(mode string) {
	instanceName := config.GetEnv("INSTANCE_ID", "")

	alertManager, err := alert.New()
	if err != nil {
		logger.Errorf("Unable to initialize the alertManager: %s", err.Error())
	} else {
		alertManager.Alert(fmt.Sprintf("Dataplane server %s entered %s mode", instanceName, mode))
	}
}

// sendRecoveryModeStat sends the recovery mode metric every 10 seconds
func sendRecoveryModeStat() {
	recoveryModeStat := stats.NewStat("recovery.mode_normal", stats.GaugeType)
	for {
		time.Sleep(10 * time.Second)
		switch CurrentMode {
		case normalMode:
			recoveryModeStat.Gauge(1)
		case degradedMode:
			recoveryModeStat.Gauge(2)
		case migrationMode:
			recoveryModeStat.Gauge(4)
		}
	}
}

// HandleRecovery decides the recovery Mode in which app should run based on earlier crashes
func HandleRecovery(forceNormal bool, forceDegraded bool, forceMigrationMode string, currTime int64) {

	enabled := config.GetBool("recovery.enabled", false)
	if !enabled {
		return
	}

	var forceMode string
	isForced := false

	//If MIGRATION_MODE environment variable is present and is equal to "import", "export", "import-export", then server mode is forced to be Migration.
	if IsValidMigrationMode(forceMigrationMode) {
		logger.Info("Setting server mode to Migration. If this is not intended remove environment variables related to Migration.")
		forceMode = migrationMode
	} else {
		forceMode = getForceRecoveryMode(forceNormal, forceDegraded)
	}

	recoveryData := getRecoveryData()
	if forceMode != "" {
		isForced = true
		recoveryData.Mode = forceMode
	} else {
		//If no mode is forced (through env or cli) and if previous mode is migration then setting server mode to normal.
		if recoveryData.Mode == migrationMode {
			recoveryData.Mode = normalMode
		}
	}
	recoveryHandler := NewRecoveryHandler(&recoveryData)

	if !isForced && recoveryHandler.HasThresholdReached() {
		logger.Info("DB Recovery: Moving to next State. Threshold reached for " + recoveryData.Mode)
		nextMode := getNextMode(recoveryData.Mode)
		if nextMode == "" {
			logger.Fatal("Threshold reached for degraded mode")
			panic("Not a valid mode")
		} else {
			recoveryData.Mode = nextMode
			recoveryHandler = NewRecoveryHandler(&recoveryData)
			alertOps(recoveryData.Mode)
		}
	}

	recoveryHandler.RecordAppStart(currTime)
	saveRecoveryData(recoveryData)
	recoveryHandler.Handle()
	logger.Infof("Starting in %s mode", recoveryData.Mode)
	CurrentMode = recoveryData.Mode
	rruntime.Go(func() {
		sendRecoveryModeStat()
	})
}
